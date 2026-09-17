package membertopup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type MemberTopupService struct {
	DB *bun.DB
}

func NewService(db *bun.DB) *MemberTopupService {
	return &MemberTopupService{DB: db}
}

func paymentBaseURL() string {
	return "http://" + os.Getenv("APP_PAYMENT_HOST") + ":" + os.Getenv("APP_PAYMENT_PORT")
}

// CreateTopup: bikin percobaan top-up baru. Tunai (PaymentGatewayCode kosong) langsung final
// -- insert member_topup_online + member_balance_ledger dalam 1 transaksi. Gateway
// (PaymentGatewayCode keisi) cuma insert member_topup_online 'pending' + minta QR ke service
// payment -- member_balance_ledger BELUM disentuh sampai settlement confirmed (lihat
// CheckStatus()).
func (s *MemberTopupService) CreateTopup(c *gin.Context, branchID int64, req CreateTopupRequestDTO) (*CreateTopupResponseDTO, error) {
	if req.PhoneNumber == "" {
		return nil, errors.New("phone_number wajib diisi")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount wajib lebih dari 0")
	}
	if req.Source != "pos" && req.Source != "kiosk" && req.Source != "mobile" {
		return nil, errors.New("source tidak valid")
	}

	memberID, err := s.resolveMemberIDByPhone(c, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	referenceNumber, err := s.generateTopupReference(c, branchID)
	if err != nil {
		return nil, err
	}

	isCash := req.PaymentGatewayCode == nil || *req.PaymentGatewayCode == ""

	if isCash {
		return s.createCashTopup(c, branchID, memberID, referenceNumber, req)
	}
	return s.createGatewayTopup(c, branchID, memberID, referenceNumber, req)
}

// resolveMemberIDByPhone: SAMA PERSIS query yang dipakai member.CheckByPhone() -- caller
// (POS/Kiosk/Mobile) cuma pegang nomor HP, gak pernah pegang/dipercaya kirim member_id
// langsung. is_active = true digate juga di sini, sama alasan kayak CheckByPhone.
func (s *MemberTopupService) resolveMemberIDByPhone(c *gin.Context, phoneNumber string) (int64, error) {
	var memberID int64
	err := s.DB.NewRaw(`SELECT id FROM master_member WHERE phone_number = ? AND is_active = true`, phoneNumber).Scan(c, &memberID)
	if err != nil {
		return 0, errors.New("member tidak ditemukan, cek nomor HP")
	}
	return memberID, nil
}

// createCashTopup: tunai -- kasir/kiosk udah terima duitnya duluan (fisik/di luar sistem),
// jadi langsung final. member_topup_online + member_balance_ledger diinsert bareng, 1 transaksi,
// biar gak mungkin salah satu doang yang kesave.
func (s *MemberTopupService) createCashTopup(c *gin.Context, branchID int64, memberID int64, referenceNumber string, req CreateTopupRequestDTO) (*CreateTopupResponseDTO, error) {
	tx, err := s.DB.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	gagal := true
	defer func() {
		if gagal {
			tx.Rollback()
		}
	}()

	now := time.Now()
	amountStr := strconv.FormatFloat(req.Amount, 'f', 2, 64)

	topup := MemberTopupOnlineModel{
		MemberID:        memberID,
		BranchID:        &branchID,
		TerminalID:      req.TerminalID,
		ReferenceNumber: referenceNumber,
		Amount:          amountStr,
		Source:          req.Source,
		Status:          "paid",
		PaidAt:          &now,
		Notes:           req.Notes,
		CreatedAt:       now,
	}
	if _, err := tx.NewInsert().Model(&topup).Exec(c); err != nil {
		return nil, err
	}

	balanceAfter, err := s.lockMemberAndInsertLedger(c, tx, memberID, &branchID, req.TerminalID, referenceNumber, "topup", req.Source, amountStr, "0")
	if err != nil {
		return nil, err
	}

	gagal = false
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateTopupResponseDTO{
		ReferenceNumber: referenceNumber,
		Status:          "paid",
		BalanceAfter:    &balanceAfter,
	}, nil
}

// createGatewayTopup: minta QR ke service payment (order_id = referenceNumber). Belum ada
// perubahan saldo di titik ini -- member_topup_online 'pending', nunggu CheckStatus() confirm.
func (s *MemberTopupService) createGatewayTopup(c *gin.Context, branchID int64, memberID int64, referenceNumber string, req CreateTopupRequestDTO) (*CreateTopupResponseDTO, error) {
	now := time.Now()
	amountStr := strconv.FormatFloat(req.Amount, 'f', 2, 64)

	topup := MemberTopupOnlineModel{
		MemberID:           memberID,
		BranchID:           &branchID,
		TerminalID:         req.TerminalID,
		ReferenceNumber:    referenceNumber,
		Amount:             amountStr,
		Source:             req.Source,
		PaymentGatewayCode: req.PaymentGatewayCode,
		Status:             "pending",
		Notes:              req.Notes,
		CreatedAt:          now,
	}
	if _, err := s.DB.NewInsert().Model(&topup).Exec(c); err != nil {
		return nil, err
	}

	amountInt := int64(req.Amount)
	body, _ := json.Marshal(map[string]any{
		"order_id":             referenceNumber,
		"payment_gateway_code": *req.PaymentGatewayCode,
		"amount":               amountInt,
		"branch_id":            branchID,
	})

	httpRes, err := http.Post(paymentBaseURL()+"/payment-gateway/qris", "application/json", bytes.NewReader(body))
	if err != nil {
		// gagal minta QR -- attempt ini gak jadi kepakai, jangan nyangkut 'pending' palsu.
		_, _ = s.DB.NewUpdate().Model((*MemberTopupOnlineModel)(nil)).
			Set("status = ?", "failed").
			Where("reference_number = ?", referenceNumber).Exec(c)
		return nil, fmt.Errorf("gagal menghubungi payment gateway: %w", err)
	}
	defer httpRes.Body.Close()

	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	parsed := paymentGatewayCreateQrisResponse{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if parsed.Code != 0 {
		_, _ = s.DB.NewUpdate().Model((*MemberTopupOnlineModel)(nil)).
			Set("status = ?", "failed").
			Where("reference_number = ?", referenceNumber).Exec(c)
		return nil, errors.New(parsed.Message)
	}

	var expiredAt *time.Time
	if parsed.Data.ExpiredAt != nil {
		if t, err := time.Parse(time.RFC3339, *parsed.Data.ExpiredAt); err == nil {
			expiredAt = &t
		}
	}
	_, err = s.DB.NewUpdate().Model((*MemberTopupOnlineModel)(nil)).
		Set("expired_at = ?", expiredAt).
		Where("reference_number = ?", referenceNumber).Exec(c)
	if err != nil {
		return nil, err
	}

	expiredAtStr := ""
	if parsed.Data.ExpiredAt != nil {
		expiredAtStr = *parsed.Data.ExpiredAt
	}
	return &CreateTopupResponseDTO{
		ReferenceNumber: referenceNumber,
		Status:          "pending",
		QRString:        parsed.Data.VendorQRString,
		QRURL:           parsed.Data.VendorQRURL,
		ExpiredAt:       &expiredAtStr,
	}, nil
}

// CheckStatus: polling status topup gateway -- dipanggil Kiosk/POS sambil nunggu customer
// scan QR. Begitu settlement confirmed, baru insert ke member_balance_ledger (saldo update).
// Idempotency guard: kalau status lokal udah 'paid', gak query ulang ke gateway / gak insert
// ledger dobel.
func (s *MemberTopupService) CheckStatus(c *gin.Context, referenceNumber string) (*CheckTopupStatusResponseDTO, error) {
	topup := new(MemberTopupOnlineModel)
	err := s.DB.NewSelect().Model(topup).Where("reference_number = ?", referenceNumber).Scan(c)
	if err != nil {
		return nil, errors.New("topup tidak ditemukan")
	}

	if topup.Status == "paid" {
		balanceAfter, err := s.getLastBalance(c, s.DB, topup.MemberID)
		if err != nil {
			return nil, err
		}
		return s.buildPaidTopupResponse(c, topup, balanceAfter)
	}
	if topup.Status != "pending" {
		return &CheckTopupStatusResponseDTO{ReferenceNumber: referenceNumber, Status: topup.Status}, nil
	}

	httpRes, err := http.Get(paymentBaseURL() + "/payment-gateway/" + referenceNumber)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi payment gateway: %w", err)
	}
	defer httpRes.Body.Close()
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	parsed := paymentGatewayStatusResponse{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if parsed.Code != 0 {
		return nil, errors.New(parsed.Message)
	}

	status := parsed.Data.Status

	// FALLBACK: expired_at lokal udah lewat tapi gateway masih bilang 'pending' -- sama pola
	// yang dipakai PaymentGatewayServices::CheckStatus() di POS Laravel (webhook mungkin gak
	// akan pernah nyampe). Cancel ke gateway + treat expired.
	if status == "pending" && topup.ExpiredAt != nil && time.Now().After(*topup.ExpiredAt) {
		_, _ = http.Post(paymentBaseURL()+"/payment-gateway/"+referenceNumber+"/cancel", "application/json", nil)
		status = "expired"
	}

	if status == "settlement" {
		balanceAfter, err := s.confirmTopupPaid(c, topup)
		if err != nil {
			return nil, err
		}
		return s.buildPaidTopupResponse(c, topup, balanceAfter)
	}

	if status == "expired" || status == "cancel" || status == "failed" {
		now := time.Now()
		_, _ = s.DB.NewUpdate().Model((*MemberTopupOnlineModel)(nil)).
			Set("status = ?", status).
			Set("cancel_at = ?", now).
			Where("reference_number = ?", referenceNumber).Exec(c)
	}

	return &CheckTopupStatusResponseDTO{ReferenceNumber: referenceNumber, Status: status}, nil
}

// confirmTopupPaid: tandain member_topup_online 'paid' + insert member_balance_ledger, 1
// transaksi. Guard status = 'pending' di WHERE update-nya jaga-jaga race (2 request check-status
// bersamaan) -- kalau RowsAffected() = 0, berarti udah kepake proses lain duluan, ambil saldo
// terkini aja tanpa insert ledger baru.
func (s *MemberTopupService) confirmTopupPaid(c *gin.Context, topup *MemberTopupOnlineModel) (string, error) {
	tx, err := s.DB.BeginTx(c, nil)
	if err != nil {
		return "", err
	}
	gagal := true
	defer func() {
		if gagal {
			tx.Rollback()
		}
	}()

	now := time.Now()
	res, err := tx.NewUpdate().Model((*MemberTopupOnlineModel)(nil)).
		Set("status = ?", "paid").
		Set("paid_at = ?", now).
		Where("reference_number = ? AND status = ?", topup.ReferenceNumber, "pending").
		Exec(c)
	if err != nil {
		return "", err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		// udah keconfirm proses lain -- gak insert ledger dobel, cukup balikin saldo terkini.
		gagal = false
		tx.Rollback()
		return s.getLastBalance(c, s.DB, topup.MemberID)
	}

	balanceAfter, err := s.lockMemberAndInsertLedger(c, tx, topup.MemberID, topup.BranchID, topup.TerminalID, topup.ReferenceNumber, "topup", topup.Source, topup.Amount, "0")
	if err != nil {
		return "", err
	}

	gagal = false
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return balanceAfter, nil
}

// lockMemberAndInsertLedger: lock baris master_member (serialize operasi saldo per-member,
// nyegah race antar 2 transaksi barengan buat member yang sama), ambil balance_after terakhir,
// insert baris baru dengan balance_after = lama + in - out (dihitung di SQL, bukan Go float,
// biar presisi NUMERIC-nya gak keganggu).
func (s *MemberTopupService) lockMemberAndInsertLedger(c *gin.Context, tx bun.Tx, memberID int64, branchID *int64, terminalID *int64, referenceNumber, transactionType, source, balanceIn, balanceOut string) (string, error) {
	var dummy int64
	if err := tx.NewRaw(`SELECT id FROM master_member WHERE id = ? FOR UPDATE`, memberID).Scan(c, &dummy); err != nil {
		return "", fmt.Errorf("member tidak ditemukan: %w", err)
	}

	var balanceAfter string
	err := tx.NewRaw(`
		INSERT INTO member_balance_ledger
			(member_id, branch_id, terminal_id, transaction_type, source, reference_number, balance_in, balance_out, balance_after)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?,
			COALESCE((SELECT balance_after FROM member_balance_ledger WHERE member_id = ? AND is_deleted = false ORDER BY created_at DESC, id DESC LIMIT 1), 0) + ? - ?)
		RETURNING balance_after
	`, memberID, branchID, terminalID, transactionType, source, referenceNumber, balanceIn, balanceOut, memberID, balanceIn, balanceOut).Scan(c, &balanceAfter)
	if err != nil {
		return "", err
	}
	return balanceAfter, nil
}

// buildPaidTopupResponse: enrich response CheckStatus buat topup yang statusnya udah "paid" --
// member_name (join master_member), amount/payment_gateway_code/terminal_id (langsung dari
// topup yang udah di-load, gak query ulang), paid_at (waktu response ini dibangun -- topup.PaidAt
// di memory mungkin belum ke-refresh abis confirmTopupPaid update DB, jadi pakai time.Now() aja,
// akurat cukup buat kebutuhan struk). Dipakai POS (Laravel, stateless) buat resolve receipt_station
// + isi struk tanpa perlu query balik ke sini.
func (s *MemberTopupService) buildPaidTopupResponse(c *gin.Context, topup *MemberTopupOnlineModel, balanceAfter string) (*CheckTopupStatusResponseDTO, error) {
	var member struct {
		Name        string `bun:"name"`
		PhoneNumber string `bun:"phone_number"`
	}
	_ = s.DB.NewRaw(`SELECT name, phone_number FROM master_member WHERE id = ?`, topup.MemberID).Scan(c, &member)

	paidAtStr := time.Now().Format(time.RFC3339)

	return &CheckTopupStatusResponseDTO{
		ReferenceNumber:    topup.ReferenceNumber,
		Status:             "paid",
		BalanceAfter:       &balanceAfter,
		TerminalID:         topup.TerminalID,
		Amount:             &topup.Amount,
		MemberName:         &member.Name,
		MemberPhoneNumber:  &member.PhoneNumber,
		PaymentGatewayCode: topup.PaymentGatewayCode,
		PaidAt:             &paidAtStr,
	}, nil
}

func (s *MemberTopupService) getLastBalance(c *gin.Context, db *bun.DB, memberID int64) (string, error) {
	var balance string
	err := db.NewRaw(`
		SELECT balance_after FROM member_balance_ledger
		WHERE member_id = ? AND is_deleted = false
		ORDER BY created_at DESC, id DESC LIMIT 1
	`, memberID).Scan(c, &balance)
	if err != nil {
		return "0.00", nil // belum pernah ada transaksi saldo sama sekali -- bukan error
	}
	return balance, nil
}

// generateTopupReference: niru persis pola helpers.GenerateReff punya sudocore2 (PREFIX +
// kode_branch + tanggal + urutan 4-digit, reset tiap hari per branch) -- gak bisa import
// langsung (beda Go module/binary), tapi baca-tulis ke tabel sequence_date & Postgres SEQUENCE
// yang SAMA (APIANDORDER connect ke DB pusat yang sama kayak sudocore2), jadi nomornya tetep
// 1 rangkaian yang konsisten & gak akan tabrakan meskipun digenerate dari proses yang beda.
func (s *MemberTopupService) generateTopupReference(c *gin.Context, branchID int64) (string, error) {
	const prefix = "tu"

	var branchCode string
	if err := s.DB.NewRaw(`SELECT code FROM master_branch WHERE id = ?`, branchID).Scan(c, &branchCode); err != nil {
		return "", err
	}
	branchCode = strings.TrimSpace(branchCode)
	if branchCode == "" {
		return "", fmt.Errorf("branch id %d belum punya kode (master_branch.code kosong), tidak bisa generate reference number", branchID)
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(branchCode) {
		return "", errors.New("kode branch tidak valid")
	}

	sequenceName := fmt.Sprintf("%s_sequence_%d", prefix, branchID)
	todayStr := time.Now().Format("2006-01-02")
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var dateStr string
	var batch int

	tx, err := s.DB.BeginTx(c, nil)
	if err != nil {
		return "", err
	}
	gagal := true
	defer func() {
		if gagal {
			tx.Rollback()
		}
	}()

	_, err = tx.NewRaw(`
		INSERT INTO sequence_date (branch_id, transaction_code, date)
		VALUES (?, ?, ?)
		ON CONFLICT (branch_id, transaction_code)
		DO NOTHING
	`, branchID, prefix, yesterdayStr).Exec(c)
	if err != nil {
		return "", err
	}

	err = tx.NewRaw(`
		SELECT date FROM sequence_date
		WHERE branch_id = ? AND transaction_code = ?
		FOR UPDATE
	`, branchID, prefix).Scan(c, &dateStr)
	if err != nil {
		return "", err
	}
	dateStr = strings.Split(dateStr, "T")[0]

	if dateStr != todayStr {
		_, err = tx.NewRaw(`UPDATE sequence_date SET date = ? WHERE branch_id = ? AND transaction_code = ?`, todayStr, branchID, prefix).Exec(c)
		if err != nil {
			return "", err
		}
		_, err = tx.NewRaw(fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS %s`, sequenceName)).Exec(c)
		if err != nil {
			return "", err
		}
		_, err = tx.NewRaw(fmt.Sprintf(`SELECT setval('%s', 1, false)`, sequenceName)).Exec(c)
		if err != nil {
			return "", err
		}
	}

	_, _ = tx.NewRaw(fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS %s`, sequenceName)).Exec(c)

	err = tx.NewRaw(fmt.Sprintf(`SELECT nextval('%s')`, sequenceName)).Scan(c, &batch)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	gagal = false

	date := strings.ReplaceAll(todayStr, "-", "")
	reference := fmt.Sprintf("%s%s%s%04d", strings.ToUpper(prefix), strings.ToUpper(branchCode), date, batch)
	return reference, nil
}

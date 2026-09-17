package mobileorder

import (
	"context"

	"github.com/uptrace/bun"
)

// GetPending: kandidat order mobile yang siap ditarik POS -- branch_id cocok, status 'paid'
// (satu-satunya status yang lolos, order gagal/expired gak pernah nyampe sini), pulled_at masih
// NULL (belum pernah di-ack lewat Ack()). 4 query terpisah (bukan JOIN raksasa) biar gampang
// nge-nest Detail->Package dan Payment di Go -- volume order per polling kecil, gak masalah.
func GetPending(ctx context.Context, db *bun.DB, branchID int64) ([]PendingOrder, error) {
	var orders []PendingOrder
	err := db.NewRaw(`
		SELECT mo.order_number, mo.branch_id, mo.member_id, mo.visit_purpose_id, mo.order_type, mo.pax,
			mo.order_fee, mo.service_charge, mo.platform_fee, mo.delivery_cost,
			mo.sub_total, mo.total_discount, mo.total_tax, mo.total_billing, mo.flag_inclusive_tax,
			mo.customer_phone_number, mo.payment_number, mo.payment_at, mo.payment_notes, mo.created_at,
			COALESCE(mm.name, '') AS member_name
		FROM mb_order mo
		LEFT JOIN master_member mm ON mm.id = mo.member_id
		WHERE mo.branch_id = ? AND mo.status = 'paid' AND mo.pulled_at IS NULL
		ORDER BY mo.created_at ASC
	`, branchID).Scan(ctx, &orders)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return orders, nil
	}

	orderNumbers := make([]string, len(orders))
	paymentNumbers := make([]string, 0, len(orders))
	for i, o := range orders {
		orderNumbers[i] = o.OrderNumber
		if o.PaymentNumber != nil {
			paymentNumbers = append(paymentNumbers, *o.PaymentNumber)
		}
	}

	var details []PendingOrderDetail
	if err := db.NewRaw(`
		SELECT ulid, order_number, pricelist_detail_id, menu_id, category_id, subcategory_id, qty,
			flag_inclusive_tax, price, tax_id, tax_type, tax_rate, tax_amount, dpp, net_dpp,
			promo_id, discount_percent, discount_amount, total, notes
		FROM mb_order_detail
		WHERE order_number IN (?)
		ORDER BY created_at ASC
	`, bun.In(orderNumbers)).Scan(ctx, &details); err != nil {
		return nil, err
	}

	detailULIDs := make([]string, len(details))
	for i, d := range details {
		detailULIDs[i] = d.ULID
	}

	var packages []PendingOrderDetailPackage
	if len(detailULIDs) > 0 {
		if err := db.NewRaw(`
			SELECT ulid, mb_order_detail_ulid, menu_package_id, menu_id, category_id, subcategory_id,
				qty, flag_inclusive_tax, price, tax_id, tax_type, tax_rate, tax_amount, dpp, net_dpp,
				promo_id, discount_percent, discount_amount, total, notes
			FROM mb_order_detail_package
			WHERE mb_order_detail_ulid IN (?)
		`, bun.In(detailULIDs)).Scan(ctx, &packages); err != nil {
			return nil, err
		}
	}

	var payments []PendingOrderPayment
	if len(paymentNumbers) > 0 {
		if err := db.NewRaw(`
			SELECT ulid, payment_number, payment_method_id, payment_amount, voucher_code,
				percent_from_billing, payment_gateway_order_id
			FROM mb_order_payment
			WHERE payment_number IN (?)
		`, bun.In(paymentNumbers)).Scan(ctx, &payments); err != nil {
			return nil, err
		}
	}

	packagesByDetail := make(map[string][]PendingOrderDetailPackage, len(details))
	for _, p := range packages {
		packagesByDetail[p.MbOrderDetailULID] = append(packagesByDetail[p.MbOrderDetailULID], p)
	}
	for i := range details {
		details[i].Package = packagesByDetail[details[i].ULID]
	}

	detailsByOrder := make(map[string][]PendingOrderDetail, len(orders))
	for _, d := range details {
		detailsByOrder[d.OrderNumber] = append(detailsByOrder[d.OrderNumber], d)
	}
	paymentByNumber := make(map[string]PendingOrderPayment, len(payments))
	for _, p := range payments {
		paymentByNumber[p.PaymentNumber] = p
	}

	for i := range orders {
		orders[i].Detail = detailsByOrder[orders[i].OrderNumber]
		if orders[i].PaymentNumber != nil {
			if p, ok := paymentByNumber[*orders[i].PaymentNumber]; ok {
				payment := p
				orders[i].Payment = &payment
			}
		}
	}

	return orders, nil
}

// Ack: tandain order udah beres ditarik POS -- WAJIB dipanggil SETELAH tr_order lokal berhasil
// keinsert, bukan sebelum (kalau insert lokal gagal di tengah jalan, order ini masih harus
// nongol lagi di GetPending berikutnya). Idempotent: WHERE pulled_at IS NULL -- panggil ulang
// (retry karena response ack sempat gak sampai walau DB kepencet) gak ngubah apa-apa lagi &
// tetap dianggap sukses (no-op), gak perlu bedain "order gak ketemu" vs "udah pernah di-ack".
func Ack(ctx context.Context, db *bun.DB, orderNumber string) error {
	_, err := db.NewRaw(`
		UPDATE mb_order SET pulled_at = now() WHERE order_number = ? AND pulled_at IS NULL
	`, orderNumber).Exec(ctx)
	return err
}

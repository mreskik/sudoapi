package member

import (
	"context"
	"database/sql"
	"errors"
)

// getMemberDetailHeader: header member + tier_name (JOIN master_member_tier_setting_detail by
// level) -- SAMA kolom query View() sudocore2 (member_services.go) ditambah tier_name. Member gak
// ketemu -> sql.ErrNoRows diteruskan apa adanya, caller (handler) yang mutusin itu bukan error.
func (h *MemberHandler) getMemberDetailHeader(ctx context.Context, memberId int) (MemberDetailResult, error) {
	result := MemberDetailResult{}
	err := h.DB.NewRaw(`
		SELECT mm.id, mm.code, mm.name, COALESCE(mm.contact_name, '') as contact_name,
		       COALESCE(mm.email, '') as email, COALESCE(mm.phone_number, '') as phone_number,
		       to_char(mm.birth_date, 'YYYY-MM-DD') as birth_date,
		       to_char(mm.created_at, 'YYYY-MM-DD') as join_date,
		       COALESCE(mm.description, '') as note,
		       mm.tier_level, mmtsd.name as tier_name
		FROM master_member mm
		LEFT JOIN master_member_tier_setting_detail mmtsd ON mmtsd.level = mm.tier_level
		WHERE mm.id = ? AND mm.is_active = true
	`, memberId).Scan(ctx, &result)
	return result, err
}

// getMemberPoint: snapshot poin terkini -- adopsi persis GetPoint() sudocore2
// (member_services.go:470-484). Belum ada histori poin sama sekali -- bukan error, "0".
func (h *MemberHandler) getMemberPoint(ctx context.Context, memberId int) (string, error) {
	var point string
	err := h.DB.NewRaw(`
		SELECT balance_after FROM member_point_ledger
		WHERE member_id = ? AND is_deleted = false
		ORDER BY created_at DESC, id DESC LIMIT 1
	`, memberId).Scan(ctx, &point)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "0", nil
		}
		return "", err
	}
	return point, nil
}

// getMemberLastOrder: breakdown LENGKAP order TERBARU member (status='paid') -- field header
// FLAT (bukan wrapper "header", 2026-10-30 permintaan eksplisit) + item_list[] (nested
// packages[] per item) + payment_method_list[], MIRROR PERSIS GetOrderHistoryDetail() sudocore2
// (member_services.go:566-629) + query payment dari posprint_service.go::ResolveReceiptData.
// 4 tahap: (1) cari order_number+payment_number TERBARU dulu, (2) header lengkap order itu,
// (3) items+packages, (4) payment method (kalau payment_number ada). Member belum pernah order
// 'paid' sama sekali -- semua field header nil, ItemList/PaymentMethodList slice kosong (BUKAN
// error).
func (h *MemberHandler) getMemberLastOrder(ctx context.Context, memberId int) (MemberLastOrderResult, error) {
	empty := MemberLastOrderResult{ItemList: []MemberLastOrderItemResult{}, PaymentMethodList: []MemberLastOrderPaymentResult{}}

	var latest struct {
		OrderNumber   string  `bun:"order_number"`
		PaymentNumber *string `bun:"payment_number"`
	}
	err := h.DB.NewRaw(`
		SELECT order_number, payment_number FROM pos_order
		WHERE member_id = ? AND status = 'paid'
		ORDER BY order_out DESC
		LIMIT 1
	`, memberId).Scan(ctx, &latest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return empty, nil
		}
		return empty, err
	}
	orderNumber := latest.OrderNumber

	result := MemberLastOrderResult{}
	err = h.DB.NewRaw(`
		SELECT po.order_number, po.branch_id, mb.name AS branch_name, po.order_out, po.status,
		       po.order_type, po.order_source,
		       po.sub_total, po.total_discount, po.total_tax, po.total_billing
		FROM pos_order po
		LEFT JOIN master_branch mb ON mb.id = po.branch_id
		WHERE po.order_number = ?
	`, orderNumber).Scan(ctx, &result)
	if err != nil {
		return empty, err
	}

	items := []MemberLastOrderItemResult{}
	err = h.DB.NewRaw(`
		SELECT pod.ulid, mi.item_name AS menu_name, mp.name AS promo_name,
		       pod.qty, pod.dpp, pod.tax_rate, pod.discount_amount, pod.total, pod.notes
		FROM pos_order_detail pod
		JOIN master_item_conversion_detail micd ON micd.id = pod.menu_id
		JOIN master_item mi ON mi.id = micd.item_id
		LEFT JOIN master_promo mp ON mp.id = pod.promo_id
		WHERE pod.order_number = ?
		ORDER BY pod.created_at ASC
	`, orderNumber).Scan(ctx, &items)
	if err != nil {
		return empty, err
	}

	if len(items) > 0 {
		packages := []MemberLastOrderPackageResult{}
		err = h.DB.NewRaw(`
			SELECT podp.tr_order_detail_ulid, mi.item_name AS menu_name, mp.name AS promo_name,
			       podp.qty, podp.dpp, podp.tax_rate, podp.discount_amount, podp.total, podp.notes
			FROM pos_order_detail_package podp
			JOIN master_item_conversion_detail micd ON micd.id = podp.menu_id
			JOIN master_item mi ON mi.id = micd.item_id
			LEFT JOIN master_promo mp ON mp.id = podp.promo_id
			WHERE podp.tr_order_detail_ulid IN (
				SELECT ulid FROM pos_order_detail WHERE order_number = ?
			)
		`, orderNumber).Scan(ctx, &packages)
		if err != nil {
			return empty, err
		}

		packagesByULID := map[string][]MemberLastOrderPackageResult{}
		for _, p := range packages {
			packagesByULID[p.TrOrderDetailULID] = append(packagesByULID[p.TrOrderDetailULID], p)
		}
		for i := range items {
			items[i].Packages = packagesByULID[items[i].ULID]
			if items[i].Packages == nil {
				items[i].Packages = []MemberLastOrderPackageResult{}
			}
		}
	}
	result.ItemList = items

	payments := []MemberLastOrderPaymentResult{}
	if latest.PaymentNumber != nil {
		err = h.DB.NewRaw(`
			SELECT mpm.name AS payment_method_name, pop.payment_amount
			FROM pos_order_payment pop
			JOIN master_payment_method mpm ON mpm.id = pop.payment_method_id
			WHERE pop.payment_number = ?
		`, *latest.PaymentNumber).Scan(ctx, &payments)
		if err != nil {
			return empty, err
		}
	}
	result.PaymentMethodList = payments

	return result, nil
}

// getMemberBestSeller: top 5 menu by qty -- adopsi persis GetBestSeller() sudocore2
// (member_services.go:640-668) + LIMIT 5 ditambahkan (sudocore2 sengaja tanpa limit, endpoint
// ini butuh top 5 pasti).
func (h *MemberHandler) getMemberBestSeller(ctx context.Context, memberId int) ([]MemberBestSellerResult, error) {
	list := []MemberBestSellerResult{}
	err := h.DB.NewRaw(`
		SELECT mi.item_name, mic.name AS category, misc.name AS sub_category,
		       SUM(pod.qty)::text AS qty,
		       MAX(po.order_out) AS last_purchased_at
		FROM pos_order_detail pod
		JOIN pos_order po ON po.order_number = pod.order_number
		JOIN master_item_conversion_detail micd ON micd.id = pod.menu_id
		JOIN master_item mi ON mi.id = micd.item_id
		LEFT JOIN master_item_category mic ON mic.id = pod.category_id
		LEFT JOIN master_item_sub_category misc ON misc.id = pod.subcategory_id
		WHERE po.member_id = ? AND po.status = 'paid'
		GROUP BY mi.id, mi.item_name, mic.name, misc.name
		ORDER BY SUM(pod.qty) DESC
		LIMIT 5
	`, memberId).Scan(ctx, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

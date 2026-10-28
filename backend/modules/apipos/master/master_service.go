package master

import (
	"APIANDORDER/backend/helpers"
	"context"

	"github.com/uptrace/bun"
)

type MasterService struct {
	DB *bun.DB
}

func NewMasterService(db *bun.DB) *MasterService {
	return &MasterService{DB: db}
}

func (this *MasterService) GetbranchList(context context.Context, user_id int) ([]Branch, error) {

	branch_list := []Branch{}
	err := this.DB.NewRaw(`SELECT
	mabdb.branch_id as id,
	concat(mad.code,' - ', mb.name)as name
FROM
	master_user_access_detail_branch mabdb
	JOIN (
	SELECT
		muad.id,
		muad.user_id,
		muad.company_id,
		muad.flag_all_branch,
		mc.code 
	FROM
		master_user_access_detail muad
		JOIN master_company mc ON mc.ID = muad.company_id 
	WHERE
	muad.user_id = ? 
	AND muad.flag_all_branch = false
	) mad ON mad.id = mabdb.access_detail_id 
	JOIN master_branch mb on mb.id = mabdb.branch_id
	
	UNION
	
	SELECT 
	mb.id,
	concat(mc.code, ' - ', mb.name) as name
	FROM master_branch mb
	JOIN master_company mc on mc.id = mb.company_id
	JOIN master_user_access_detail muad on muad.company_id = mb.company_id
	WHERE muad.user_id = ? AND muad.flag_all_branch = true AND mb.status = '1'

	`, user_id, user_id).Scan(context, &branch_list)
	if err != nil {
		return branch_list, err
	}
	return branch_list, nil
}

// branchDataQuery: SELECT inti data branch (brand, address, printing, dkk) buat GetDataBranch
// -- TIDAK include token sama sekali, karena sumber token-nya beda tergantung caller (install
// generate baru, sync rutin pakai token yang request ini sendiri sudah pakai). Lihat
// GetDataBranch() dan GetDataBranchForInstall() di bawah.
func (this *MasterService) branchDataQuery(context context.Context, branch_id int) (BranchData, error) {
	BranchData := BranchData{}

	// JOIN master_brand (2026-10-12): di-ubah dari INNER JOIN jadi LEFT JOIN -- branch yang
	// brand_id-nya NULL/gak match baris master_brand manapun TETAP harus bisa ke-pull (data
	// branch lain tetap lengkap, field brand_* doang yang kosong), bukan GAGAL TOTAL kayak
	// sebelumnya (INNER JOIN bikin seluruh baris branch ilang kalau brand_id-nya gak match).
	err := this.DB.NewRaw(`
 SELECT

mb.id as branch_id,
mb.code as branch_code,
mb.name as branch_name,
mb.brand_id as brand_id,
COALESCE(mbd.code, '') as brand_code,
COALESCE(mbd.name, '') as brand_name,
COALESCE(mbd.brand_color, '') as brand_color,
COALESCE(mbd.primary_color, '') as primary_color,
COALESCE(mbd.secondary_color, '') as secondary_color,
mb.address as address,
mb.telp as phone,
mb.printing_header,
mb.printing_footer,
mb.company_id,
COALESCE(mb.logo_header_src, '') as logo_header_src,
COALESCE(mb.image_footer_src, '') as image_footer_src,
COALESCE(mbs.domain_receipt, '') as domain_receipt,
COALESCE(mbs.whatsapp_template_text, '') as whatsapp_template_text

FROM master_branch mb
LEFT JOIN master_brand mbd on mbd.id = mb.brand_id
LEFT JOIN master_branch_setting mbs on mbs.branch_id = mb.id

WHERE mb.id = ?
 `, branch_id).Scan(context, &BranchData)

	return BranchData, err
}

// GetDataBranch: dipanggil SYNC RUTIN (sync_pull/get_data_branch, lewat BranchTokenAuth, token
// header yang SUDAH tervalidasi/dipakai request ini sendiri) -- TIDAK generate token baru, TIDAK
// insert baris session baru. Token dibalikin APA ADANYA dari token yang dipakai request ini
// (currentToken, diteruskan dari handler yang baca header Authorization) -- POS (SetupServices::
// getDatabranch(), shared fungsi install+sync) nyimpen Token hasil response ini ke mr_branch.token
// lokal, jadi KALAU DIBALIKIN KOSONG bakal NIMPA token valid yang lagi dipakai device itu sendiri.
func (this *MasterService) GetDataBranch(context context.Context, branch_id int, currentToken string) (BranchData, error) {
	BranchData, err := this.branchDataQuery(context, branch_id)
	if err != nil {
		return BranchData, err
	}

	BranchData.Token = currentToken
	return BranchData, nil
}

// GetDataBranchForInstall: dipanggil SETUP WIZARD (setup/get_data_branch, login username/
// password) -- device BARU, belum punya token. Generate token baru (2026-10-16, sistem session
// dinamis multi-device) + insert 1 baris baru ke master_branch_pos_session per panggilan
// (BUKAN baca master_branch.token statis lagi) -- device_hostname dari body POS, created_by
// dari userId hasil login (SetupService.CekLogin). master_branch.token LAMA DIBIARKAN ADA di
// skema (belum dihapus), sekadar gak dipakai lagi di jalur ini.
func (this *MasterService) GetDataBranchForInstall(context context.Context, branch_id int, userId int64, deviceHostname string) (BranchData, error) {
	BranchData, err := this.branchDataQuery(context, branch_id)
	if err != nil {
		return BranchData, err
	}

	token, err := helpers.GenerateToken()
	if err != nil {
		return BranchData, err
	}

	session := BranchPosSessionModel{
		BranchId:   int64(branch_id),
		Token:      token,
		FlagActive: true,
		CreatedBy:  &userId,
	}
	if deviceHostname != "" {
		session.DeviceHostname = &deviceHostname
	}

	if _, err := this.DB.NewInsert().Model(&session).Exec(context); err != nil {
		return BranchData, err
	}

	BranchData.Token = token
	return BranchData, nil
}

func (this *MasterService) GetStationList(context context.Context, branch_id int) ([]Station, error) {
	StationList := []Station{}

	err := this.DB.NewRaw(`SELECT 
id,branch_id,name,printer_name,printer_connection,printing_mode,printer_type,port,auto_cut,cash_drawer,line_character
FROM master_station WHERE branch_id = ? and is_active = true`, branch_id).Scan(context, &StationList)
	if err != nil {
		return StationList, err
	}
	return StationList, nil
}

// GetMasterCategory: master_item_category sekarang shared/global (company_id dihapus dari
// tabelnya, sudocore2 migration 269) -- gak ada lagi join ke branch/company, branch_id cuma
// dipertahankan di signature biar caller gak perlu diubah.
func (this *MasterService) GetMasterCategory(context context.Context, branch_id int) ([]CategoryDTO, error) {
	CategoryList := []CategoryDTO{}

	err := this.DB.NewRaw(`
	SELECT
	mic.id,
	mic.name,
	mic.code,
	mic.type_category_id
	from master_item_category mic
	`).Scan(context, &CategoryList)

	if err != nil {
		return CategoryList, err
	}
	return CategoryList, nil
}

// GetMasterSubCategory: master_item_sub_category sekarang shared/global (company_id dihapus,
// sudocore2 migration 269).
func (this *MasterService) GetMasterSubCategory(context context.Context, branch_id int) ([]SubCategoryDTO, error) {
	data := []SubCategoryDTO{}
	err := this.DB.NewRaw(`SELECT
msc.id,
msc.name,
COALESCE(msc.icon_src, '') as icon_src,
COALESCE(msc.banner_src, '') as banner_src
from master_item_sub_category msc`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetTableSection(context context.Context, branch_id int) ([]TableSection, error) {
	data := []TableSection{}

	err := this.DB.NewRaw(`SELECT
id,branch_id,name,table_checker_station_id,main_checker_station_id,layout_width,layout_height,layout_x,layout_y,layout_image_src,is_active, type, can_hold, print_category_setting_link
FROM master_table_section WHERE branch_id = ?`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this MasterService) GetTable(context context.Context, branch_id int) ([]TableSectionTable, error) {
	data := []TableSectionTable{}

	err := this.DB.NewRaw(`SELECT

ta.id,
ta.table_section_id,
ta.name,
ta.type,
ta.height,
ta.width,
ta.pos_x,
ta.pos_y,
ta.table_seat,
ta.minimum_billing,
ta.available_for_book

FROM master_table_section_layout_table ta
JOIN master_table_section mts on mts.id = ta.table_section_id
WHERE mts.branch_id = ? and mts.is_active = true`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetMasterTax: master_tax sekarang shared/global (company_id dihapus, sudocore2 migration 275).
func (this *MasterService) GetMasterTax(context context.Context, branch_id int) ([]MasterTax, error) {
	data := []MasterTax{}
	err := this.DB.NewRaw(`SELECT
mt.id, mt.name, mt.rate
FROM master_tax mt`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterTerminal(context context.Context, branch_id int) ([]MasterTerminal, error) {
	data := []MasterTerminal{}
	err := this.DB.NewRaw(`SELECT
id,
name,
branch_id,
device_id,
pos_type_id,
is_active,
is_used,
table_section_id,
receipt_station_id,
flag_printer_frontend

from master_terminal_id
WHERE branch_id = ?`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//////

// GetItem: master_item sekarang shared/global (company_id dihapus, sudocore2 migration 268).
func (this *MasterService) GetItem(context context.Context, branch_id int) ([]MasterItem, error) {
	data := []MasterItem{}
	err := this.DB.NewRaw(`
	SELECT

mi.id,
mi.item_longname as short_name,
mi.item_name as name,
mi.item_code as code,
COALESCE(mi.item_description, '') as description,
mi.item_category as category_id,
mi.item_subcategory as subcategory_id,
mi.item_bom as bom_id,
mi.item_color as menu_color,
mi.image,
COALESCE(mi.icon_src, '') as icon_src,
mi.use_tax as tax_type

FROM master_item mi
WHERE mi.is_deleted = false and mi.item_type = 'menu'
	`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetItemConv(context context.Context, branch_id int) ([]ItemConv, error) {
	data := []ItemConv{}
	err := this.DB.NewRaw(`SELECT
micd.id,
micd.item_id
FROM master_item_conversion_detail micd
JOIN master_item mi on micd.item_id = mi.id`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetItemPackage(context context.Context, branch_id int) ([]ItemPackage, error) {
	data := []ItemPackage{}
	err := this.DB.NewRaw(`SELECT
mip.id,
mip.item_id
FROM master_item_package mip
JOIN master_item mi on mi.id = mip.item_id`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetItemPackageGroup(context context.Context, branch_id int) ([]ItemPackageGroup, error) {
	data := []ItemPackageGroup{}
	err := this.DB.NewRaw(`SELECT
mipg.id,
mipg.item_package_id,
mipg.name,
mipg.min_qty,
mipg.max_qty

FROM master_item_package_group mipg
JOIN master_item_package mip on mip.id = mipg.item_package_id
JOIN master_item mi on mi.id = mip.item_id`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetItemPackageDetail(context context.Context, branch_id int) ([]ItemPackageDetail, error) {
	data := []ItemPackageDetail{}
	err := this.DB.NewRaw(`SELECT
mipd.id,
mipd.package_group_id,
mipd.item_conversion_detail_id as item_conv_detail_id,
mipd.price,
mipd.flag_all_menu_template,
mipd.default_item

FROM master_item_package_detail mipd
JOIN master_item_package_group mipg on mipg.id = mipd.package_group_id
JOIN master_item_package mip on mip.id = mipg.item_package_id
JOIN master_item mi on mi.id = mip.item_id`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

// GetItemPackageDetailPricelist: override harga sub-item package per pricelist (2026-08-26) --
// sumbernya master_item_package_detail_menu_template (ERP), field menu_template_id di-alias
// "pricelist_id" di response (lihat catatan di ItemPackageDetailPricelist DTO).
func (this *MasterService) GetItemPackageDetailPricelist(context context.Context, branch_id int) ([]ItemPackageDetailPricelist, error) {
	data := []ItemPackageDetailPricelist{}
	err := this.DB.NewRaw(`SELECT
mt.id,
mt.master_item_package_detail_id as item_package_detail_id,
mt.menu_template_id as pricelist_id,
mt.price

FROM master_item_package_detail_menu_template mt
JOIN master_item_package_detail mipd on mipd.id = mt.master_item_package_detail_id
JOIN master_item_package_group mipg on mipg.id = mipd.package_group_id
JOIN master_item_package mip on mip.id = mipg.item_package_id
JOIN master_item mi on mi.id = mip.item_id`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetPriceList(context context.Context, branch_id int) ([]MasterPricelist, error) {
	data := []MasterPricelist{}
	err := this.DB.NewRaw(`SELECT
mp.id,
mp.name,
mp.inclusive_price as is_inclusive

FROM master_pricelist mp`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetPriceListDetail(context context.Context, branch_id int) ([]masterPricelistDetail, error) {
	data := []masterPricelistDetail{}
	err := this.DB.NewRaw(`SELECT
mpd.id,
mpd.menu_template_id as pricelist_id,
mpd.item_conversion_detail_id as item_conv_detail_id,
mpd.price,
mpd.pos,
mpd.qr_order

FROM master_pricelist_detail mpd
join master_pricelist mp on mp.id = mpd.menu_template_id`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

// GetPriceListSubCategorySequence -- urutan tampil subcategory per menu template
// (master_pricelist_sub_category_sequence), ditarik bareng pricelist/pricelist detail
// (2026-10-15). Gak difilter branch_id (sama kayak GetPriceList/GetPriceListDetail di atas --
// menu template sekarang shared/global, lihat MASTER MENU TEMPLATE.md sudocore2).
func (this *MasterService) GetPriceListSubCategorySequence(context context.Context, branch_id int) ([]MasterPricelistSubCategorySequence, error) {
	data := []MasterPricelistSubCategorySequence{}
	err := this.DB.NewRaw(`SELECT
pscs.id,
pscs.menu_template_id as pricelist_id,
pscs.sub_category_id,
pscs.sequence

FROM master_pricelist_sub_category_sequence pscs`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

//payment method

func (this *MasterService) GetMasterPaymentMethod(context context.Context, branch_id int) ([]MasterPaymentMethod, error) {
	data := []MasterPaymentMethod{}
	err := this.DB.NewRaw(`
		SELECT

mpm.id,
mpm.name,
mpm.code,
mpm.coa_accout_id,
mpm.payment_method_type_id,
mpm.mdr,
mpm.printer_count,
mpm.open_cash_drawer,
mpm.group_payment_id,
mpm.color_theme,
mpm.use_authorization,
mpm.fixed_amount,
mpm.verification_code_mandatory,
mpm.card_number_code_mandatory,
mpm.payment_gateway_code

FROM  master_payment_method mpm
LEFT JOIN  master_payment_method_branches mpmd on mpmd.payment_method_id = mpm.id and mpmd.branch_id = ?
WHERE mpm.is_active = true and (mpm.flag_all_branch = true or mpmd.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPaymentMethodGroup(context context.Context, branch_id int) ([]MasterPaymentMethodGroup, error) {
	data := []MasterPaymentMethodGroup{}
	err := this.DB.NewRaw(`
SELECT
mpg.id, mpg.name
FROM master_payment_method_group mpg
WHERE mpg.is_active = TRUE
`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPaymentMethodType(context context.Context, branch_id int) ([]MasterPaymentMethodType, error) {
	data := []MasterPaymentMethodType{}
	err := this.DB.NewRaw(`
SELECT
id,name
FROM master_payment_method_type`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPaymentMethodVisitPurposes(context context.Context, branch_id int) ([]MasterPaymentMethodVisitPurposes, error) {
	data := []MasterPaymentMethodVisitPurposes{}
	err := this.DB.NewRaw(`SELECT
mpvp.id,
mpvp.payment_method_id,
mpvp.visitpurpose_id as visit_purpose_id
FROM master_payment_method_visit_purposes mpvp`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterBranchVisitPurpose(context context.Context, branch_id int) ([]MasterBranchVisitPurpose, error) {
	data := []MasterBranchVisitPurpose{}
	err := this.DB.NewRaw(`SELECT
id,
visit_purpose_id,
service_charge,
vat,
pb1,
order_fee,
menu_template_id as pricelist_id,
flag_kiosk
FROM master_branch_visit_purpose
WHERE branch_id = ? and is_active = true`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterBranchOpsSetting(context context.Context, branch_id int) ([]MasterBranchOpsSetting, error) {
	data := []MasterBranchOpsSetting{}
	err := this.DB.NewRaw(`SELECT id, day, status, open_time, closed_time
	FROM master_branch_ops_setting
	WHERE branch_id = ?`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

// GetMasterImage, GetMasterImageCustomerDisplay, GetMasterImageKiosk -- 3 endpoint flat
// terpisah buat data master_image per-channel (2026-08-24, ganti dari header -> image_list ->
// apply_for jadi header -> customer_display / kiosk, 2 tabel eksplisit sejajar, gak ada lagi
// join grandchild apply_for). Sama kayak master_promo: flag_all_branches -- kalau true, image
// ini berlaku lintas SEMUA branch/company (bukan cuma company branch itu sendiri), gak ada
// company_id di master_image.
func (this *MasterService) GetMasterImage(context context.Context, branch_id int) ([]MasterImage, error) {
	data := []MasterImage{}
	err := this.DB.NewRaw(`
	SELECT mi.id, mi.name, mi.is_active
	FROM master_image mi
	LEFT JOIN master_image_branches mib ON mib.master_image_id = mi.id AND mib.branch_id = ?
	WHERE mi.is_active = true AND (mi.flag_all_branches = true OR mib.id IS NOT NULL)
	`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterImageCustomerDisplay(context context.Context, branch_id int) ([]MasterImageCustomerDisplay, error) {
	data := []MasterImageCustomerDisplay{}
	err := this.DB.NewRaw(`
	SELECT micd.id, micd.master_image_id, micd.name, micd.banner_src, micd.sequence
	FROM master_image_customer_display micd
	JOIN master_image mi ON mi.id = micd.master_image_id
	LEFT JOIN master_image_branches mib ON mib.master_image_id = mi.id AND mib.branch_id = ?
	WHERE mi.is_active = true AND (mi.flag_all_branches = true OR mib.id IS NOT NULL)
	`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterImageKiosk(context context.Context, branch_id int) ([]MasterImageKiosk, error) {
	data := []MasterImageKiosk{}
	err := this.DB.NewRaw(`
	SELECT mik.id, mik.master_image_id, mik.name, mik.banner_src, mik.sequence
	FROM master_image_kiosk mik
	JOIN master_image mi ON mi.id = mik.master_image_id
	LEFT JOIN master_image_branches mib ON mib.master_image_id = mi.id AND mib.branch_id = ?
	WHERE mi.is_active = true AND (mi.flag_all_branches = true OR mib.id IS NOT NULL)
	`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterVisitPurpose(context context.Context, branch_id int) ([]MasterVisitPurpose, error) {
	data := []MasterVisitPurpose{}
	err := this.DB.NewRaw(`SELECT
vp.id,
vp.name
FROM master_visit_purpose vp
WHERE vp.is_active = true`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetTableSectionPrintCategorySetting(context context.Context, branch_id int) ([]MasterTableSectionPrintCategorySetting, error) {

	data_table_section, err := this.GetTableSection(context, branch_id)
	if err != nil {
		return nil, err
	}

	result := []MasterTableSectionPrintCategorySetting{}

	for _, table_section := range data_table_section {
		temp := []MasterTableSectionPrintCategorySetting{}

		// Jika print_category_setting_link null atau 0, ambil milik table section sendiri
		// Jika ada link, ambil milik table section yang dilink
		var target_table_section_id int
		if table_section.PrintCategorySettingLink == nil || *table_section.PrintCategorySettingLink == 0 {
			target_table_section_id = table_section.ID
		} else {
			target_table_section_id = int(*table_section.PrintCategorySettingLink)
		}

		err = this.DB.NewRaw(`
	SELECT
		mtsp.*
	FROM
		master_table_section_print_category_setting mtsp
	WHERE
		mtsp.table_section_id = ?
	`, target_table_section_id).Scan(context, &temp)

		if err != nil {
			return result, err
		}

		// Override table_section_id ke ID asli table section ini,
		// supaya consumer tahu setting ini milik table section yang mana
		for i := range temp {
			temp[i].TableSectionId = table_section.ID
		}

		result = append(result, temp...)
	}

	return result, nil
}

func (this *MasterService) GetMasterUserList(context context.Context, branch_id int) ([]MasterUser, error) {
	data := []MasterUser{}
	err := this.DB.NewRaw(`
SELECT DISTINCT
a.user_id as id,
c.username,
c.fullname,
c.role_id,
c.email,
c.sandi
FROM master_user_access_detail a
JOIN master_users c on c.id = a.user_id
JOIN master_user_access_detail_branch b on b.access_detail_id = a.id

WHERE a.flag_all_branch = false and b.branch_id = ? and c.flag_pos = true

UNION

SELECT DISTINCT
a.user_id as id,
c.username,
c.fullname,
c.role_id,
c.email,
c.sandi
FROM master_user_access_detail a
JOIN master_users c on c.id = a.user_id
JOIN master_branch b on b.id = ?
WHERE a.flag_all_branch = TRUE and c.flag_pos = TRUE AND a.company_id = b.company_id
	`, branch_id, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterRoleAccess(context context.Context, branch_id int) ([]MasterRoleAccess, error) {
	data := []MasterRoleAccess{}
	err := this.DB.NewRaw(`
		SELECT 
		mra.id,
		mra.role_id,
		mra.menu_id,
		mra.view,
		mra.insert,
		mra.update,
		mra.approve,
		mra.delete
		FROM master_role_access mra 
		JOIN (SELECT DISTINCT
		mu.role_id
		FROM master_users mu
		JOIN master_users_branches mub on mu.id = mub.user_id
		WHERE mub.branch_id = ? and mu.flag_pos = true) mu on mra.role_id = mu.role_id
	`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromo(context context.Context, branch_id int) ([]MasterPromo, error) {
	data := []MasterPromo{}
	err := this.DB.NewRaw(`
SELECT
mp.id,
mp.name,
mp.code,
mp.type,
mp.type_rupiah_amount,
mp.type_percent_use_limit,
mp.type_percent_rate,
mp.type_percent_limit_amount,
mp.type_freeitem_item_id,
mp.min_buy_amount,
mp.min_point_amount,
mp.flag_include_package,
mp.promo_for,
mp.apply_limit_per_day,
mp.apply_limit_per_item,
mp.period_start,
mp.period_end,
mp.flag_all_branches,
mp.flag_all_visit_purposes,
mp.flag_all_type_members,
mp.flag_all_tiers,
mp.flag_required_member,
mp.flag_all_days,
mp.flag_all_times,
mp.flag_apply_to_all,
mp.is_active,
mp.created_at,
mp.created_by,
mp.updated_at,
mp.updated_by

FROM master_promo mp
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoBranches(context context.Context, branch_id int) ([]MasterPromoBranches, error) {
	data := []MasterPromoBranches{}
	err := this.DB.NewRaw(`SELECT
mpb.id,
mpb.promo_id,
mpb.branch_id
FROM master_promo_branches mpb
JOIN master_promo mp on mp.id = mpb.promo_id
WHERE mpb.branch_id = ? and mp.is_active = true`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoVisitPurposes(context context.Context, branch_id int) ([]MasterPromoVisitPurposes, error) {
	data := []MasterPromoVisitPurposes{}
	err := this.DB.NewRaw(`SELECT
mpvp.id,
mpvp.promo_id,
mpvp.visit_purpose_id
FROM master_promo_visit_purposes mpvp
JOIN master_promo mp on mp.id = mpvp.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoTypeMembers(context context.Context, branch_id int) ([]MasterPromoTypeMembers, error) {
	data := []MasterPromoTypeMembers{}
	err := this.DB.NewRaw(`SELECT
mptm.id,
mptm.promo_id,
mptm.type_member_id
FROM master_promo_type_members mptm
JOIN master_promo mp on mp.id = mptm.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoTier(context context.Context, branch_id int) ([]MasterPromoTier, error) {
	data := []MasterPromoTier{}
	err := this.DB.NewRaw(`SELECT
mptr.id,
mptr.promo_id,
mptr.tier_level
FROM master_promo_tier mptr
JOIN master_promo mp on mp.id = mptr.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoCategories(context context.Context, branch_id int) ([]MasterPromoCategories, error) {
	data := []MasterPromoCategories{}
	err := this.DB.NewRaw(`SELECT
mpc.id,
mpc.promo_id,
mpc.category_id
FROM master_promo_categories mpc
JOIN master_promo mp on mp.id = mpc.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoSubCategories(context context.Context, branch_id int) ([]MasterPromoSubCategories, error) {
	data := []MasterPromoSubCategories{}
	err := this.DB.NewRaw(`SELECT
mpsc.id,
mpsc.promo_id,
mpsc.sub_category_id
FROM master_promo_sub_categories mpsc
JOIN master_promo mp on mp.id = mpsc.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoItems(context context.Context, branch_id int) ([]MasterPromoItems, error) {
	data := []MasterPromoItems{}
	err := this.DB.NewRaw(`SELECT
mpi.id,
mpi.promo_id,
mpi.item_id
FROM master_promo_items mpi
JOIN master_promo mp on mp.id = mpi.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoDays(context context.Context, branch_id int) ([]MasterPromoDays, error) {
	data := []MasterPromoDays{}
	err := this.DB.NewRaw(`SELECT
mpd.id,
mpd.promo_id,
mpd.day
FROM master_promo_days mpd
JOIN master_promo mp on mp.id = mpd.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoTimes(context context.Context, branch_id int) ([]MasterPromoTimes, error) {
	data := []MasterPromoTimes{}
	err := this.DB.NewRaw(`SELECT
mpt.id,
mpt.promo_id,
mpt.time_start,
mpt.time_end
FROM master_promo_times mpt
JOIN master_promo mp on mp.id = mpt.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterPromoApplyTo(context context.Context, branch_id int) ([]MasterPromoApplyTo, error) {
	data := []MasterPromoApplyTo{}
	err := this.DB.NewRaw(`SELECT
mpat.id,
mpat.promo_id,
mpat.apply_to
FROM master_promo_apply_to mpat
JOIN master_promo mp on mp.id = mpat.promo_id
LEFT JOIN master_promo_branches mpb on mpb.promo_id = mp.id and mpb.branch_id = ?
WHERE mp.is_active = true and (mp.flag_all_branches = true or mpb.id is not null)`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterMemberType(context context.Context, branch_id int) ([]MasterMemberType, error) {
	data := []MasterMemberType{}
	err := this.DB.NewRaw(`SELECT
id,
name,
is_active,
created_at,
created_by
FROM master_member_type
WHERE is_active = true`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterMember(context context.Context, branch_id int) ([]MasterMember, error) {
	data := []MasterMember{}
	err := this.DB.NewRaw(`SELECT
id,
member_type_id,
code,
name,
contact_name,
email,
phone_number,
tier_level,
is_active,
created_at,
created_by,
updated_at,
updated_by
FROM master_member
WHERE is_active = true`).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterMenuApp(context context.Context, branch_id int) ([]MasterMenuApp, error) {
	data := []MasterMenuApp{}
	err := this.DB.NewRaw(`
SELECT
id, menu, submenu
FROM master_menu`).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// /////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////

// GetMasterNotesMenu -- header, difilter flag_active + branch (SAMA POLA GetMasterPromo()):
// keluar kalau flag_all_branch = true ATAU ada baris master_notes_menu_branches buat branch
// ini. POS gak nerima flag_all_branch/flag_active sendiri (lihat MasterNotesMenu) -- baris ini
// nyampe ke POS artinya UDAH PASTI aktif & berlaku buat branch itu, gak perlu field tambahan.
func (this *MasterService) GetMasterNotesMenu(context context.Context, branch_id int) ([]MasterNotesMenu, error) {
	data := []MasterNotesMenu{}
	err := this.DB.NewRaw(`
SELECT
mnm.id,
mnm.name,
mnm.applies_to,
mnm.created_at,
mnm.created_by,
mnm.updated_at,
mnm.updated_by
FROM master_notes_menu mnm
LEFT JOIN master_notes_menu_branches mnb ON mnb.master_notes_menu_id = mnm.id AND mnb.branch_id = ?
WHERE mnm.flag_active = true AND (mnm.flag_all_branch = true OR mnb.id IS NOT NULL)`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// GetMasterNotesMenuCategories -- JOIN balik ke header + LEFT JOIN branches (sama pola
// GetMasterPromoVisitPurposes()) biar categories yang kepull cuma punya notes menu yang emang
// relevan buat branch ini (aktif + all_branch/match branch).
func (this *MasterService) GetMasterNotesMenuCategories(context context.Context, branch_id int) ([]MasterNotesMenuCategories, error) {
	data := []MasterNotesMenuCategories{}
	err := this.DB.NewRaw(`
SELECT
mnc.id,
mnc.master_notes_menu_id AS notes_menu_id,
mnc.category_id
FROM master_notes_menu_categories mnc
JOIN master_notes_menu mnm ON mnm.id = mnc.master_notes_menu_id
LEFT JOIN master_notes_menu_branches mnb ON mnb.master_notes_menu_id = mnm.id AND mnb.branch_id = ?
WHERE mnm.flag_active = true AND (mnm.flag_all_branch = true OR mnb.id IS NOT NULL)`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterNotesMenuSubCategories(context context.Context, branch_id int) ([]MasterNotesMenuSubCategories, error) {
	data := []MasterNotesMenuSubCategories{}
	err := this.DB.NewRaw(`
SELECT
mns.id,
mns.master_notes_menu_id AS notes_menu_id,
mns.sub_category_id
FROM master_notes_menu_subcategories mns
JOIN master_notes_menu mnm ON mnm.id = mns.master_notes_menu_id
LEFT JOIN master_notes_menu_branches mnb ON mnb.master_notes_menu_id = mnm.id AND mnb.branch_id = ?
WHERE mnm.flag_active = true AND (mnm.flag_all_branch = true OR mnb.id IS NOT NULL)`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterNotesMenuDetail(context context.Context, branch_id int) ([]MasterNotesMenuDetail, error) {
	data := []MasterNotesMenuDetail{}
	err := this.DB.NewRaw(`
SELECT
mnd.id,
mnd.master_notes_menu_id AS notes_menu_id,
mnd.short_notes,
mnd.full_notes
FROM master_notes_menu_detail mnd
JOIN master_notes_menu mnm ON mnm.id = mnd.master_notes_menu_id
LEFT JOIN master_notes_menu_branches mnb ON mnb.master_notes_menu_id = mnm.id AND mnb.branch_id = ?
WHERE mnm.flag_active = true AND (mnm.flag_all_branch = true OR mnb.id IS NOT NULL)`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

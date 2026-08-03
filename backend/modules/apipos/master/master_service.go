package master

import (
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

func (this *MasterService) GetDataBranch(context context.Context, branch_id int) (BranchData, error) {
	BranchData := BranchData{}

	err := this.DB.NewRaw(`
 SELECT

mb.id as branch_id,
mb.code as branch_code,
mb.name as branch_name,
mbd.code as brand_code,
mbd.name as brand_name,
mb.address as address,
mb.telp as phone,
mb.printing_header,
mb.printing_footer,
mb.company_id,
mb.token,
COALESCE(mb.logo_header_src, '') as logo_header_src,
COALESCE(mb.image_footer_src, '') as image_footer_src

FROM master_branch mb 
JOIN master_brand mbd on mbd.id = mb.brand_id

WHERE mb.id = ?
 `, branch_id).Scan(context, &BranchData)

	if err != nil {
		return BranchData, err
	}

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

// GET DATA CATEGORY FROM BRANCH REFERSE KE COMPANY ID
func (this *MasterService) GetMasterCategory(context context.Context, branch_id int) ([]CategoryDTO, error) {
	CategoryList := []CategoryDTO{}

	err := this.DB.NewRaw(`
	SELECT
	mic.id,
	mic.name,
	mic.code,
	mic.type_category_id
	from master_item_category mic
	JOIN (SELECT mb.company_id FROM master_branch mb WHERE mb.id = ?) mbb on mbb.company_id = mic.company_id
	`, branch_id).Scan(context, &CategoryList)

	if err != nil {
		return CategoryList, err
	}
	return CategoryList, nil
}

func (this *MasterService) GetMasterSubCategory(context context.Context, branch_id int) ([]SubCategoryDTO, error) {
	data := []SubCategoryDTO{}
	err := this.DB.NewRaw(`SELECT
msc.id,
msc.name
from master_item_sub_category msc
JOIN master_branch mb on msc.company_id = mb.company_id AND mb.id = ?`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetTableSection(context context.Context, branch_id int) ([]TableSection, error) {
	data := []TableSection{}

	err := this.DB.NewRaw(`SELECT
id,branch_id,name,table_checker_station_id,main_checker_station_id,layout_width,layout_height,layout_image_src,is_active, type, can_hold
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

func (this *MasterService) GetMasterTax(context context.Context, branch_id int) ([]MasterTax, error) {
	data := []MasterTax{}
	err := this.DB.NewRaw(`SELECT
mt.id, mt.name, mt.rate
FROM master_tax mt
JOIN master_branch mb on mb.company_id = mt.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)
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
is_used

from master_terminal_id
WHERE branch_id = ?`, branch_id).Scan(context, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//////

func (this *MasterService) GetItem(context context.Context, branch_id int) ([]MasterItem, error) {
	data := []MasterItem{}
	err := this.DB.NewRaw(`
	SELECT

mi.id,
mi.item_longname as short_name,
mi.item_name as name,
mi.item_code as code,
mi.item_category as category_id,
mi.item_subcategory as subcategory_id,
mi.item_bom as bom_id,
mi.item_color as menu_color,
mi.image,
mi.use_tax as tax_type

FROM master_item mi
JOIN master_branch mb on mb.company_id = mi.company_id
WHERE mb.id = ? and mi.is_deleted = false and mi.item_type = 'menu'
	`, branch_id).Scan(context, &data)
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
JOIN master_item mi on micd.item_id = mi.id
JOIN master_branch mb on mb.company_id = mi.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)
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
JOIN master_item mi on mi.id = mip.item_id
JOIN master_branch mb on mb.company_id = mi.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)
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
JOIN master_item mi on mi.id = mip.item_id
JOIN master_branch mb on mb.company_id = mi.company_id

WHERE mb.id = ?`, branch_id).Scan(context, &data)
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
mipd.price

FROM master_item_package_detail mipd
JOIN master_item_package_group mipg on mipg.id = mipd.package_group_id
JOIN master_item_package mip on mip.id = mipg.item_package_id
JOIN master_item mi on mi.id = mip.item_id
JOIN master_branch mb on mb.company_id = mi.company_id

WHERE mb.id = ?`, branch_id).Scan(context, &data)

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

FROM master_pricelist mp
JOIN master_branch mb on mb.company_id = mp.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)

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
join master_pricelist mp on mp.id = mpd.menu_template_id
JOIN master_branch mb on mb.company_id = mp.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)

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
mpm.card_number_code_mandatory

FROM  master_payment_method mpm
JOIN  master_payment_method_branches mpmd on mpmd.payment_method_id = mpm.id
WHERE mpmd.branch_id = ? and mpm.is_active = true`, branch_id).Scan(context, &data)

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
JOIN master_branch mb on mb.company_id = mpg.company_id
WHERE mb.id = ? and mpg.is_active = TRUE
`, branch_id).Scan(context, &data)

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
FROM master_payment_method_visit_purposes mpvp
JOIN master_branch mb on mb.company_id = mpvp.company_id
WHERE mb.id = ?`, branch_id).Scan(context, &data)

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
menu_template_id as pricelist_id
FROM master_branch_visit_purpose
WHERE branch_id = ? and is_active = true`, branch_id).Scan(context, &data)

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
JOIN master_branch mb on mb.company_id = vp.company_id
WHERE mb.id = ? and vp.is_active = true`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetTableSectionPrintCategorySetting(context context.Context, branch_id int) ([]MasterTableSectionPrintCategorySetting, error) {
	data := []MasterTableSectionPrintCategorySetting{}
	err := this.DB.NewRaw(`
	SELECT
		mtsp.* 
	FROM
		master_table_section_print_category_setting mtsp
		JOIN master_table_section mts ON mts.ID = mtsp.table_section_id 
	WHERE
		mts.branch_id = ?
	`, branch_id).Scan(context, &data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func (this *MasterService) GetMasterUserList (context context.Context, branch_id int)([]MasterUser, error){
	data :=  []MasterUser{}
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
	return data,nil
}

func (this *MasterService) GetMasterRoleAccess (context context.Context, branch_id int)([]MasterRoleAccess, error){
	data :=  []MasterRoleAccess{}
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
	return data,nil
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
mp.flag_all_days,
mp.flag_all_times,
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

func (this *MasterService) GetMasterMenuApp (context context.Context, branch_id int)([]MasterMenuApp, error){
	data := []MasterMenuApp{}
	err := this.DB.NewRaw(`
SELECT
id, menu, submenu
FROM master_menu`).Scan(context, &data)
if err != nil {
	return data, err
}
return data,nil
}
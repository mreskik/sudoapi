package master

import "time"

type CategoryDTO struct {
	ID             int    `bun:"id" json:"id"`
	Code           string `bun:"code" json:"code"`
	Name           string `bun:"name" json:"name"`
	TypeCategoryID int    `bun:"type_category_id" json:"type_category_id"`
}

type SubCategoryDTO struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}

type Branch struct {
	Id   int    `bun:"id"`
	Name string `bun:"name"`
}

type BranchData struct {
	BranchID       int    `bun:"branch_id"`
	BranchCode     string `bun:"branch_code"`
	BranchName     string `bun:"branch_name"`
	BrandCode      string `bun:"brand_code"`
	BrandName      string `bun:"brand_name"`
	Address        string `bun:"address"`
	Phone          string `bun:"phone"`
	PrintingHeader string `bun:"printing_header"`
	PrintingFooter string `bun:"printing_footer"`
	CompanyCode    string `bun:"company_code"`
	CompanyId      int    `bun:"company_id"`
	Token          string `bun:"token"`
	LogoHeaderSrc  string `bun:"logo_header_src" json:"LogoHeaderSrc"`
	ImageFooterSrc string `bun:"image_footer_src" json:"ImageFooterSrc"`
}

type Station struct {
	StationID         int    `bun:"id,pk,autoincrement"`
	BranchID          int    `bun:"branch_id,notnull"`
	StationName       string `bun:"name,notnull"`
	PrinterName       string `bun:"printer_name,nullzero"`
	PrinterType       int    `bun:"printer_type,nullzero"`
	PrinterConnection int    `bun:"printer_connection,nullzero"`
	PrintingMode      int    `bun:"printing_mode,nullzero"`
	Port              string `bun:"port,nullzero"`
	AutoCut           bool   `bun:"auto_cut,notnull,default:false"`
	CashDrawer        bool   `bun:"cash_drawer,notnull,default:false"`
	LineCharacter     int    `bun:"line_character,nullzero"`
}

type TableSection struct {
	ID                       int    `bun:"id,pk,autoincrement" json:"id"`
	BranchID                 int    `bun:"branch_id" json:"branch_id"`
	Name                     string `bun:"name" json:"name"`
	TableCheckerStationId    int    `bun:"table_checker_station_id" json:"tablechecker_station_id"`
	MainCheckerStationId     int    `bun:"main_checker_station_id" json:"mainchecker_station_id"`
	LayoutWidth              int    `bun:"layout_width" json:"layout_width"`
	LayoutHeight             int    `bun:"layout_height" json:"layout_height"`
	LayoutImageSrc           string `bun:"layout_image_src" json:"layout_image_src"`
	IsActive                 bool   `bun:"is_active" json:"is_active"`
	Type                     string `bun:"type" json:"type"`
	CanHold                  bool   `bun:"can_hold" json:"can_hold"`
	PrintCategorySettingLink *int64 `bun:"print_category_setting_link" json:"print_category_setting_link"`
}

type TableSectionTable struct {
	ID               int    `bun:"id,pk,autoincrement" json:"id"`
	TableSectionID   int    `bun:"table_section_id" json:"table_section_id"`
	Name             string `bun:"name" json:"name"`
	Type             string `bun:"type" json:"type"`
	Height           int    `bun:"height" json:"height"`
	Width            int    `bun:"width" json:"width"`
	PosX             int    `bun:"pos_x" json:"pos_x"`
	PosY             int    `bun:"pos_y" json:"pos_y"`
	TableSeat        int    `bun:"table_seat" json:"table_seat"`
	MinimumBilling   string `bun:"minimum_billing,type:numeric(18,2)" json:"minimum_billing"`
	AvailableForBook bool   `bun:"available_for_book" json:"available_for_book"`
}

type MasterTax struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
	Rate string `bun:"rate" json:"rate"`
}

type MasterTerminal struct {
	ID        int    `bun:"id" json:"id"`
	Name      string `bun:"name" json:"name"`
	BranchID  int    `bun:"branch_id" json:"branch_id"`
	DeviceID  string `bun:"device_id" json:"device_id"`
	POSTypeID int    `bun:"pos_type_id" json:"pos_type_id"`
	IsActive  bool   `bun:"is_active" json:"is_active"`
	IsUsed    bool   `bun:"is_used" json:"is_used"`
}

// /////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////

type MasterItem struct {
	ID            int    `bun:"id,pk" json:"id"`
	ShortName     string `bun:"short_name" json:"short_name"`
	Name          string `bun:"name" json:"name"`
	Code          string `bun:"code" json:"code"`
	CategoryID    int    `bun:"category_id" json:"category_id"`
	SubcategoryID int    `bun:"subcategory_id" json:"subcategory_id"`
	BomID         int    `bun:"bom_id" json:"bom_id"`
	MenuColor     string `bun:"menu_color" json:"menu_color"`
	Image         string `bun:"image" json:"image"`
	TaxType       string `bun:"tax_type" json:"tax_type"`
}

type ItemConv struct {
	ID     int `bun:"id" json:"id"`
	ItemId int `bun:"item_id" json:"item_id"`
}

type ItemPackage struct {
	ID     int `bun:"id" json:"id"`
	ItemId int `bun:"item_id" json:"item_id"`
}

type ItemPackageGroup struct {
	ID            int    `bun:"id" json:"id"`
	ItemPackageID int    `bun:"item_package_id" json:"item_package_id"`
	Name          string `bun:"name" json:"name"`
	MinQTY        int    `bun:"min_qty" json:"min_qty"`
	MaxQTY        int    `bun:"max_qty" json:"max_qty"`
}

type ItemPackageDetail struct {
	ID               int    `bun:"id" json:"id"`
	PackageGroupID   int    `bun:"package_group_id" json:"package_group_id"`
	ItemConvDetailID int    `bun:"item_conv_detail_id" json:"item_conv_detail_id"`
	Price            string `bun:"price,type:numeric(20,2)" json:"price"`
}

//

type MasterPricelist struct {
	ID          int    `bun:"id" json:"id"`
	Name        string `bun:"name" json:"name"`
	IsInclusive bool   `bun:"is_inclusive" json:"is_inclusive"`
}

type masterPricelistDetail struct {
	ID               int    `bun:"id" json:"id"`
	PricelistID      int    `bun:"pricelist_id" json:"pricelist_id"`
	ItemConvDetailID int    `bun:"item_conv_detail_id" json:"item_conv_detail_id"`
	Price            string `bun:"price" json:"price"`
	POS              bool   `bun:"pos" json:"pos"`
	QROrder          bool   `bun:"qr_order" json:"qr_order"`
}

///

type MasterPaymentMethod struct {
	ID                        int     `bun:"id,pk" json:"id"`
	Name                      string  `bun:"name" json:"name"`
	Code                      string  `bun:"code" json:"code"`
	CoaAccountID              int     `bun:"coa_accout_id" json:"coa_account_id"`
	PaymentMethodTypeID       int     `bun:"payment_method_type_id" json:"payment_method_type_id"`
	MDR                       string  `bun:"mdr" json:"mdr"`
	PrinterCount              int     `bun:"printer_count" json:"printer_count"`
	OpenCashDrawer            bool    `bun:"open_cash_drawer" json:"open_cash_drawer"`
	GroupPaymentID            int     `bun:"group_payment_id" json:"group_payment_id"`
	ColorTheme                *string `bun:"color_theme" json:"color_theme"`
	UseAuthorization          bool    `bun:"use_authorization" json:"use_authorization"`
	FixedAmount               string  `bun:"fixed_amount" json:"fixed_amount"`
	VerificationCodeMandatory bool    `bun:"verification_code_mandatory" json:"verification_code_mandatory"`
	CardNumberCodeMandatory   bool    `bun:"card_number_code_mandatory" json:"card_number_code_mandatory"`
}

type MasterPaymentMethodGroup struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}

type MasterPaymentMethodType struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}

type MasterPaymentMethodVisitPurposes struct {
	ID              int `bun:"id" json:"id"`
	PaymentMethodID int `bun:"payment_method_id" json:"payment_method_id"`
	VisitPurposeID  int `bun:"visit_purpose_id" json:"visit_purpose_id"`
}

type MasterBranchVisitPurpose struct {
	ID             int    `bun:"id" json:"id"`
	VisitPurposeID int    `bun:"visit_purpose_id" json:"visit_purpose_id"`
	ServiceCharge  int    `bun:"service_charge" json:"service_charge"`
	VAT            int    `bun:"vat" json:"vat"`
	PB1            int    `bun:"pb1" json:"pb1"`
	OrderFee       string `bun:"order_fee" json:"order_fee"`
	PriceListID    int    `bun:"pricelist_id" json:"pricelist_id"`
}

type MasterVisitPurpose struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}

type MasterTableSectionPrintCategorySetting struct {
	ID             int `bun:"id" json:"id"`
	TableSectionId int `bun:"table_section_id" json:"table_section_id"`
	SubCategoryId  int `bun:"sub_category_id" json:"sub_category_id"`
	StationId      int `bun:"station_id" json:"station_id"`
	CategoryId     int `bun:"category_id" json:"category_id"`
}

type MasterUser struct {
	ID       int    `bun:"id" json:"id"`
	Username string `bun:"username" json:"username"`
	Fullname string `bun:"fullname" json:"fullname"`
	RoleID   int    `bun:"role_id" json:"role_id"`
	Email    string `bun:"email" json:"email"`
	Sandi    string `bun:"sandi" json:"sandi"`
}

type MasterRoleAccess struct {
	ID      int  `bun:"id" json:"id"`
	RoleID  int  `bun:"role_id" json:"role_id"`
	MenuID  int  `bun:"menu_id" json:"menu_id"`
	View    bool `bun:"view" json:"view"`
	Insert  bool `bun:"insert" json:"insert"`
	Update  bool `bun:"update" json:"update"`
	Delete  bool `bun:"delete" json:"delete"`
	Approve bool `bun:"approve" json:"approve"`
}

type MasterMenuApp struct {
	ID      int    `bun:"id" json:"id"`
	Menu    string `bun:"menu" json:"menu"`
	SubMenu string `bun:"submenu" json:"submenu"`
}

// /////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////

type MasterPromo struct {
	ID   int    `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
	Code string `bun:"code" json:"code"`
	Type string `bun:"type" json:"type"` // rupiah, percent, freeitem

	TypeRupiahAmount       string  `bun:"type_rupiah_amount" json:"type_rupiah_amount"`
	TypePercentUseLimit    bool    `bun:"type_percent_use_limit" json:"type_percent_use_limit"`
	TypePercentRate        *string `bun:"type_percent_rate" json:"type_percent_rate,omitempty"`
	TypePercentLimitAmount string  `bun:"type_percent_limit_amount" json:"type_percent_limit_amount"`
	TypeFreeitemItemID     *int    `bun:"type_freeitem_item_id" json:"type_freeitem_item_id,omitempty"` // ITEM ID, bukan CONV ID

	MinBuyAmount   string `bun:"min_buy_amount" json:"min_buy_amount"`
	MinPointAmount int    `bun:"min_point_amount" json:"min_point_amount"`

	FlagIncludePackage bool   `bun:"flag_include_package" json:"flag_include_package"`
	PromoFor           string `bun:"promo_for" json:"promo_for"` // category, sub_category, item
	ApplyLimitPerDay   int    `bun:"apply_limit_per_day" json:"apply_limit_per_day"`
	ApplyLimitPerItem  int    `bun:"apply_limit_per_item" json:"apply_limit_per_item"`

	PeriodStart string `bun:"period_start" json:"period_start"` // DATE, format: "2006-01-02"
	PeriodEnd   string `bun:"period_end" json:"period_end"`     // DATE, format: "2006-01-02"

	FlagAllBranches      bool `bun:"flag_all_branches" json:"flag_all_branches"`
	FlagAllVisitPurposes bool `bun:"flag_all_visit_purposes" json:"flag_all_visit_purposes"`
	FlagAllTypeMembers   bool `bun:"flag_all_type_members" json:"flag_all_type_members"`

	FlagAllDays  bool `bun:"flag_all_days" json:"flag_all_days"`
	FlagAllTimes bool `bun:"flag_all_times" json:"flag_all_times"`

	IsActive bool `bun:"is_active" json:"is_active"`

	CreatedAt time.Time  `bun:"created_at" json:"created_at"`
	CreatedBy *int       `bun:"created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time `bun:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *int       `bun:"updated_by" json:"updated_by,omitempty"`

	Branches      []MasterPromoBranches      `json:"branches,omitempty"`
	VisitPurposes []MasterPromoVisitPurposes `json:"visit_purposes,omitempty"`
	TypeMembers   []MasterPromoTypeMembers   `json:"type_members,omitempty"`
	Categories    []MasterPromoCategories    `json:"categories,omitempty"`
	SubCategories []MasterPromoSubCategories `json:"sub_categories,omitempty"`
	Items         []MasterPromoItems         `json:"items,omitempty"`
	Days          []MasterPromoDays          `json:"days,omitempty"`
	Times         []MasterPromoTimes         `json:"times,omitempty"`
}

type MasterPromoBranches struct {
	ID       int `bun:"id" json:"id"`
	PromoID  int `bun:"promo_id" json:"promo_id"`
	BranchID int `bun:"branch_id" json:"branch_id"`
}

type MasterPromoVisitPurposes struct {
	ID             int `bun:"id" json:"id"`
	PromoID        int `bun:"promo_id" json:"promo_id"`
	VisitPurposeID int `bun:"visit_purpose_id" json:"visit_purpose_id"`
}

type MasterPromoTypeMembers struct {
	ID           int `bun:"id" json:"id"`
	PromoID      int `bun:"promo_id" json:"promo_id"`
	TypeMemberID int `bun:"type_member_id" json:"type_member_id"`
}

type MasterPromoCategories struct {
	ID         int `bun:"id" json:"id"`
	PromoID    int `bun:"promo_id" json:"promo_id"`
	CategoryID int `bun:"category_id" json:"category_id"`
}

type MasterPromoSubCategories struct {
	ID            int `bun:"id" json:"id"`
	PromoID       int `bun:"promo_id" json:"promo_id"`
	SubCategoryID int `bun:"sub_category_id" json:"sub_category_id"`
}

type MasterPromoItems struct {
	ID      int `bun:"id" json:"id"`
	PromoID int `bun:"promo_id" json:"promo_id"`
	ItemID  int `bun:"item_id" json:"item_id"`
}

type MasterPromoDays struct {
	ID      int    `bun:"id" json:"id"`
	PromoID int    `bun:"promo_id" json:"promo_id"`
	Day     string `bun:"day" json:"day"` // senin, selasa, rabu, kamis, jumat, sabtu, minggu
}

type MasterPromoTimes struct {
	ID        int    `bun:"id" json:"id"`
	PromoID   int    `bun:"promo_id" json:"promo_id"`
	TimeStart string `bun:"time_start" json:"time_start"` // TIME, format: "15:04:05"
	TimeEnd   string `bun:"time_end" json:"time_end"`     // TIME, format: "15:04:05"
}

// /////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////

type MasterMemberType struct {
	ID        int       `bun:"id" json:"id"`
	Name      string    `bun:"name" json:"name"`
	IsActive  bool      `bun:"is_active" json:"is_active"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
	CreatedBy *int      `bun:"created_by" json:"created_by,omitempty"`
}

type MasterMember struct {
	ID           int        `bun:"id" json:"id"`
	MemberTypeID *int       `bun:"member_type_id" json:"member_type_id,omitempty"`
	Code         string     `bun:"code" json:"code"`
	Name         string     `bun:"name" json:"name"`
	ContactName  string     `bun:"contact_name" json:"contact_name,omitempty"`
	Email        string     `bun:"email" json:"email,omitempty"`
	PhoneNumber  string     `bun:"phone_number" json:"phone_number,omitempty"`
	IsActive     bool       `bun:"is_active" json:"is_active"`
	CreatedAt    time.Time  `bun:"created_at" json:"created_at"`
	CreatedBy    *int       `bun:"created_by" json:"created_by,omitempty"`
	UpdatedAt    *time.Time `bun:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy    *int       `bun:"updated_by" json:"updated_by,omitempty"`
}

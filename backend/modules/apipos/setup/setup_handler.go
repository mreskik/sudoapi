package setup

import (
	"APIANDORDER/backend/helpers"
	"APIANDORDER/backend/modules/apipos/master"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type SetupHandler struct {
	setupService  *SetupService
	masterService *master.MasterService
}

func NewHandler(db *bun.DB) *SetupHandler {
	return &SetupHandler{setupService: NewSetupService(db), masterService: master.NewMasterService(db)}
}

func (this *SetupHandler) GetBranchList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}

	login_status, user_id, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	branch_list := []master.Branch{}
	if login_status {
		branch_list, _ = this.masterService.GetbranchList(c, user_id)
	} else {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	c.JSON(200, res.Success().SetData(branch_list))
}

func (this *SetupHandler) GetDataBranch(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	branch_data, err := this.masterService.GetDataBranch(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("error get data branch!"))
		return
	}

	c.JSON(200, res.Success().SetData(branch_data))
}

func (this *SetupHandler) GetStationList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	station_list, err := this.masterService.GetStationList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal tarik daftar station 1"))
		return
	}

	c.JSON(200, res.Success().SetData(station_list))
}

func (this *SetupHandler) GetCategoryList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_category, err := this.masterService.GetMasterCategory(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data category!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_category))
}
func (this *SetupHandler) GetSubCategoryList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_subcategory, err := this.masterService.GetMasterSubCategory(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data sub category!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_subcategory))
}

func (this *SetupHandler) GetTableSectionList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tablesection, err := this.masterService.GetTableSection(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table section!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tablesection))
}

func (this *SetupHandler) GetTable(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_table, err := this.masterService.GetTable(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table section!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_table))
}

func (this *SetupHandler) GetMasterTax(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterTax(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data table master tax!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterTerminal(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterTerminal(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  terminal!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetItem(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetItem(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  item!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetItemConv(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetItemConv(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  item conv!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetItemPackage(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetItemPackage(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  item package!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetItemPackageGroup(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetItemPackageGroup(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  item package group!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetItemPackageDetail(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetItemPackageDetail(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  ItemPackage detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetPriceList(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetPriceList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  pricelist!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetPriceListDetail(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetPriceListDetail(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  pricelist detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterPaymentMethod(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterPaymentMethod(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  payment method detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterPaymentMethodGroup(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterPaymentMethodGroup(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  payment method group detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterPaymentMethodType(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterPaymentMethodType(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  payment method type detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterPaymentMethodVisitPurposes(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterPaymentMethodVisitPurposes(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  Payment Method Visit Purpose detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterBranchVisitPurpose(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterBranchVisitPurpose(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  pricelist detail!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}

func (this *SetupHandler) GetMasterVisitPurpose(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterVisitPurpose(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  master visit purpose!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}



func (this *SetupHandler) GetMasterTableSectionPrintCategorySetting(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetTableSectionPrintCategorySetting(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  master table section print category setting!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}


func (this *SetupHandler) GetMasterUser(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterUserList(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  master user section print category setting!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}


func (this *SetupHandler) GetMasterRoleAccess(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterRoleAccess(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  master MasterRoleAccess!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}


func (this *SetupHandler) GetMenuApp(c *gin.Context) {
	res := helpers.NewResponse()

	var login Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(200, res.GeneralError())
		return
	}
	login_status, _, _ := this.setupService.CekLogin(c, login.Username, login.Password)
	if !login_status {
		c.JSON(200, res.GeneralError().SetMessage("username or password is incorrect!"))
		return
	}

	branch_id, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("branch id salah!"))
		return
	}

	data_tax, err := this.masterService.GetMasterMenuApp(c, branch_id)
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage("gagal ambil data  master MasterMenuApp!"))
		return
	}

	c.JSON(200, res.Success().SetData(data_tax))
}
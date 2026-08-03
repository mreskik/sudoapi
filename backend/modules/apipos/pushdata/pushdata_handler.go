package pushdata

import (
	"APIANDORDER/backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type PushDataHandler struct {
	PushDataService *PushDataService
}

func NewHander(db *bun.DB) *PushDataHandler {
	return &PushDataHandler{PushDataService: NewPushDataService(db)}
}


func (this *PushDataHandler) PushDataPosOrder(c *gin.Context){
	res:= helpers.NewResponse()
	list_data_order :=  PosOrderDTO{}

	err:= c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}


	err = this.PushDataService.PushPOSOrder(c, list_data_order)
	
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	
	c.JSON(200, res.SetCode(0).SetMessage("success push data data order!"))
	
}


func (this *PushDataHandler) PushDataPosOrderDetail(c *gin.Context){
	res:= helpers.NewResponse()
	list_data_order :=  PosOrderDetailDTO{}

	err:= c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}


	err = this.PushDataService.PushPOSOrderDetail(c, list_data_order)
	
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	
	c.JSON(200, res.SetMessage("success push data!"))
	
}


func (this *PushDataHandler) PushDataPosOrderDetailPackage(c *gin.Context){
	res:= helpers.NewResponse()
	list_data_order :=  PosOrderDetailPackageDTO{}

	err:= c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}


	err = this.PushDataService.PushPOSOrderDetailPackage(c, list_data_order)
	
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	
	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
	
}

func (this *PushDataHandler) PushDataPosOrderPayment(c *gin.Context){
	res:= helpers.NewResponse()
	list_data_order :=  PosOrderPaymentDTO{}

	err:= c.ShouldBindJSON(&list_data_order)
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}

	err = this.PushDataService.PushPOSOrderPayment(c, list_data_order)
	
	if err != nil {
		res.SetCode(100)
		c.JSON(200, res.SetMessage(err.Error()))
		return
	}
	
	c.JSON(200, res.SetMessage("success push data!").SetCode(0))
}
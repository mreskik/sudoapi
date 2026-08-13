package pushdata

type PosOrderDTO struct {
	ListOrder []PosOrderModel `json:"list_order"`
}

type PosOrderDetailDTO struct {
	ListOrderDetail []PosOrderDetailModel `json:"list_order_detail"`
}

type PosOrderDetailPackageDTO struct {
	ListOrderDetailPackage []PosOrderDetailPackageModel `json:"list_order_detail_package"`
}

type PosOrderPaymentDTO struct {
	ListOrderPayment []PosOrderPaymentModel `json:"list_order_payment"`
}

type PosDayShiftDTO struct {
	ListDayShift []PosDayShiftModel `json:"list_dayshift"`
}

type PosDayShiftDetailDTO struct {
	ListDayShiftDetail []PosDayShiftDetailModel `json:"list_dayshift_detail"`
}

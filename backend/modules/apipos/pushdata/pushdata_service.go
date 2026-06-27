package pushdata

import (
	"context"

	"github.com/uptrace/bun"
)

type PushDataService struct {
	DB *bun.DB
}

func NewPushDataService(db *bun.DB) *PushDataService {
	return &PushDataService{
		DB: db,
	}
}

type BranchDetails struct{
	Id int `bun:"id"`
	CompanyId int `bun:"company_id"`
	BranchName string `bun:"branch_name"`
}


func(this *PushDataService) PushPOSOrder(context context.Context,  data_order_request PosOrderDTO)error{

	tx, err := this.DB.BeginTx(context, nil)
	if err != nil {
		return err
	}
	
	
	if(len(data_order_request.ListOrder) >0){	
		
		DETAIL_BRANCH := BranchDetails{}
		err:= tx.NewRaw(`
		SELECT id,company_id,name as branch_name from master_branch where id = ?
		`, data_order_request.ListOrder[0].BranchID).Scan(context, &DETAIL_BRANCH)
		if err != nil {
			return err
		}
		
		// now := time.Now()
		// timestamp_now:= now.Format("20060102150405")+fmt.Sprintf("%06d", now.Nanosecond()/1_000)
		// // 20260624104530123456 contoh hasil

		// PREFIX := "POS_"+DETAIL_BRANCH.BranchName+timestamp_now

		// PREFIX = strings.ReplaceAll(PREFIX, " ", "")
		

		
	for _,v := range data_order_request.ListOrder {
		// v.EndDayReferenceNumber = &PREFIX
		v.CompanyId = &DETAIL_BRANCH.CompanyId
		_,err= tx.NewInsert().Model(&v).On("CONFLICT (order_number) DO UPDATE").Exec(context)
		if err != nil {
			return err
		}
	}
	}
	


	err= tx.Commit()
	if err != nil {
		return err
	}
	return nil
}



func(this *PushDataService) PushPOSOrderDetail(context context.Context, data_order_detail_request PosOrderDetailDTO)error{
	tx, err := this.DB.BeginTx(context, nil)
	if err != nil {
		return err
	}

	// var list_order_number  []string
	// for _,v := range data_order_detail_request.ListOrderDetail{
	// 	list_order_number = append(list_order_number, v.OrderNumber)
	// }
	// var list_order_number_merger []string

	// seen := make(map[string]struct{})

	// for _, item := range list_order_number {
	// 	if _, exists := seen[item]; !exists {
	// 		seen[item] = struct{}{}
	// 		list_order_number_merger = append(list_order_number_merger, item)
	// 	}
	// }
	

	for _,v := range data_order_detail_request.ListOrderDetail {
		_,err= tx.NewInsert().Model(&v).On("CONFLICT (ulid) DO UPDATE").Exec(context)
		if err != nil {
			return err
		}
	}

	err= tx.Commit()
	if err != nil {
		return err
	}
	return nil
}


func(this *PushDataService) PushPOSOrderDetailPackage(context context.Context, data_order_detail_package PosOrderDetailPackageDTO)error{
	tx, err := this.DB.BeginTx(context, nil)
	if err != nil {
		return err
	}

	for _,v := range data_order_detail_package.ListOrderDetailPackage {
		_,err= tx.NewInsert().Model(&v).On("CONFLICT (ulid) DO UPDATE").Exec(context)
		if err != nil {
			return err
		}
	}

	err= tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func(this *PushDataService) PushPOSOrderPayment(context context.Context, data_order_payment PosOrderPaymentDTO)error{
	tx, err := this.DB.BeginTx(context, nil)
	if err != nil {
		return err
	}

	var list_payment_number  []string

	for _,v := range data_order_payment.ListOrderPayment{
		list_payment_number = append(list_payment_number, v.PaymentNumber)
	}
	var list_payment_number_merger []string

	seen := make(map[string]struct{})

	for _, item := range list_payment_number {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			list_payment_number_merger = append(list_payment_number_merger, item)
		}
	}

	_,err = tx.NewRaw("delete from pos_order_payment where payment_number IN (?)", bun.In(list_payment_number)).Exec(context)
	if err != nil {
		return err
	}

	for _,v := range data_order_payment.ListOrderPayment {
		_,err= tx.NewInsert().Model(&v).On("CONFLICT (ulid) DO UPDATE").Exec(context)
		if err != nil {
			return err
		}
	}

	err= tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
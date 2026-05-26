package setup

import (
	"APIANDORDER/backend/helpers"
	"context"

	"github.com/uptrace/bun"
)

type SetupService struct {
	DB *bun.DB
}

func NewSetupService(db *bun.DB) *SetupService {
	return &SetupService{DB: db}
}

func (this *SetupService) CekLogin(context context.Context, username, password string) (bool, int, error) {

	datauser := DataUser{}

	err := this.DB.NewRaw("SELECT id, username, password from master_users where username = ? and status='1'", username).Scan(context, &datauser)

	if err != nil {
		return false, 0, err
	}

	if datauser.Id == 0 {
		return false, 0, err
	}

	if !helpers.CompareHashAndPassword(datauser.Password, password) {
		return false, 0, err
	}
	return true, datauser.Id, err
}



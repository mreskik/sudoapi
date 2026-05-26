package setup

import (
	"APIANDORDER/backend/config"
	"APIANDORDER/tester"
	"context"
	"testing"
)

func TestSetupService(t *testing.T) {

	testere := tester.New()
	defer func() {
		testere.Close()
	}()

	t.Run("Testing GetDataBranch()", func(t *testing.T) {
		// t.Log("============== testing GetDataBranch =============")
		setupservice := NewSetupService(config.DB)
		data, _, err := setupservice.CekLogin(context.Background(), "user", "user")
		if err != nil {
			t.Errorf("======= ERROR TEST : %s", err)
			return
		}
		t.Log(data)
	})

}

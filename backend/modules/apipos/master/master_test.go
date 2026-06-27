package master

import (
	"APIANDORDER/tester"
	"testing"
)

func TestMasterService(t *testing.T) {
	testere := tester.New()
	defer func() {
		testere.Close()
	}()

	t.Run("TESTING GET ROLE ACCESS BY BRANCH", func(t *testing.T) {
		categoryservices := NewMasterService(testere.DB)

		data, err := categoryservices.GetMasterRoleAccess(t.Context(), 14)

		if err != nil {
			t.Errorf("error : %s", err)
			return
		}
		t.Log(data)
	})
}

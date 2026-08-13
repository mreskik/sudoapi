package endday

import (
	"APIANDORDER/backend/helpers"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type EnddayHandler struct{}

func NewHandler() *EnddayHandler {
	return &EnddayHandler{}
}

// sudocoreResponse bentuk response dari sudocore2 endpoint /api/pos/endday-jurnal &
// /api/pos/revert-jurnal (masih gaya lama {status, message}, bukan {code, message, data}).
type sudocoreResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func sudocoreBaseURL() string {
	return "http://" + os.Getenv("APP_SUDOCORE_HOST") + ":" + os.Getenv("APP_SUDOCORE_PORT")
}

// forwardToSudocore terusin token yang sama yang divalidasi middleware.BranchTokenAuth di
// sini ke sudocore2 -- sudocore2 validasi ulang token itu sendiri (defense in depth, endpoint
// sudocore2 ini bisa diakses langsung, bukan cuma lewat APIANDORDER).
func forwardToSudocore(url string, token string) (*sudocoreResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	httpRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	parsed := sudocoreResponse{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func tokenFromHeader(c *gin.Context) string {
	return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
}

// RequestEndDay minta sudocore2 jalankan jurnal endday buat 1 dayshift -- dipanggil POS pas
// dayout (lewat SyncController/DayShiftServices di Laravel).
func (h *EnddayHandler) RequestEndDay(c *gin.Context) {
	res := helpers.NewResponse()

	branchID := c.Param("branch_id")
	dayshiftUlid := c.Param("dayshift_ulid")

	url := sudocoreBaseURL() + "/api/pos/endday-jurnal/" + branchID + "/" + dayshiftUlid
	parsed, err := forwardToSudocore(url, tokenFromHeader(c))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage(err.Error()))
		return
	}
	if parsed.Status != 0 {
		c.JSON(200, res.GeneralError().SetMessage(parsed.Message))
		return
	}
	c.JSON(200, res.Success().SetMessage(parsed.Message))
}

// RequestEndDayRevert minta sudocore2 batalkan jurnal endday buat 1 dayshift -- aksi manual
// (belum ada auto-trigger dari POS), tetap lewat jalur token yang sama.
func (h *EnddayHandler) RequestEndDayRevert(c *gin.Context) {
	res := helpers.NewResponse()

	branchID := c.Param("branch_id")
	dayshiftUlid := c.Param("dayshift_ulid")

	url := sudocoreBaseURL() + "/api/pos/revert-jurnal/" + branchID + "/" + dayshiftUlid
	parsed, err := forwardToSudocore(url, tokenFromHeader(c))
	if err != nil {
		c.JSON(200, res.GeneralError().SetMessage(err.Error()))
		return
	}
	if parsed.Status != 0 {
		c.JSON(200, res.GeneralError().SetMessage(parsed.Message))
		return
	}
	c.JSON(200, res.Success().SetMessage(parsed.Message))
}

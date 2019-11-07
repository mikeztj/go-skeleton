package skeleton

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	skeletonEntity "github.com/vilbert/ProdukOnline/internal/entity/skeleton"
	"github.com/vilbert/ProdukOnline/pkg/response"
)

// ISkeletonSvc is an interface to Skeleton Service
type ISkeletonSvc interface {
	GetAllBridgingProduct(ctx context.Context, page int, length int) ([]skeletonEntity.Skeleton, interface{}, error)
}

type (
	// Handler ...
	Handler struct {
		skeletonSvc ISkeletonSvc
	}

	skeletonResp struct {
		BPCode            string  `json:"bp_code"`
		BPName            string  `json:"bp_name"`
		BPName2           string  `json:"bp_name2"`
		BPDeptCode        string  `json:"bp_deptcode"`
		BPSource          string  `json:"bp_source"`
		BPActiveYN        string  `json:"bp_activeyn"`
		BPUserID          string  `json:"bp_userid"`
		BPLastUpdate      string  `json:"bp_lastupdate"`
		BPDataAktifYN     string  `json:"bp_dataaktifyn"`
		BPArsipRecordDate string  `json:"bp_arsiprecorddate"`
		BPPurgeRecordDate *string `json:"bp_purgerecorddate"`
	}
)

// New for bridging product handler initialization
func New(is ISkeletonSvc) *Handler {
	return &Handler{
		skeletonSvc: is,
	}
}

// SkeletonHandler return user data
func (h *Handler) SkeletonHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp     *response.Response
		page     int
		length   int
		result   interface{}
		metadata interface{}
		err      error
		errRes   response.Error
	)
	resp = &response.Response{}
	defer resp.RenderJSON(w, r)

	// Check if request method is GET
	if r.Method == http.MethodGet {
		// Check if page and length exists in URL parameters
		_, pageOK := r.URL.Query()["page"]
		_, lengthOK := r.URL.Query()["length"]

		// If page and length exists
		if pageOK && lengthOK {
			// Get page and length data and convert to reassure data type is int
			page, _ = strconv.Atoi(r.FormValue("page"))
			length, _ = strconv.Atoi(r.FormValue("length"))
		} else {
			page, length = 0, 0
		}

		result, metadata, err = h.skeletonSvc.GetAllBridgingProduct(context.Background(), page, length)
	}

	// If anything from service or data return an error
	if err != nil {
		// Error response handling
		errRes = response.Error{
			Code:   101,
			Msg:    "Data Not Found",
			Status: true,
		}
		// If service returns an error
		if strings.Contains(err.Error(), "service") {
			// Replace error with server error
			errRes = response.Error{
				Code:   201,
				Msg:    "Failed to process request due to server error",
				Status: true,
			}
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.Error = errRes
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
}

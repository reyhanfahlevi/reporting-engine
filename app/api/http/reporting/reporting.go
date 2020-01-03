package reporting

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tokopedia/td-report-engine/app/api"
	"github.com/tokopedia/td-report-engine/internal/usecase/reporting"
)

type errorResponse struct {
	ErrorMsg string `json:"error_msg"`
}

type defaultResponse struct {
	Success bool `json:"success"`
}

// Handler reporting handler usecase
type Handler struct {
	reporting api.ReportingUseCase
}

// New will instantiate the reporting usecase handler
func New(ucReporting api.ReportingUseCase) *Handler {
	return &Handler{
		reporting: ucReporting,
	}
}

// HandlerStoreReport func
func (h *Handler) HandlerStoreReport(r gin.IRoutes) gin.IRoutes {
	return r.POST("/store", func(c *gin.Context) {
		req := make(map[string]interface{})
		data := make(map[string]interface{})

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{ErrorMsg: err.Error()})
			return
		}

		data, ok := req["data"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, errorResponse{ErrorMsg: "failed parsing data data"})
			return
		}

		reportType, ok := req["report_type"]
		if !ok {
			c.JSON(http.StatusBadRequest, errorResponse{ErrorMsg: "report type is required"})
			return
		}

		err = h.reporting.SaveReport(c, reporting.ParamSaveReport{
			ReportType: reportType.(string),
			Data:       data,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{ErrorMsg: err.Error()})
			return
		}

		c.JSON(http.StatusOK, defaultResponse{Success: true})
	})
}

// HandlerGetReport func
func (h *Handler) HandlerGetReport(r gin.IRoutes) gin.IRoutes {
	return r.GET("/list", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		param := reporting.ParamGetReports{
			ReportType: c.Query("report_type"),
			Page:       page,
			Limit:      limit,
		}

		resp, err := h.reporting.GetReports(c, param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{ErrorMsg: err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})
}

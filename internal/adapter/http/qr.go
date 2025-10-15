package httpadapter

import (
	"context"
	"strings"

	"strconv"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type QrService interface {
	ExtractContentsToDB(ctx context.Context, inventoryData models.InventoryData) error
}

type QrHandler struct {
	qrService QrService
}

func NewQrHandler(qrService QrService) *QrHandler {
	return &QrHandler{
		qrService: qrService,
	}
}

func (q *QrHandler) ExtractInventory(c *gin.Context) {

	item := c.Query("item")
	quantity, _ := strconv.ParseFloat(strings.TrimSpace(c.Query("quantity")), 32)
	unit := c.Query("unit")
	wholesale_price_per_quantity, _ := strconv.ParseFloat(strings.TrimSpace(c.Query("wholesale_price_per_quantity")), 32)
	total_cost_of_product, err := strconv.ParseFloat(strings.TrimSpace(c.Query("total_cost_of_product")), 32)

	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "Bad request, the request is invalid",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	inventoryDataBody := models.InventoryData{
		Item:                      item,
		Quantity:                  quantity,
		Unit:                      unit,
		WholesalePricePerQuantity: wholesale_price_per_quantity,
		TotalCostOfProduct:        total_cost_of_product,
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusOK,
		Message: "successfully extracted inventory data from QR code",
		Data:    inventoryDataBody,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)

}

func (q *QrHandler) Extract(c *gin.Context) {

	ctx := c.Request.Context()
	var inventoryDataBody models.InventoryData
	if err := c.BindJSON(&inventoryDataBody); err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "Bad request, the request is invalid",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}
	err := q.qrService.ExtractContentsToDB(ctx, inventoryDataBody)

	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "Unable to store QR data to Database",
			Data:    err.Error(),
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusCreated,
		Message: "Stored inventory data to Database",
		Data:    inventoryDataBody,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}

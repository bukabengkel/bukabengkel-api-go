package response

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
)

type productExportLogResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	File      string    `json:"file"`
	DoneAt    time.Time `json:"doneAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ProductExportLogsListResponse(products *[]models.ProductExportLog) *[]productExportLogResponse {
	var responses = make([]productExportLogResponse, 0)
	for _, product := range *products {
		response := &productExportLogResponse{
			Status:    product.Status.String(),
			File:      product.PathFile,
			DoneAt:    product.DoneAt.Local(),
			CreatedAt: product.CreatedAt.Local(),
			UpdatedAt: product.UpdatedAt.Local(),
		}
		responses = append(responses, *response)
	}

	return &responses
}

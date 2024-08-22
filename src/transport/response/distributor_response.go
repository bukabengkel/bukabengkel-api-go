package response

import (
	"github.com/peang/bukabengkel-api-go/src/models"
)

type distributorDetailResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DistributorDetailResponse(distributor *models.Distributor) *distributorDetailResponse {
	response := &distributorDetailResponse{
		ID:   distributor.Key,
		Name: distributor.Name,
	}

	return response
}

func DistributorListResponse(distributors *[]models.Distributor) []distributorDetailResponse {
	var responses = make([]distributorDetailResponse, 0)
	for _, product := range *distributors {
		response := DistributorDetailResponse(&product)
		responses = append(responses, *response)
	}
	return responses
}

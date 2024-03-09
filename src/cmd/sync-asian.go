package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/spf13/cobra"
)

type CategoryResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Content struct {
		Data []struct {
			ID    int    `json:"id"`
			Code  string `json:"code"`
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"data"`
	} `json:"content"`
}

type ProductListResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Content struct {
		Data []Product `json:"data"`
		Meta Meta      `json:"meta"`
	} `json:"content"`
}

type Product struct {
	ID             int     `json:"id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Qty            int     `json:"qty"`
	AvailableQty   int     `json:"available_qty"`
	BasePrice      int     `json:"base_price"`
	Unit           string  `json:"unit"`
	Images         string  `json:"images"`
	Price          []Price `json:"price"`
	WishlistStatus string  `json:"wishlist_status"`
	Catcode        string  `json:"catcode"`
	RankingStatus  string  `json:"ranking_status"`
	Volume         string  `json:"volume"`
	Weight         string  `json:"weight"`
}

type Price struct {
	Qty   int `json:"qty"`
	Price int `json:"price"`
}

type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type SyncAsian struct {
	productDistributorRepo         *repository.ProductDistributorRepository
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository
}

func NewSyncAsian(
	productDistributorRepo *repository.ProductDistributorRepository,
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository,
) *SyncAsian {
	return &SyncAsian{
		productDistributorRepo:         productDistributorRepo,
		productCategoryDistributorRepo: productCategoryDistributorRepo,
	}
}

func (s *SyncAsian) Execute(cmd *cobra.Command, args []string) {
	// s.productCategoryDistributorRepo.UpdateWithCondition(
	// 	repository.ProductCategoryDistributorRepositoryFilter{
	// 		DistributorID: utils.Uint64(1),
	// 	},
	// 	repository.ProductCategoryDistributorRepositoryValues{
	// 		RemoteUpdate: utils.Boolean(false),
	// 	},
	// )

	// s.syncCategory(1)
	// s.syncCategory(2)

	s.syncProduct(1)
}

func (s *SyncAsian) syncCategory(cat uint) {
	var errorCount uint
	resp, err := utils.HttpGetWithRetry(fmt.Sprintf("https://api-mobile.asian-accessory.com/category/list/%v", cat), "GET", 5)
	if err != nil {
		log.Fatalf("Failed to Fetch %v", err)
	}

	var response CategoryResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		log.Fatalf("Failed to parse %v", err)
	}

	for _, category := range response.Content.Data {
		pc, err := s.productCategoryDistributorRepo.FindOne(repository.ProductCategoryDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
			Code:          &category.Code,
		})
		if err != nil {
			log.Fatal(err)
			fmt.Printf("Skipping %v", category.Name)
			errorCount++
			continue
		}

		if pc == nil {
			newPc := models.ProductCategoryDistributor{
				ExternalID:    strconv.Itoa(category.ID),
				DistributorID: 1,
				Name:          category.Name,
				Code:          category.Code,
				Description:   "",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				RemoteUpdate:  true,
			}

			s.productCategoryDistributorRepo.Save(&newPc)
		} else {
			pc.Name = category.Name
			pc.RemoteUpdate = true

			s.productCategoryDistributorRepo.Update(pc)
		}
	}
}

func (s *SyncAsian) syncProduct(cat uint) {
	var errorCount uint
	var page uint
	hasResult := true
	for hasResult {
		resp, err := utils.HttpGetWithRetry(fmt.Sprintf("https://api-mobile.asian-accessory.com/category/product/%v?page=%v", cat, page), "GET", 5)
		if err != nil {
			log.Fatalf("Failed to Fetch %v", err)
		}

		var response ProductListResponse
		err = json.Unmarshal([]byte(resp), &response)
		if err != nil {
			log.Fatalf("Failed to parse %v", err)
		}

		if len(response.Content.Data) == 0 {
			hasResult = false
			break
		}

		for _, product := range response.Content.Data {
			p, err := s.productDistributorRepo.FindOne(repository.ProductDistributorRepositoryFilter{
				DistributorID: utils.Uint64(1),
				Code:          &product.Code,
			})

			if err != nil {
				log.Fatal(err)
				fmt.Printf("Skipping %v", product.Name)
				errorCount++
				continue
			}

			if p == nil {
				cat, err := s.productCategoryDistributorRepo.FindOne(repository.ProductCategoryDistributorRepositoryFilter{
					DistributorID: utils.Uint64(1),
					Code:          &product.Catcode,
				})

				if err != nil {
					log.Fatal(err)
					fmt.Printf("Skipping %v", product.Name)
					errorCount++
					continue
				}

				weight, _ := strconv.Atoi(product.Weight)
				volume, _ := strconv.Atoi(product.Volume)

				newProductDistributor := models.ProductDistributor{
					ExternalID:       strconv.Itoa(int(product.ID)),
					Key:              uuid.NewString(),
					DistributorID:    1,
					CategoryID:       *cat.ID,
					Name:             product.Name,
					Code:             product.Code,
					Description:      "",
					Unit:             product.Unit,
					Thumbnail:        "",
					Images:           nil,
					Price:            float64(product.BasePrice),
					BulkPrice:        nil,
					Weight:           float64(weight),
					Volume:           float64(volume),
					Stock:            float64(product.Qty),
					IsStockUnlimited: false,
					Status:           models.ProductActive,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
					RemoteUpdate:     true,
				}

				s.productDistributorRepo.Save(&newProductDistributor)
			} else {
				continue
			}
		}

		page++
	}
}

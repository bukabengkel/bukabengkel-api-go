package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
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
	imageRepo                      *repository.ImageRepository
	s3service                      *file_service.S3Service
}

func NewSyncAsian(
	productDistributorRepo *repository.ProductDistributorRepository,
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository,
	imageRepo *repository.ImageRepository,
	s3service *file_service.S3Service,
) *SyncAsian {
	return &SyncAsian{
		productDistributorRepo:         productDistributorRepo,
		productCategoryDistributorRepo: productCategoryDistributorRepo,
		imageRepo:                      imageRepo,
		s3service:                      s3service,
	}
}

func (s *SyncAsian) Execute(cmd *cobra.Command, args []string) {
	s.productCategoryDistributorRepo.UpdateWithCondition(
		repository.ProductCategoryDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
		},
		repository.ProductCategoryDistributorRepositoryValues{
			RemoteUpdate: utils.Boolean(false),
		},
	)

	s.productDistributorRepo.UpdateWithCondition(
		repository.ProductDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
		},
		repository.ProductDistributorRepositoryValues{
			RemoteUpdate: utils.Boolean(false),
		},
	)

	s.syncCategory(1)
	s.syncCategory(2)

	s.syncProduct(1)
	s.syncProduct(2)
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
		resp, err := utils.HttpGetWithRetry(fmt.Sprintf("https://api-mobile.asian-accessory.com/category/product/%v?page=%v&per_page=50", cat, page), "GET", 5)
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

		var wg sync.WaitGroup
		wg.Add(len(response.Content.Data))

		for _, product := range response.Content.Data {
			go func(product Product) {
				defer wg.Done()

				p, err := s.productDistributorRepo.FindOne(repository.ProductDistributorRepositoryFilter{
					DistributorID: utils.Uint64(1),
					Code:          &product.Code,
				})

				if err != nil {
					log.Fatal(err)
					fmt.Printf("Skipping %v, Error Getting Product", product.Name)
					errorCount++
					return
				}

				if p == nil {
					cat, err := s.productCategoryDistributorRepo.FindOne(repository.ProductCategoryDistributorRepositoryFilter{
						DistributorID: utils.Uint64(1),
						Code:          &product.Catcode,
					})

					if cat == nil {
						log.Fatal(err)
						fmt.Printf("Skipping %v, Unknown Product Category", product.Catcode)
						errorCount++
						return
					}

					if err != nil {
						log.Fatal(err)
						fmt.Printf("Skipping %v, Error Getting Product Category", product.Name)
						errorCount++
						return
					}

					weight, _ := strconv.ParseFloat(product.Weight, 64)
					volume, _ := strconv.ParseFloat(product.Volume, 64)
					img, err := s.s3service.Upload(models.ImageProductCategory, product.Images)
					if err != nil {
						log.Fatal(err)
						fmt.Printf("Skipping %v, Error Uploading to S3", product.Name)
						errorCount++
						return
					}

					var bulkPrice []models.ProductBulkPrice
					for _, price := range product.Price {
						bulkPrice = append(bulkPrice, models.ProductBulkPrice{
							Qty:   int64(price.Qty),
							Price: float32(price.Price),
						})
					}

					newProductDistributor := models.ProductDistributor{
						ExternalID:       strconv.Itoa(int(product.ID)),
						Key:              uuid.NewString(),
						DistributorID:    1,
						CategoryID:       *cat.ID,
						Name:             product.Name,
						Code:             product.Code,
						Description:      "",
						Unit:             product.Unit,
						Thumbnail:        img.Key,
						Images:           []string{img.Key},
						Price:            float64(product.BasePrice),
						BulkPrice:        bulkPrice,
						Weight:           float64(weight),
						Volume:           float64(volume),
						Stock:            float64(product.AvailableQty),
						IsStockUnlimited: false,
						Status:           models.ProductActive,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
						RemoteUpdate:     true,
					}

					s.productDistributorRepo.Save(&newProductDistributor)
				} else {
					p.Name = product.Name
					p.Code = product.Code
					p.Unit = product.Unit
					p.Stock = float64(product.AvailableQty)
					p.Price = float64(product.BasePrice)
					p.RemoteUpdate = true

					s.productDistributorRepo.Update(p)
				}
			}(product)
		}

		wg.Wait()
		time.Sleep(1 * time.Second)
		page++
	}
}

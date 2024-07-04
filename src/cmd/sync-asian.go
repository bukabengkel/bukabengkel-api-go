package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		Data []CategoryResponseData `json:"data"`
	} `json:"content"`
}

type CategoryResponseData struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ProductListResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Content struct {
		Data []ProductResponseData `json:"data"`
		Meta Meta                  `json:"meta"`
	} `json:"content"`
}

type ProductResponseData struct {
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
	logger                         utils.Logger
	productDistributorRepo         *repository.ProductDistributorRepository
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository
	imageRepo                      *repository.ImageRepository
	fileService                    file_service.FileService
}

func NewSyncAsian(
	logger utils.Logger,
	productDistributorRepo *repository.ProductDistributorRepository,
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository,
	imageRepo *repository.ImageRepository,
	fileService file_service.FileService,
) *SyncAsian {
	return &SyncAsian{
		logger:                         logger,
		productDistributorRepo:         productDistributorRepo,
		productCategoryDistributorRepo: productCategoryDistributorRepo,
		imageRepo:                      imageRepo,
		fileService:                    fileService,
	}
}

func (s *SyncAsian) ExecuteViaCmd(cmd *cobra.Command, args []string) {
	s.Execute()
}

func (s *SyncAsian) Execute() {
	file, err := os.Create("./src/cmd/logs/sync-asian-error.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

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

	var wg sync.WaitGroup
	var mt sync.Mutex
	var createdCategory, updatedCategory, createdProduct, updatedProduct int

	errorChannel := make(chan error)
	wg.Add(1)
	go s.syncProductErrorWriter(&wg, errorChannel)

	// Category Sync
	categoryChannel := make(chan CategoryResponseData)
	categoryWorkerNum := 2
	for i := 0; i < categoryWorkerNum; i++ {
		wg.Add(1)
		go s.syncCategory(&wg, &mt, categoryChannel, errorChannel, &createdCategory, &updatedCategory)
	}

	// Product Sync
	productChannel := make(chan ProductResponseData)
	productWorkerNum := 50
	for i := 0; i < productWorkerNum; i++ {
		wg.Add(1)
		go s.syncProduct(
			&wg,
			&mt,
			productChannel,
			errorChannel,
			&createdProduct,
			&updatedProduct,
		)
	}

	s.getCategory(1, categoryChannel, errorChannel)
	s.getCategory(2, categoryChannel, errorChannel)
	fmt.Println("")
	close(categoryChannel)

	s.getProduct(1, productChannel, errorChannel)
	s.getProduct(2, productChannel, errorChannel)
	fmt.Println("")
	close(productChannel)

	close(errorChannel)
	wg.Wait()

	s.remove()

}

func (s *SyncAsian) getCategory(
	cat uint,
	ch chan<- CategoryResponseData,
	chErr chan<- error,
) {
	resp, err := utils.HttpGetWithRetry(fmt.Sprintf("https://api-mobile.asian-accessory.com/category/list/%v", cat), "GET", 5)
	if err != nil {
		chErr <- fmt.Errorf("failed_to_fetch %v", err)
		return
	}

	var response CategoryResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		chErr <- fmt.Errorf("failed_to_parse %v", err)
		return
	}

	for _, category := range response.Content.Data {
		ch <- category
	}
}

func (s *SyncAsian) syncCategory(
	wg *sync.WaitGroup,
	mt *sync.Mutex,
	ch <-chan CategoryResponseData,
	chErr chan<- error,
	totalCreated *int,
	totalUpdated *int,
) {
	defer wg.Done()

	for category := range ch {
		fmt.Printf("\rCategory Created %v Category Updated %v", *totalCreated, *totalUpdated)

		pc, err := s.productCategoryDistributorRepo.FindOne(repository.ProductCategoryDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
			Code:          &category.Code,
		})
		if err != nil {
			chErr <- fmt.Errorf("error_find_category %v", err)
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

			mt.Lock()
			*totalCreated++
			mt.Unlock()
		} else {
			pc.Name = category.Name
			pc.RemoteUpdate = true

			s.productCategoryDistributorRepo.Update(pc)

			mt.Lock()
			*totalUpdated++
			mt.Unlock()
		}
	}
}

func (s *SyncAsian) getProduct(
	cat uint,
	ch chan<- ProductResponseData,
	chErr chan<- error,
) {
	var errorCount uint

	page := 1
	hasResult := true
	for hasResult {
		var errorCountPerPage uint

		resp, err := utils.HttpGetWithRetry(fmt.Sprintf("https://api-mobile.asian-accessory.com/category/product/%v?page=%v&per_page=200&sort=Terbaru", cat, page), "GET", 5)
		if err != nil {
			errNew := fmt.Errorf("failed_to_fetch %v", err)
			chErr <- errNew
			return
		}

		var response ProductListResponse
		err = json.Unmarshal([]byte(resp), &response)
		if err != nil {
			errNew := fmt.Errorf("failed_to_parse %v", err)
			chErr <- errNew
			return
		}

		if len(response.Content.Data) == 0 {
			hasResult = false
			break
		}

		for _, product := range response.Content.Data {
			ch <- product
		}

		page++
		errorCount += errorCountPerPage
	}
}

func (s *SyncAsian) syncProduct(
	wg *sync.WaitGroup,
	mt *sync.Mutex,
	ch <-chan ProductResponseData,
	chErr chan<- error,
	totalCreated *int,
	totalUpdated *int,
) {
	defer wg.Done()

	for product := range ch {
		fmt.Printf("\rProduct Created %v Product Updated %v", *totalCreated, *totalUpdated)

		p, err := s.productDistributorRepo.FindOne(repository.ProductDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
			Code:          &product.Code,
		})
		if err != nil {
			newErr := fmt.Errorf("error_getting_product;%s;%v", product.Catcode, err.Error())
			chErr <- newErr
			return
		}

		if p == nil {
			cat, err := s.productCategoryDistributorRepo.FindOne(repository.ProductCategoryDistributorRepositoryFilter{
				DistributorID: utils.Uint64(1),
				Code:          &product.Catcode,
			})

			if cat == nil {
				newErr := fmt.Errorf("unknown_product_category;%s;%v", product.Catcode, err.Error())
				chErr <- newErr
				return
			}

			if err != nil {
				newErr := fmt.Errorf("error_getting_product_category;%s;%v", product.Name, err.Error())
				chErr <- newErr
				return
			}

			weight, _ := strconv.ParseFloat(product.Weight, 64)
			volume, _ := strconv.ParseFloat(product.Volume, 64)
			var img file_service.FileUploadResponse
			if product.Images != "" {
				imgPtr, err := s.fileService.Upload(models.ImageProductDistributor, product.Images)
				if err != nil {
					newErr := fmt.Errorf("error_uploading_image;%s;%v", product.Name, err.Error())
					chErr <- newErr
					return
				}
				img = *imgPtr
			} else {
				img.Key = ""
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

			_, err = s.productDistributorRepo.Save(&newProductDistributor)
			if err != nil {
				newErr := fmt.Errorf("error_inserting_product;%s;%v", product.Name, err.Error())
				chErr <- newErr
				return
			}

			mt.Lock()
			*totalCreated++
			mt.Unlock()
		} else {
			p.Name = product.Name
			p.Unit = product.Unit
			p.Stock = float64(product.AvailableQty)
			p.Price = float64(product.BasePrice)
			p.RemoteUpdate = true
			p.UpdatedAt = time.Now()

			_, err := s.productDistributorRepo.Update(p)
			if err != nil {
				newErr := fmt.Errorf("error_updating_product;%s;%v", product.Name, err.Error())
				chErr <- newErr
				return
			}

			mt.Lock()
			*totalUpdated++
			mt.Unlock()
		}
	}
}

func (s *SyncAsian) syncProductErrorWriter(wg *sync.WaitGroup, ch <-chan error) {
	defer wg.Done()

	file, err := os.OpenFile("./src/cmd/logs/sync-asian-error.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for errorStr := range ch {
		_, err = file.WriteString(errorStr.Error() + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
func (s *SyncAsian) remove() {
	hasResult := true

	for hasResult {
		products, _, err := s.productDistributorRepo.List(context.TODO(), 1, 10, "id", repository.ProductDistributorRepositoryFilter{
			DistributorID: utils.Uint64(1),
			RemoteUpdate:  utils.Boolean(false),
		})

		if err != nil {
			log.Fatal(err)
		}

		if len(*products) == 0 {
			hasResult = false
			break
		}

		for _, product := range *products {
			if len(product.Images) > 0 {
				s.fileService.Delete(product.Images[0])
			}
			s.productDistributorRepo.Delete(&product)
		}
	}
}

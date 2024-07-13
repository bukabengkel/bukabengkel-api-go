package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/xuri/excelize/v2"
)

func ExportProduct(
	productRepository *repository.ProductRepository,
	productExportLogRepository *repository.ProductExportLogRepository,
	productLogId uint64,
	storeId uint64,
	categoryId *uint64,
) {
	filename := uuid.New().String()
	ch := make(chan models.Product, 10)
	wg := sync.WaitGroup{}

	excelFile := excelize.NewFile()
	sheetName := "Sheet1"
	excelFile.NewSheet(sheetName)

	headerStyle, err := excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F9CB9C"},
			Pattern: 1,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create style: %v", err)
	}

	headerErrStyle, err := excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFC7CE"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: "#FF0000",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create style: %v", err)
	}

	excelFile.SetCellStyle(sheetName, "A1", "Q1", headerStyle)
	excelFile.SetCellStyle(sheetName, "R1", "R1", headerErrStyle)
	excelFile.SetColWidth(sheetName, "A", "R", 35)

	excelFile.SetCellValue(sheetName, "A1", "ID Produk")
	excelFile.SetCellValue(sheetName, "B1", "Nama Produk")
	excelFile.SetCellValue(sheetName, "C1", "Merk")
	excelFile.SetCellValue(sheetName, "D1", "Kategori")
	excelFile.SetCellValue(sheetName, "E1", "Unit")
	excelFile.SetCellValue(sheetName, "F1", "Deskripsi")
	excelFile.SetCellValue(sheetName, "G1", "Harga Beli")
	excelFile.SetCellValue(sheetName, "H1", "Harga Jual")
	excelFile.SetCellValue(sheetName, "I1", "Produk Menggunakan Stok")
	excelFile.SetCellValue(sheetName, "J1", "Stok")
	excelFile.SetCellValue(sheetName, "K1", "Stok Minimal")
	excelFile.SetCellValue(sheetName, "L1", "Status")
	excelFile.SetCellValue(sheetName, "M1", "Gambar 1")
	excelFile.SetCellValue(sheetName, "N1", "Gambar 2")
	excelFile.SetCellValue(sheetName, "O1", "Gambar 3")
	excelFile.SetCellValue(sheetName, "P1", "Gambar 4")
	excelFile.SetCellValue(sheetName, "Q1", "Gambar 5")
	excelFile.SetCellValue(sheetName, "R1", "Error")

	wg.Add(1)
	go toExcel(&wg, ch, excelFile, sheetName)

	filter := repository.ProductRepositoryFilter{
		StoreID: &storeId,
	}
	if categoryId != nil {
		filter.CategoryId = categoryId
	}

	var done bool
	for !done {
		products, count, _ := productRepository.List(context.TODO(), 1, 20, "id", filter)
		if count == 0 {
			done = true
			break
		}

		for _, product := range *products {
			ch <- product
		}

		if count < 20 {
			done = true
			break
		}
	}
	close(ch)

	wg.Wait()

	if err := excelFile.SaveAs(fmt.Sprintf("tmp/%s.xlsx", filename)); err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}
}

func toExcel(wg *sync.WaitGroup, ch <-chan models.Product, exc *excelize.File, sheetName string) {
	defer wg.Done()

	cellNumber := 2
	for product := range ch {
		exc.SetCellValue(sheetName, fmt.Sprintf("A%d", cellNumber), product.Key)
		exc.SetCellValue(sheetName, fmt.Sprintf("B%d", cellNumber), product.Name)

		if product.Brand != nil {
			exc.SetCellValue(sheetName, fmt.Sprintf("C%d", cellNumber), product.Brand.Name)
		}

		exc.SetCellValue(sheetName, fmt.Sprintf("D%d", cellNumber), product.Category.Name)
		exc.SetCellValue(sheetName, fmt.Sprintf("E%d", cellNumber), product.Unit)
		exc.SetCellValue(sheetName, fmt.Sprintf("F%d", cellNumber), product.Description)
		exc.SetCellValue(sheetName, fmt.Sprintf("G%d", cellNumber), product.Price)
		exc.SetCellValue(sheetName, fmt.Sprintf("H%d", cellNumber), product.SellPrice)
		exc.SetCellValue(sheetName, fmt.Sprintf("I%d", cellNumber), product.IsStockUnlimited)
		exc.SetCellValue(sheetName, fmt.Sprintf("J%d", cellNumber), product.Stock)
		exc.SetCellValue(sheetName, fmt.Sprintf("K%d", cellNumber), product.StockMinimum)
		exc.SetCellValue(sheetName, fmt.Sprintf("L%d", cellNumber), product.Status.String())

		if len(product.Images) > 1 && product.Images[0].Path != "" {
			exc.SetCellValue(sheetName, fmt.Sprintf("M%d", cellNumber), product.Images[0].Path)
		}

		if len(product.Images) > 2 && product.Images[1].Path != "" {
			exc.SetCellValue(sheetName, fmt.Sprintf("M%d", cellNumber), product.Images[1].Path)
		}

		if len(product.Images) > 3 && product.Images[2].Path != "" {
			exc.SetCellValue(sheetName, fmt.Sprintf("M%d", cellNumber), product.Images[2].Path)
		}

		if len(product.Images) > 4 && product.Images[3].Path != "" {
			exc.SetCellValue(sheetName, fmt.Sprintf("M%d", cellNumber), product.Images[3].Path)
		}

		if len(product.Images) > 5 && product.Images[4].Path != "" {
			exc.SetCellValue(sheetName, fmt.Sprintf("M%d", cellNumber), product.Images[4].Path)
		}

		// // Set the active sheet
		// exc.SetActiveSheet(index)
		cellNumber++
	}
}

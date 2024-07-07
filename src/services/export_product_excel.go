package services

import (
	"context"
	"fmt"
	"log"
	"sync"

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

	if err := excelFile.SaveAs("test.xlsx"); err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}
}

func toExcel(wg *sync.WaitGroup, ch <-chan models.Product, exc *excelize.File, sheetName string) {
	defer wg.Done()

	cellNumber := 2
	for product := range ch {
		fmt.Println(product.Name)

		// Set values in the sheet
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
		// f.SetCellValue(sheetName, "A1", "Name")
		// f.SetCellValue(sheetName, "B1", "Age")
		// f.SetCellValue(sheetName, "A2", "John Doe")
		// f.SetCellValue(sheetName, "B2", 29)
		// f.SetCellValue(sheetName, "A3", "Jane Smith")
		// f.SetCellValue(sheetName, "B3", 34)

		// // Set the active sheet
		// f.SetActiveSheet(index)
		cellNumber++
	}
}

package cmd

import (
	"fmt"

	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/spf13/cobra"
)

type SyncAsian struct {
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository
}

func NewSyncAsian(
	productCategoryDistributorRepo *repository.ProductCategoryDistributorRepository,
) *SyncAsian {
	return &SyncAsian{
		productCategoryDistributorRepo: productCategoryDistributorRepo,
	}
}

func (s *SyncAsian) Execute(cmd *cobra.Command, args []string) {
	fmt.Println("Sync Asian Update II")
}

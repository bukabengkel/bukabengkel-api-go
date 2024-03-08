package cmd

import (
	"fmt"

	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/spf13/cobra"
)

type SyncAsian struct {
	productRepo *repository.ProductRepository
}

func NewSyncAsian(
	productRepo *repository.ProductRepository,
) *SyncAsian {
	return &SyncAsian{
		productRepo: productRepo,
	}
}

func (s *SyncAsian) Execute(cmd *cobra.Command, args []string) {
	fmt.Println("Sync Asian Update II")
}

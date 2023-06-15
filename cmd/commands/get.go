package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get curency snapshot",
		Long:  `Get curency snapshot`,
		Run: func(cmd *cobra.Command, args []string) {
			GetCurrencySnapshot()
		},
	}
	input string
)

const (
	DBName = "valueshift"
)

func init() {
	getCmd.Flags().StringVarP(&input, "input", "i", "USD", "Input currency lable")

	rootCmd.AddCommand(getCmd)
}

func GetCurrencySnapshot() {
	result, err := currencyService.GetCurrencySnapshot(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

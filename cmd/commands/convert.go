package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert value -v currency -i to -o",
		Long:  `Convert value flag=v to currency input flag=i and output flag=o.`,
		Run: func(cmd *cobra.Command, args []string) {
			GetCurrencySnapshot()
		},
	}
	input  string
	output string
	value  float64
)

const (
	DBName = "valueshift"
)

func init() {
	convertCmd.Flags().StringVarP(&input, "input", "i", "USD", "Input currency lable")
	convertCmd.Flags().StringVarP(&output, "output", "o", "USD", "Output currency lable")
	convertCmd.Flags().Float64VarP(&value, "value", "v", 0, "Target value")

	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")
	convertCmd.MarkFlagRequired("value")

	rootCmd.AddCommand(convertCmd)
}

func GetCurrencySnapshot() {
	result, err := currencyService.Converting(input, output, value)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

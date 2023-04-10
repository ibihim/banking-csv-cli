package main

import (
	"fmt"

	"github.com/ibihim/banking-csv-cli/pkg/cmd"
)

func main() {
	if err := cmd.BankingCommand().Execute(); err != nil {
		fmt.Println(err)
	}
}

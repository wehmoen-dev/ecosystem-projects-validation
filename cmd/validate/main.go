package main

import (
	"flag"
	"fmt"
	"github.com/wehmoen-dev/ecosystem-projects-validation/internal/pkg/data"
	"os"
)

var input string

func main() {
	flag.StringVar(&input, "input", "", "The input string")
	flag.Parse()

	if input == "" {
		panic("No input provided")
	}

	content, err := os.ReadFile(input)

	if err != nil {
		panic(err)
	}

	errors := data.ValidateStructure(content)

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Print("valid")
	}

}

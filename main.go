package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/tadyjp/bigquery-vis/bigquery"
)

func main() {
	if len(os.Args) != 2 {
		pp.Fatalln(len(os.Args))
		pp.Fatalln(os.Args)
		panic("invalid args")
	}

	statements := bigquery.Parse(os.Args[1])
	fmt.Println("======================")
	pp.Println(statements)
}

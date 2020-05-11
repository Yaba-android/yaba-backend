package main

import (
	"fmt"
	"strconv"
	"strings"
)

/********************************
*
* 		TOOLS FUNC
*
********************************/

func concatStringInt(str string, num int64) string {
	return strings.Join([]string{str, strconv.FormatInt(num, 10)}, "")
}

func printStringMap(m map[string]string) {
	for i, n := range m {
		fmt.Print(i)
		fmt.Println(n)
	}
}

func printBookSlice(slice []Book) {
	for i, book := range slice {
		fmt.Print(i)
		fmt.Print(book)
	}
}

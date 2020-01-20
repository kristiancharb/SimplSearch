package search

import (
	"fmt"
)

func (index *Index) Search(query string) {
	terms := tokenize(query)
	fmt.Println(terms)
}

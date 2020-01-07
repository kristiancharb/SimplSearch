package search

import (
	"fmt"
	"regexp"
	"strings"
)

func (index *Index) Search(query string) {
	terms := tokenize(query)
	fmt.Println(terms)
}

func tokenize(query string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	query = reg.ReplaceAllString(query, "")
	return strings.Fields(query)
}

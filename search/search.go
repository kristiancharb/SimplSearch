package search

import (
	"fmt"
)

type queryTerm struct {
	frequency    int
	normalizedTF float64
	tfidf        float64
}

func (index *Index) Search(query string) {
	terms := tokenize(query)
	termsMap := index.getQueryTermsMap(terms)

	for term := range termsMap {
		fmt.Printf("%v: {%v}\n", term, termsMap[term])
	}
}

func (index *Index) getQueryTermsMap(terms []string) map[string]*queryTerm {
	termsMap := make(map[string]*queryTerm)
	indexTerms := index.Terms

	for _, term := range terms {
		if _, present := termsMap[term]; !present {
			termsMap[term] = &queryTerm{}
		}
		termsMap[term].frequency++
	}

	for _, term := range terms {
		if indexTermInfo, present := indexTerms[term]; present {
			normalizedTF := float64(termsMap[term].frequency) / float64(len(terms))
			termsMap[term].normalizedTF = normalizedTF
			termsMap[term].tfidf = indexTermInfo.IDF * normalizedTF
		}
	}

	return termsMap
}

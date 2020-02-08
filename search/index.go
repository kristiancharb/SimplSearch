package search

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type IndexStore struct {
	store map[string]*Index
}

type Index struct {
	Name     string
	DocCount int
	Docs     []string
	Terms    map[string]*TermInfo
}

type TermInfo struct {
	Postings []*Posting
	IDF      float64
	numDocs  int
}

type Posting struct {
	DocID        int
	Frequency    int
	Positions    []int
	NormalizedTF float64
}

func NewIndexStore() *IndexStore {
	return &IndexStore{make(map[string]*Index)}
}

func (indexStore *IndexStore) NewIndex(name string) {
	indexStore.store[name] = &Index{
		Name:  name,
		Terms: make(map[string]*TermInfo),
	}
}

func (indexStore *IndexStore) AddDocument(indexName string, contents string) {
	index := indexStore.store[indexName]
	index.Docs = append(index.Docs, contents)
	index.DocCount++
	docID := len(index.Docs) - 1
	terms := tokenize(contents)
	index.add(terms, docID)
}

func (index *Index) add(terms []string, docID int) {
	pos := 0

	for _, term := range terms {
		var posting *Posting
		if _, present := index.Terms[term]; !present {
			posting = &Posting{DocID: docID}
			termInfo := &TermInfo{Postings: []*Posting{posting}, numDocs: 1}
			index.Terms[term] = termInfo
		} else {
			postings := index.Terms[term].Postings
			posting = getPostingForDoc(postings, docID)
			if posting == nil {
				posting = &Posting{DocID: docID}
				index.Terms[term].Postings = append(postings, posting)
			}
			index.Terms[term].numDocs++
		}

		posting.Frequency++
		posting.NormalizedTF = float64(posting.Frequency) / float64(len(terms))
		posting.Positions = append(posting.Positions, pos)
		pos++
	}

	allTerms := getAllTerms(index.Terms)
	numTerms := len(allTerms)
	for _, term := range allTerms {
		if termInfo, present := index.Terms[term]; present {
			termInfo.IDF = 1.0 + math.Log(float64(numTerms)/float64(termInfo.numDocs))
		}
	}
}

func tokenize(contents string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	contents = reg.ReplaceAllString(contents, "")
	return strings.Fields(contents)
}

func getPostingForDoc(postingList []*Posting, docID int) *Posting {
	for _, posting := range postingList {
		if posting.DocID == docID {
			return posting
		}
	}
	return nil
}

func getAllTerms(termsMap map[string]*TermInfo) []string {
	terms := make([]string, len(termsMap))

	i := 0
	for term := range termsMap {
		terms[i] = term
		i++
	}
	return terms
}

func (index Index) String() string {
	terms := ""
	for term, termInfo := range index.Terms {
		terms += fmt.Sprintf("%v: (%v, %v) => \n", term, termInfo.numDocs, termInfo.IDF)
		postings := termInfo.Postings
		for _, posting := range postings {
			terms += fmt.Sprintf("DocID: %v\n", posting.DocID)
			terms += fmt.Sprintf("Frequency: %v\n", posting.Frequency)
			terms += fmt.Sprintf("Positions: %v\n", posting.Positions)
		}
		terms += "\n"
	}
	return fmt.Sprintf(
		"++++++++++++++++++++++++++++++++++\n"+
			"Name: %v\n"+
			"DocCount: %v\n"+
			"Docs: %v\n"+
			"Terms: \n%v\n"+
			"++++++++++++++++++++++++++++++++++\n",
		index.Name, index.DocCount, index.Docs, terms,
	)
}

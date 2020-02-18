package search

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type IndexStore struct {
	store map[string]*Index
	db    *DBWrapper
}

type Index struct {
	Name     string
	DocCount int
	// Docs     []string
	Terms map[string]*TermInfo
}

type TermInfo struct {
	Postings map[int64]*Posting
	IDF      float64
	numDocs  int
}

type Posting struct {
	DocID        int64
	Frequency    int
	Positions    []int
	NormalizedTF float64
}

func NewIndexStore() *IndexStore {
	store := &IndexStore{
		store: make(map[string]*Index),
		db:    InitDb(),
	}
	store.InitIndex()
	return store
}

func (indexStore *IndexStore) NewIndex(name string) {
	indexStore.store[name] = &Index{
		Name:  name,
		Terms: make(map[string]*TermInfo),
	}
}

func (indexStore *IndexStore) AddDocument(indexName string, title string, contents string, docID int64) {
	index := indexStore.store[indexName]
	if docID == -1 {
		docID = indexStore.db.InsertDoc(indexName, title, contents)
	}
	index.DocCount++
	terms := tokenize(contents)
	index.add(terms, docID)
}

func (indexStore *IndexStore) UpdateIndex(indexName string) {
	index := indexStore.store[indexName]
	allTerms := getAllTerms(index.Terms)
	numTerms := len(allTerms)

	for _, term := range allTerms {
		if termInfo, present := index.Terms[term]; present {
			termInfo.IDF = 1.0 + math.Log(float64(numTerms)/float64(termInfo.numDocs))
		}
	}
}

func (index *Index) add(terms []string, docID int64) {
	pos := 0

	for _, term := range terms {
		var posting *Posting
		if _, present := index.Terms[term]; !present {
			posting = &Posting{DocID: docID}
			termInfo := &TermInfo{
				Postings: make(map[int64]*Posting),
				numDocs:  1,
			}
			termInfo.Postings[docID] = posting
			index.Terms[term] = termInfo
		} else {
			postings := index.Terms[term].Postings
			if posting, present = postings[docID]; !present {
				posting = &Posting{DocID: docID}
				postings[docID] = posting
			}
			index.Terms[term].numDocs++
		}

		posting.Frequency++
		posting.NormalizedTF = float64(posting.Frequency) / float64(len(terms))
		posting.Positions = append(posting.Positions, pos)
		pos++
	}
}

func tokenize(contents string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	contents = reg.ReplaceAllString(contents, "")
	contents = strings.ToLower(contents)
	return strings.Fields(contents)
}

func getPostingForDoc(postingList []*Posting, docID int64) *Posting {
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
			"Terms: \n%v\n"+
			"++++++++++++++++++++++++++++++++++\n",
		index.Name, index.DocCount, terms,
	)
}

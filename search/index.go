package search

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// IndexStore contains index map and DB reference
type IndexStore struct {
	store map[string]*Index
	db    *DBWrapper
}

// Index is an individual inverted index
type Index struct {
	Name     string
	DocCount int
	Terms    map[string]*TermInfo
}

// TermInfo is info for each term in index
type TermInfo struct {
	Postings map[int64]*Posting
	IDF      float64
	numDocs  int
}

// Posting is an individual occurence of a term in a particular doc
type Posting struct {
	DocID        int64
	Frequency    int
	NormalizedTF float64
}

// NewIndexStore initializes the search engine
func NewIndexStore() *IndexStore {
	store := &IndexStore{
		store: make(map[string]*Index),
		db:    InitDb(),
	}
	store.InitIndex()
	return store
}

// NewIndex adds a new index to a store
func (indexStore *IndexStore) NewIndex(name string) {
	indexStore.store[name] = &Index{
		Name:  name,
		Terms: make(map[string]*TermInfo),
	}
}

// AddDocument adds a new document to a particular index
func (indexStore *IndexStore) AddDocument(indexName string, title string, contents string, docID int64) {
	index := indexStore.store[indexName]
	if docID == -1 {
		docID = indexStore.db.InsertDoc(indexName, title, contents)
	}
	index.DocCount++
	terms := tokenize(contents)
	index.add(terms, docID)
}

// UpdateIndex calculates IDF for each term
// This should be called once per transaction
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
	for _, term := range terms {
		var posting *Posting
		if _, present := index.Terms[term]; !present {
			// Term isn't present in index
			// Create new term entry in map and insert new posting
			posting = &Posting{DocID: docID}
			termInfo := &TermInfo{
				Postings: make(map[int64]*Posting),
				numDocs:  1,
			}
			termInfo.Postings[docID] = posting
			index.Terms[term] = termInfo
		} else {
			// Term exists in index
			// Create new posting
			postings := index.Terms[term].Postings
			if posting, present = postings[docID]; !present {
				posting = &Posting{DocID: docID}
				postings[docID] = posting
			}
			index.Terms[term].numDocs++
		}
		posting.Frequency++
		posting.NormalizedTF = float64(posting.Frequency) / float64(len(terms))
	}
}

// Remove punctuation, make lowercase, and split individual words into slice
func tokenize(contents string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	contents = reg.ReplaceAllString(contents, "")
	contents = strings.ToLower(contents)
	return strings.Fields(contents)
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

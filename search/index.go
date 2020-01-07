package search

import (
	"fmt"
	"strings"
	"unicode"
)

type Index struct {
	Name     string
	DocCount int
	Docs     []string
	Terms    map[string][]*Posting
}

type Posting struct {
	DocID        int
	Frequency    int
	Positions    []int
	StartOffsets []int
	EndOffsets   []int
}

func (index *Index) AddDocument(contents string) {
	index.Docs = append(index.Docs, contents)
	index.DocCount++
	docID := len(index.Docs) - 1
	index.tokenize(contents, docID)
	fmt.Println(index)
}

func (index *Index) tokenize(contents string, docID int) {
	var currWord strings.Builder
	pos := 0

	for i, currRune := range contents {
		if unicode.IsLetter(currRune) || unicode.IsDigit(currRune) {
			currWord.WriteRune(currRune)
		} else if unicode.IsSpace(currRune) && currWord.Len() > 0 {
			term := currWord.String()

			var posting *Posting
			if _, present := index.Terms[term]; !present {
				posting = &Posting{DocID: docID}
				index.Terms[term] = []*Posting{posting}
			} else {
				posting = getPostingForSameDoc(index.Terms[term], docID)
			}
			if posting == nil {
				posting = &Posting{DocID: docID}
				index.Terms[term] = append(index.Terms[term], posting)
			}

			posting.Frequency++
			posting.Positions = append(posting.Positions, pos)
			posting.StartOffsets = append(posting.StartOffsets, i-len(term))
			posting.EndOffsets = append(posting.EndOffsets, i)
			pos++
			currWord.Reset()
		}
	}
}

func getPostingForSameDoc(postingList []*Posting, docID int) *Posting {
	for _, posting := range postingList {
		if posting.DocID == docID {
			return posting
		}
	}
	return nil
}

func (index Index) String() string {
	terms := ""
	for term, postingList := range index.Terms {
		terms += term + " => \n"
		for _, posting := range postingList {
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

package search

import (
	"fmt"
	"math"
	"sort"
)

type queryTerm struct {
	frequency    int
	normalizedTF float64
	tfidf        float64
}

type ResponseDoc struct {
	Id       int
	Contents string
}

type QueryResponse struct {
	Length int
	Docs   []ResponseDoc
}

func (index *Index) Search(query string) QueryResponse {
	terms := tokenize(query)
	queryVector, docVectorMap := index.getVectors(terms)
	docRankings := index.getDocRankings(queryVector, docVectorMap)
	response := QueryResponse{}
	for _, docID := range docRankings {
		response.Docs = append(response.Docs, ResponseDoc{
			Id:       docID,
			Contents: index.Docs[docID],
		})
	}
	response.Length = len(response.Docs)
	return response
}

func (index *Index) getVectors(terms []string) ([]float64, map[int][]float64) {
	queryTermFrequencies := make(map[string]int)

	for _, term := range terms {
		if _, present := queryTermFrequencies[term]; !present {
			queryTermFrequencies[term] = 0
		}
		queryTermFrequencies[term]++
	}

	numUniqueTerms := len(queryTermFrequencies)
	queryVector := make([]float64, numUniqueTerms)
	docVectorMap := make(map[int][]float64)
	termIndex := 0
	for term := range queryTermFrequencies {
		if indexTermInfo, present := index.Terms[term]; present {
			normalizedTF := float64(queryTermFrequencies[term]) / float64(numUniqueTerms)
			queryVector[termIndex] = indexTermInfo.IDF * normalizedTF

			for _, posting := range indexTermInfo.Postings {
				if _, present := docVectorMap[posting.DocID]; !present {
					docVectorMap[posting.DocID] = make([]float64, numUniqueTerms)
				}
				tfidf := posting.NormalizedTF * indexTermInfo.IDF
				docVectorMap[posting.DocID][termIndex] = tfidf
			}
		}
		termIndex++
	}
	return queryVector, docVectorMap
}

func (index *Index) getDocRankings(queryVector []float64, docVectorMap map[int][]float64) []int {
	var docs []int
	docScores := make(map[int]float64)
	for docID, docVector := range docVectorMap {
		docs = append(docs, docID)
		docScores[docID] = dotProduct(queryVector, docVector) / (magnitude(queryVector) * magnitude(docVector))
	}
	fmt.Println("Doc Scores:")
	fmt.Println(docScores)
	sort.Slice(docs, func(i, j int) bool {
		return docScores[i] > docScores[j]
	})
	fmt.Println(docs)
	return docs
}

func dotProduct(x []float64, y []float64) float64 {
	sum := 0.0
	if len(x) != len(y) {
		return sum
	}

	for i := range x {
		sum += x[i] * y[i]
	}
	return sum
}

func magnitude(x []float64) float64 {
	sum := 0.0

	for i := range x {
		sum += x[i] * x[i]
	}
	return math.Sqrt(sum)
}

package search

import (
	// "fmt"
	"math"
	"sort"
)

type queryTerm struct {
	frequency    int
	normalizedTF float64
	tfidf        float64
}

type QueryResponse struct {
	Length int
	Docs   []Doc
}

func (indexStore *IndexStore) Search(indexName string, query string, limit int) QueryResponse {
	index := indexStore.store[indexName]
	terms := tokenize(query)
	queryVector, docVectorMap := index.getVectors(terms)
	docs := indexStore.getDocsRanked(indexName, queryVector, docVectorMap, limit)
	response := QueryResponse{
		Length: len(docs),
		Docs:   docs,
	}
	response.Length = len(response.Docs)
	return response
}

func (index *Index) getVectors(terms []string) ([]float64, map[int64][]float64) {
	queryTermFrequencies := make(map[string]int)

	for _, term := range terms {
		if _, present := queryTermFrequencies[term]; !present {
			queryTermFrequencies[term] = 0
		}
		queryTermFrequencies[term]++
	}

	numUniqueTerms := len(queryTermFrequencies)
	queryVector := make([]float64, numUniqueTerms)
	docVectorMap := make(map[int64][]float64)
	termIndex := 0
	for term := range queryTermFrequencies {
		if indexTermInfo, present := index.Terms[term]; present {
			normalizedTF := float64(queryTermFrequencies[term]) / float64(numUniqueTerms)
			queryVector[termIndex] = indexTermInfo.IDF * normalizedTF

			for docId, posting := range indexTermInfo.Postings {
				if _, present := docVectorMap[docId]; !present {
					docVectorMap[docId] = make([]float64, numUniqueTerms)
				}
				tfidf := posting.NormalizedTF * indexTermInfo.IDF
				docVectorMap[docId][termIndex] = tfidf
			}
		}
		termIndex++
	}
	return queryVector, docVectorMap
}

func (indexStore *IndexStore) getDocsRanked(
	indexName string,
	queryVector []float64,
	docVectorMap map[int64][]float64,
	limit int,
) []Doc {
	var docIDs []int64
	docScores := make(map[int64]float64)
	for docID, docVector := range docVectorMap {
		docIDs = append(docIDs, docID)
		docScores[docID] = dotProduct(queryVector, docVector) / (magnitude(queryVector) * magnitude(docVector))
	}
	sort.Slice(docIDs, func(i, j int) bool {
		return docScores[docIDs[i]] > docScores[docIDs[j]]
	})
	if limit > 0 && limit < len(docIDs) {
		docIDs = docIDs[:limit]
	}
	docs := indexStore.db.getDocs(docIDs)
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

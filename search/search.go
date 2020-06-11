package search

import (
	"math"
	"sort"
)

// QueryResponse is the query response sent to the client
type QueryResponse struct {
	Length     int
	NumResults int
	Docs       []Doc
}

// Search queries the given index and returns a QueryResponse for the client
func (indexStore *IndexStore) Search(indexName string, query string, start int, end int) QueryResponse {
	index := indexStore.store[indexName]
	terms := tokenize(query)
	queryVector, docVectorMap := index.getVectors(terms)
	docIDs := indexStore.getDocsRanked(indexName, queryVector, docVectorMap)

	docs := []Doc{}
	numResults := 0
	if len(docIDs) > 0 && start <= len(docIDs)-1 {
		end = min(end, len(docIDs))
		numResults = len(docIDs)
		docIDs = docIDs[start:end]
		docs = indexStore.db.getDocs(docIDs)
	}

	response := QueryResponse{
		Length:     len(docs),
		NumResults: numResults,
		Docs:       docs,
	}
	response.Length = len(response.Docs)
	return response
}

// Returns the query vector and a map of docId: vector
// Vectors are arrays of TFIDF values for each term in the document
func (index *Index) getVectors(terms []string) ([]float64, map[int64][]float64) {
	queryTermFrequencies := getQueryTermFrequencies(terms)
	numUniqueTerms := len(queryTermFrequencies)
	queryVector := make([]float64, numUniqueTerms)
	docVectorMap := make(map[int64][]float64)
	termIndex := 0

	for term := range queryTermFrequencies {
		if indexTermInfo, present := index.Terms[term]; present {
			// Calculate query TFIDF value for current term
			normalizedTF := float64(queryTermFrequencies[term]) / float64(numUniqueTerms)
			queryVector[termIndex] = indexTermInfo.IDF * normalizedTF

			// For each posting, get TFIDF value for current term
			// Insert value into correspond vector from doc vector map
			for docID, posting := range indexTermInfo.Postings {
				if _, present := docVectorMap[docID]; !present {
					docVectorMap[docID] = make([]float64, numUniqueTerms)
				}
				tfidf := posting.NormalizedTF * indexTermInfo.IDF
				docVectorMap[docID][termIndex] = tfidf
			}
		}
		termIndex++
	}
	return queryVector, docVectorMap
}

func getQueryTermFrequencies(terms []string) map[string]int {
	queryTermFrequencies := make(map[string]int)
	for _, term := range terms {
		if _, present := queryTermFrequencies[term]; !present {
			queryTermFrequencies[term] = 0
		}
		queryTermFrequencies[term]++
	}
	return queryTermFrequencies
}

// Rank docs by score and return requested range of docs
func (indexStore *IndexStore) getDocsRanked(
	indexName string,
	queryVector []float64,
	docVectorMap map[int64][]float64,
) []int64 {
	var docIDs []int64
	docScores := make(map[int64]float64)
	for docID, docVector := range docVectorMap {
		docIDs = append(docIDs, docID)
		// Doc score for a doc is the cosine similarity between query and doc vectors
		docScores[docID] = dotProduct(queryVector, docVector) / (magnitude(queryVector) * magnitude(docVector))
	}
	// Use stable sort for consistent search results
	sort.SliceStable(docIDs, func(i, j int) bool {
		iScore := docScores[docIDs[i]]
		jScore := docScores[docIDs[j]]
		if iScore != jScore {
			return iScore > jScore
		}
		return docIDs[i] > docIDs[j]
	})
	return docIDs
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

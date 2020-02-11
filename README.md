# SimplSearch

## A simple Go search engine

## Approach
### Structs
Index: 
```
{
    name string
    docCount int
    documents []string
    terms map[string][]Posting
}
```
Posting:
```
{
    docId int
    frequency int
    positions []int
    startOffsets []int
    endOffsets []int
}
```



### Querying
- Split up query string into individual search terms omitting whitespace and punctuation
- For each search term, get posting list
- Iterate over query terms 
    - Build query vector (slice) by looking up IDF for current term in index
    - Build map of doc vectors
        - Iterate through postings
        - If doc for current posting isn't in map, insert 
        - Calculate tfidf and insert it into doc vector for current doc at the term index
    - Build slice of unique doc IDs for docs included in results
- Sort doc IDs slice by cosine similarity between doc and query vector 

### Resources
- https://sease.io/2015/07/exploring-solr-internals-lucene.html
- https://medium.com/@deangelaneves/how-to-build-a-search-engine-from-scratch-in-python-part-1-96eb240f9ecb
- https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/s


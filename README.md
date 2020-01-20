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

### Example
```
INDEX:

hello => [docId: 1, tf-idf: 1.2], [docId: 2, tf-idf: 3],
how => [docId: 3, tf-idf: 2.1]
are => [docId: 1, tf-idf: 2.0]
you => [docId: 1, tf-idf: 0.1], [docId: 2, tf-idf: 3], [docId: 3, tf-idf: 2.5]

QUERY:
hello you

docs: 
```


### Querying
- Split up query string into individual search terms omitting whitespace and punctuation
- For each search term, get posting list
- Each posting has an tf-idf score that shows how important that term is to the document
- 

### Resources
- https://sease.io/2015/07/exploring-solr-internals-lucene.html
- https://medium.com/@deangelaneves/how-to-build-a-search-engine-from-scratch-in-python-part-1-96eb240f9ecb
- https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/s


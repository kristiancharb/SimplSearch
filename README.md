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
    offsets []Offset
}
```
Offset:
```
{
    start int
    end int
}
```




### Resources
- https://sease.io/2015/07/exploring-solr-internals-lucene.html
- https://medium.com/@deangelaneves/how-to-build-a-search-engine-from-scratch-in-python-part-1-96eb240f9ecb


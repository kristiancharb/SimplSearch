# SimplSearch

## A simple Go search engine Go

## Approach
### Structs
Index: 
```
{
    name string
    segments []Segment
}
```
Segment:
```
{
    docCount int
    documents map[int]Document
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

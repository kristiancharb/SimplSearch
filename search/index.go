package search

import "fmt"

type Index struct {
	Name     string
	Segments []Segment
}

func (index *Index) AddDocument(contents string) {
	if index.Segments == nil {
		index.Segments = append(index.Segments, Segment{
			DocCount: 0,
			Docs:     make(map[int]string),
		})
	}
	segment := index.getSegmentForAdd()
	fmt.Printf("%v\n", segment)
	segment.Docs[segment.DocCount] = contents
	segment.DocCount++
}

func (index *Index) getSegmentForAdd() *Segment {
	i := len(index.Segments) - 1
	return &index.Segments[i]
}

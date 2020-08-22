# SimplSearch

## A simple Go search engine

SimplSearch is a full text search server built with Go and SQLite. 

SQLite is used only for document storage, the search and indexing logic was built completely from scratch. 

[See it in action](https://search.kristiancharb.com/)

### API

- Create a new index

  POST: /index
  ```
  {
    "name": string
  }
  ```
  
- Add document to index

  POST: /index/<name>
  ```
  {
    "title": string,
    "contents": string
  }
  ```
- Search an index
  
  POST: /index/<name?
  ```
  {
    "query": string,
    "start": int (must be > 0,
    "end": 10 (must be > 0, > start)
  }
  ```
  Response:
  ```
  {
    "length": int (number of results returned),
    "numResults": int (total number of results for query),
    "docs": [
      {
        "id": int,
        "index": string,
        "title": string,
        "contents": string,
      }, ...
    ]
  }
  ```


### Indexing

When a new document is added, it's split up into individual tokens. These tokens are just the individual words in the document converted to lowercase with all punctuation removed. The tokens are then inserted into the inverted index.

The inverted index maps individual search terms to the following info:
- Posting map: A posting represents an individual occurence of a term in a document. This data structure maps document IDs to the frequency of the term in the particular document and the normalized term frequency (TF: `term frequency ÷ number of terms in document`).
- Inverse Document Frequency (IDF): `1 + log(total number of terms in index ÷ number of docs containing term)`
- Number of documents

### Searching

When an index is queried, the query is also split up into individual tokens. Then each token is used as a key into the inverted index. All the posting maps for each token are merged to create one master list of document IDs. This list is then sorted based on document rank. These documents are then fetched from the database and returned to the client. 

### Ranking 

Documents are ranked based on how "similar" they are to the query. To measure similarity, first a term vector is constructed for each document and for the query. The term vector is an array containing the TF-IDF values for each term in the document/query. A high TF-IDF value shows that a term shows up often in a document (high TF) but doesn't show up often in other documents (high IDF). The final measure of similarity is the cosine similarity between a particular document vector and the query vector. Before being returned to the client, documents are sorted based on this measure. 


### Resources
- https://sease.io/2015/07/exploring-solr-internals-lucene.html
- https://medium.com/@deangelaneves/how-to-build-a-search-engine-from-scratch-in-python-part-1-96eb240f9ecb
- https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/s


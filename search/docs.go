package search

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"log"
	"strings"
)

// DBWrapper wraps the db struct so we can use it for methods
type DBWrapper struct {
	db *sql.DB
}

// Doc represents the DB schema for the docs table
type Doc struct {
	Id       int64
	Index    string
	Title    string
	Contents string
}

// InitDb creates a new DB instance and the docs table
func InitDb() *DBWrapper {
	db, err := sql.Open("sqlite3", "./docs.db")
	if err != nil {
		fmt.Println("Couldn't open DB")
		log.Fatal(err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY, index_name VARCHAR(255), title TEXT, contents TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Println("Couldn't create table")
		log.Fatal(err)
	}
	return &DBWrapper{db}
}

// InitIndex reads all docs from DB and creates their inverted indexes
func (indexStore *IndexStore) InitIndex() {
	docs := indexStore.db.getAllDocs()
	uniqueIndexes := make(map[string]int)
	i := 0
	for _, doc := range docs {
		if _, present := indexStore.store[doc.Index]; !present {
			indexStore.NewIndex(doc.Index)
		}
		indexStore.AddDocument(doc.Index, doc.Title, doc.Contents, doc.Id)
		uniqueIndexes[doc.Index] = 1

		if i%50 == 0 {
			fmt.Printf("%v documents indexed\n", i)
		}
		i++
	}

	for index := range uniqueIndexes {
		indexStore.UpdateIndex(index)
	}
}

// InsertDoc inserts a document into the DB
func (wrapper *DBWrapper) InsertDoc(index string, title string, contents string) int64 {
	statement, _ := wrapper.db.Prepare("INSERT INTO docs (index_name, title, contents) VALUES (?, ?, ?)")
	result, _ := statement.Exec(index, title, contents)
	id, _ := result.LastInsertId()
	return id
}

func (wrapper *DBWrapper) getAllDocs() []*Doc {
	var docs []*Doc
	rows, err := wrapper.db.Query("SELECT id, index_name, title, contents FROM docs")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for rows.Next() {
		var doc Doc
		rows.Scan(&doc.Id, &doc.Index, &doc.Title, &doc.Contents)
		docs = append(docs, &doc)
	}
	return docs
}

func (wrapper *DBWrapper) getDocs(docIDs []int64) []Doc {
	var docs []Doc
	for _, id := range docIDs {
		var doc Doc
		stmt := "SELECT id, index_name, title, contents FROM docs WHERE id=$1"
		row := wrapper.db.QueryRow(stmt, id)
		err := row.Scan(&doc.Id, &doc.Index, &doc.Title, &doc.Contents)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		docs = append(docs, doc)
	}
	return docs
}

func sliceToString(slice []int64) string {
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", ", ", -1), "[]")
}

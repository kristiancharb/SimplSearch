package search

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type DBWrapper struct {
	db *sql.DB
}

type Doc struct {
	Id       int64
	Index    string
	Title    string
	Contents string
}

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
	args := make([]interface{}, len(docIDs))
	for i, id := range docIDs {
		args[i] = id
	}
	if len(args) < 1 {
		return docs
	}
	stmt := `SELECT id, index_name, title, contents FROM docs WHERE id IN (?` + strings.Repeat(",?", len(args)-1) + `)`
	rows, err := wrapper.db.Query(stmt, args...)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var doc Doc
		rows.Scan(&doc.Id, &doc.Index, &doc.Title, &doc.Contents)
		docs = append(docs, doc)
	}
	return docs
}

func sliceToString(slice []int64) string {
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", ", ", -1), "[]")
}

// func Test() {
// 	wrapper := InitDb()
// 	wrapper.InsertDoc("test", "cool title", "i love java")
// 	wrapper.InsertDoc("test", "cooler title", "sql is good with java")
// 	docs := wrapper.getAllDocs()
// 	fmt.Println(*docs[3])
// }

package search

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// "fmt"
	// "os"
	// "strconv"
	// "strings"
)

func InitDb() {
	db, _ := sql.Open("sqlite3", "./docs.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
}

// func init() *IndexStore {
// 	files := getAllFiles("./docs")
// 	indexStore := IndexStore{make(map[string]*Index)}

// 	for _, file := range files {
// 		indexName, docID := readFileName(file)
// 	}
// }

// func getAllFiles(dir string) []os.FileInfo {
// 	f, err := os.Open(dir)
// 	if err != nil {
// 		fmt.Println("Directory not found")
// 	}
// 	files, err := f.Readdir(-1)
// 	f.Close()
// 	if err != nil {
// 		fmt.Println("There was an error reading the files in the directory")
// 	}
// 	return files
// }

// func readFileName(file os.FileInfo) (string, int) {
// 	name := strings.SplitN(file.Name(), "_", -1)
// 	docID, _ := strconv.Atoi(name[1])
// 	return name[0], docID
// }

// func getDoc(indexName string, docID int) string {

// }

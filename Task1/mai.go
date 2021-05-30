package main
import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	dtbs, _ := sql.Open("sqlite3","mydb.db")
	statement, _ := dtbs.Prepare("CREATE TABLE IF NOT EXISTS Users(rollno INTEGER,name TEXT)")
	statement.Exec()
	addStudent(dtbs,190353,"neil")
	addStudent(dtbs,190122,"feil")
}
func addStudent(db *sql.DB,rollno int, name string) {
	query := "INSERT INTO Users(rollno, name) VALUES(?,?)"
	statement, _ := db.Prepare(query)
	statement.Exec(rollno,name)
}


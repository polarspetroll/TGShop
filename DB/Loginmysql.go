package DB

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func QueryLogin(username, password string) (bool, error) {
	password = fmt.Sprintf("%X", sha256.Sum256([]byte(password)))
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return false, err
	}
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT username FROM login WHERE username=? AND password=?", username, password)
	if err != nil {
		return false, err
	}
	return q.Next(), err
}

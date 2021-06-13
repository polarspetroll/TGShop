package DB

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type QueryOutput struct {
	Id    []int64
	Name  []string
	Price []string
	Stat  []int
	Fname []string
}

var (
	mysqlusr    = os.Getenv("MYSQLUSR")
	mysqlpasswd = os.Getenv("MYSQLPASSWD")
	mysqlhost   = os.Getenv("MYSQLADDR")
	mysqlport   = os.Getenv("MYSQLPORT")
)
var mysqlconncetion = fmt.Sprintf("%v:%v@tcp(%v:%v)/TGShop", mysqlusr, mysqlpasswd, mysqlhost, mysqlport)

func Insert(name, price, filename string, stat int) (affect int64) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	CheckErr(err)
	defer mysqldb.Close()
	p, err := mysqldb.Prepare("INSERT INTO products(name, price, stat, filename) VALUES(?, ?, ?, ?)")
	CheckErr(err)
	row, err := p.Exec(name, price, stat, filename)
	CheckErr(err)
	affect, err = row.RowsAffected()
	CheckErr(err)
	return affect

}

func ListQuery() (output QueryOutput) {
	var id int64
	var name, price string
	var stat int
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	CheckErr(err)
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT name, id, price, stat FROM products")
	CheckErr(err)
	defer q.Close()
	for q.Next() {
		q.Scan(&name, &id, &price, &stat)
		output.Name = append(output.Name, name)
		output.Id = append(output.Id, id)
		output.Price = append(output.Price, price)
		output.Stat = append(output.Stat, stat)
	}
	return output
}

func QueryById(id int64) (out QueryOutput) {
	var name, price, filename string
	var stat int
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	CheckErr(err)
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT name, price, stat, filename FROM products WHERE id=?", id)
	CheckErr(err)
	defer q.Close()
	for q.Next() {
		q.Scan(&name, &price, &stat, &filename)
		out.Price = append(out.Price, price)
		out.Name = append(out.Name, name)
		out.Stat = append(out.Stat, stat)
		out.Fname = append(out.Fname, filename)
	}
	return out
}

func Update(name, price string, stat int, id int64) (affect int64) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	CheckErr(err)
	defer mysqldb.Close()
	p, err := mysqldb.Prepare("UPDATE products SET name=?, price=?, stat=? WHERE id=?")
	CheckErr(err)
	row, err := p.Exec(name, price, stat, id)
	CheckErr(err)
	affect, err = row.RowsAffected()
	CheckErr(err)
	return affect
}

func Delete(id int64) (affect int64) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	CheckErr(err)
	defer mysqldb.Close()
	d, err := mysqldb.Prepare("DELETE FROM products WHERE id=?")
	CheckErr(err)
	row, err := d.Exec(id)
	affect, err = row.RowsAffected()
	CheckErr(err)
	return affect
}

//**********************************************************************************************//
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//**********************************************************************************************//

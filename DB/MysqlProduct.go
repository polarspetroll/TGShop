package DB

import (
	"database/sql"
	"fmt"
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

var OutError QueryOutput

var (
	mysqlusr    = os.Getenv("MYSQLUSR")
	mysqlpasswd = os.Getenv("MYSQLPASSWD")
	mysqlhost   = os.Getenv("MYSQLADDR")
	mysqlport   = os.Getenv("MYSQLPORT")
)
var mysqlconncetion = fmt.Sprintf("%v:%v@tcp(%v:%v)/TGShop", mysqlusr, mysqlpasswd, mysqlhost, mysqlport)

func Insert(name, price, filename string, stat int) (affect int64, err error) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return 0, err
	}
	defer mysqldb.Close()
	p, err := mysqldb.Prepare("INSERT INTO products(name, price, stat, filename) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	row, err := p.Exec(name, price, stat, filename)
	if err != nil {
		return 0, err
	}
	affect, err = row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, err

}

func ListQuery() (output QueryOutput, err error) {
	var id int64
	var name, price string
	var stat int
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return OutError, err
	}
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT name, id, price, stat FROM products")
	if err != nil {
		return OutError, err
	}
	defer q.Close()
	for q.Next() {
		q.Scan(&name, &id, &price, &stat)
		output.Name = append(output.Name, name)
		output.Id = append(output.Id, id)
		output.Price = append(output.Price, price)
		output.Stat = append(output.Stat, stat)
	}
	return output, err
}

func QueryById(id int64) (out QueryOutput, err error) {
	var name, price, filename string
	var stat int
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return OutError, err
	}
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT name, price, stat, filename FROM products WHERE id=?", id)
	if err != nil {
		return OutError, err
	}
	defer q.Close()
	for q.Next() {
		q.Scan(&name, &price, &stat, &filename)
		out.Price = append(out.Price, price)
		out.Name = append(out.Name, name)
		out.Stat = append(out.Stat, stat)
		out.Fname = append(out.Fname, filename)
	}
	return out, err
}

func Update(name, price string, stat int, id int64) (affect int64, err error) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return 0, err
	}
	defer mysqldb.Close()
	p, err := mysqldb.Prepare("UPDATE products SET name=?, price=?, stat=? WHERE id=?")
	if err != nil {
		return 0, err
	}
	row, err := p.Exec(name, price, stat, id)
	if err != nil {
		return 0, err
	}
	affect, err = row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, err
}

func Delete(id int64) (affect int64, err error) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return 0, err
	}
	defer mysqldb.Close()
	d, err := mysqldb.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		return 0, err
	}
	row, err := d.Exec(id)
	affect, err = row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, err
}

func GetFileName(id int64) (name string, err error) {
	mysqldb, err := sql.Open("mysql", mysqlconncetion)
	if err != nil {
		return "", err
	}
	defer mysqldb.Close()
	q, err := mysqldb.Query("SELECT filename FROM products WHERE id=?", id)
	if err != nil {
		return "", err
	}
	if q.Next() {
		q.Scan(&name)
		return name, err
	}
	return "", err
}

package controllers

import (
	//"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	r "github.com/revel/revel"
	"log"
	"myapp/app/models"
	//"time"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

var (
	Dbm *gorp.DbMap
)
const (
	DATE_FORMAT     = "Jan _2, 2006"
	SQL_DATE_FORMAT = "2006-01-02"
)

func InitDB() {

	db, err := sql.Open("mysql", "golangkorea:1234@tcp(127.0.0.1:3306)/golangtest")
	checkErr(err, "sql.Open failed")
	Dbm = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	t := Dbm.AddTableWithName(models.Board{}, "board").SetKeys(true, "Id")
	t.ColMap("DayWrite").Transient = true
	err = Dbm.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	Dbm.TraceOn("[gorp]", r.INFO)
	

	log.Println("gorp 초기화 잘 되었다~~ 꺼윽~~")
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		log.Println("패닉 발생!? 비긴에서?! ")
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		log.Println("패닉 발생!? 커밋에서?! ")
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Println("패닉 발생!? 롤백에서?! ")
		panic(err)
	}
	c.Txn = nil
	return nil
}

package models

import (
	"time"
	"github.com/coopernurse/gorp"
	"fmt"
	"log"
)

type Board struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64 `db:"post_id"`
	Title   string
	Body    string
	Nick    string
	DayWriteStr string
	// Transient
	DayWrite  time.Time
}

const (
	DATE_FORMAT     = "Jan _2, 2006"
	SQL_DATE_FORMAT = "2006-01-02"
)

func (b *Board) PreInsert(_ gorp.SqlExecutor) error {
	log.Println("테이블 추가전 PreInsert")
	b.DayWrite = time.Now()
	b.DayWriteStr = b.DayWrite.Format(SQL_DATE_FORMAT)
	log.Println("테이블 추가후 ")
	log.Println(b)
	return nil
}

func (b *Board) PostGet(_ gorp.SqlExecutor) error {
	var (
		err error
	)
	if b.DayWrite, err = time.Parse(SQL_DATE_FORMAT, b.DayWriteStr); err != nil {
		log.Println("쓴날 에러 검증에러 에러난다...")
		return fmt.Errorf("Error parsing check in date '%s':", b.DayWriteStr, err)
	}
	return nil
}


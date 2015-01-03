package controllers

import (
	"github.com/revel/revel"
	"log"
	"myapp/app/models"
	"myapp/app/routes"
)

type Board struct {
	App
}

// 기본페이지 . 리스트
func (c Board) Index() revel.Result {
	results, err := c.Txn.Select(models.Board{}, `select * from board order by post_id desc`)
	if err != nil {
		panic(err)
	}
	
	var articles []*models.Board
	for _, r := range results {
		b := r.(*models.Board)
		articles = append(articles, b)
	}
	log.Println(articles)
	return c.Render(articles)
}

//하나의 글 읽기
func (c Board) Article(Id int) revel.Result {
	article := c.loadBoardById(Id)
	return c.Render(article)
}

//글 하나 쓰기 폼
func (c Board) FormWrite() revel.Result {
	return c.Render()
}

//글쓰기
func (c Board) Post(board models.Board) revel.Result {
	log.Println(board)
	err := c.Txn.Insert(&board)
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.Board.Index())
}


//글삭제
func (c Board) Delete(Id int64) revel.Result {
	_, err := c.Txn.Delete(&models.Board{Id: Id})
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.Board.Index())
}

//글 수정폼
func (c Board) FormUpdate(Id int) revel.Result{
	article := c.loadBoardById(Id)
	return c.Render(article)
}

//글 수정처리
func (c Board) Update(Id int, Title, Body string) revel.Result{
	_, err := c.Txn.Exec("update board set Title = ?, Body=? where post_id = ?", Title, Body, Id)
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.Board.Article(Id))
}


//글로드
func (c Board) loadBoardById(id int) *models.Board {
	h, err := c.Txn.Get(models.Board{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return h.(*models.Board)
}

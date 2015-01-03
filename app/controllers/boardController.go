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
const (
	COUNT_PER_PAGE = 10;
)

// 기본페이지 . 리스트
func (c Board) Index() revel.Result {	
	articles := c.Page(1)
	log.Println(articles)
	return c.Render(articles)
}

func (c Board) Page(requestPage int64) []*models.Board {
	//페이지 요청후 페이징 계산 알고리즘..복붙복붙;;
	totalArticleCount, err := c.Txn.SelectInt("select count(*) from board")
    checkErr(err, "select count(*) failed")
    log.Println("Rows count:", totalArticleCount)

    beginPage  := (requestPage - 1) / COUNT_PER_PAGE * COUNT_PER_PAGE + 1
    endPage := beginPage + (COUNT_PER_PAGE-1)
    totalPageCount := totalArticleCount / COUNT_PER_PAGE;
	if endPage > totalPageCount{
		endPage = totalPageCount
	}
	firstRow := (requestPage - 1) * COUNT_PER_PAGE 
	endRow := firstRow + COUNT_PER_PAGE 
	if endRow > totalArticleCount{
		endRow = totalArticleCount
	}
	//여기서부터는 sql
	results, err := c.Txn.Select(models.Board{}, 
		`select * from board order by post_id desc limit ?, ?`, firstRow, endRow-firstRow)
	if err != nil {
		panic(err)
	}
	
	var articles []*models.Board
	for _, r := range results {
		b := r.(*models.Board)
		articles = append(articles, b)
	}
	return articles

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

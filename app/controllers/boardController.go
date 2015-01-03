package controllers

import (
	"github.com/revel/revel"
	"log"
	"myapp/app/models"
	"myapp/app/routes"
	"strconv"
)

type Board struct {
	App
}
const (
	COUNT_PER_PAGE = 10;
)

// 기본페이지 . 리스트
func (c Board) Index(RequestPage int64) revel.Result {	
	if RequestPage==0{RequestPage++}
	articles , pageInfo := c.Page(RequestPage)
	log.Println("게시글들과 페이지 정보 ",articles, pageInfo)
	return c.Render(articles, pageInfo)
}

func (c Board) Page(requestPage int64) ([]*models.Board, models.PageInfo){
	//페이지 요청후 페이징 계산 알고리즘..복붙복붙;;
	totalArticleCount, err := c.Txn.SelectInt("select count(*) from board")
    checkErr(err, "select count(*) failed")
    totalPageCount := totalArticleCount / COUNT_PER_PAGE;
    if (totalPageCount % COUNT_PER_PAGE) != 0{
    	totalPageCount++;
    }
    beginPage  := (requestPage - 1) / COUNT_PER_PAGE * COUNT_PER_PAGE + 1
    endPage := beginPage + (COUNT_PER_PAGE-1)
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
	var pageinfo  models.PageInfo
	pageinfo.BeginPage = beginPage
	pageinfo.EndPage = endPage
	pageinfo.TotalPageCount = totalPageCount
	
	
	return articles, pageinfo
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
	return c.Redirect(routes.Board.Index(1))
}
//글 더미 작성
func (c Board) Dummy() revel.Result{
	log.Println("더미를 작성합니다.")
	var board models.Board
	for i:=0;i<100;i++{
		board.Title = "test"+strconv.Itoa(i)
		board.Body ="testtest"
		board.Nick = "nick"
		err := c.Txn.Insert(&board)
		if err != nil {
			panic(err)
		}
	}

	return c.Redirect(routes.Board.Index(1))
}



//글삭제
func (c Board) Delete(Id int64) revel.Result {
	_, err := c.Txn.Delete(&models.Board{Id: Id})
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.Board.Index(1))
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

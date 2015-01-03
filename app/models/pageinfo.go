package models 

import ("fmt")
type PageInfo struct{
	BeginPage int64
	EndPage int64
	TotalPageCount int64
}

func (b PageInfo) PrevBeginPage() int64{
	return b.BeginPage - int64(1)
}
func (b PageInfo) NextEndPage() int64{
	return b.EndPage + int64(1)
}

func (b PageInfo) Pagenation() []int64{
	var pageRow = make([]int64, b.EndPage - b.BeginPage +1)
	for i:=0; b.BeginPage+int64(i) <= b.EndPage;i++ {
		pageRow[i]=b.BeginPage+ int64(i)
	}	
	return pageRow
}


func (b PageInfo) String() string {
	return fmt.Sprintf("PageInfo(%d,%d, %d)", b.BeginPage, b.EndPage, b.TotalPageCount)
}
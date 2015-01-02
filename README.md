revel의 booking 을 조금 단순화시켜서 간단히 CRUD 만 만들어보았습니다.

mysql 버젼이고. 

gorp.go 에 소스를 보시면 db설정이 있습니다. 예)
db, err := sql.Open("mysql", "golangkorea:1234@tcp(127.0.0.1:3306)/golangtest")

이 부분에 맞춰서 database 스키마를 만들어주셔야 되구요. 
ID 가 golangkorea 
비밀번호가 1234 입니다. 

로컬호스트로 접속합니다.

mysql 모듈같은 경우 go get github.com/go-sql-driver/mysql
으로 다운을 받으셔야 할 지도 모릅니다. 

초기버젼은 간단히 게시판에 쓸 CRUD 를 만들고 그 후로는 페이징을 적용해볼까 합니다. 
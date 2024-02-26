package gos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DBInfo struct { //ADBInfo구조체 선언
	DBUser string
	DBPW   string
	DBIP   string
	DBPort int
	DBName string
}

type DBProc struct { //ADBProc 구조체 선언
	DBInfo DBInfo  //구조체 안 구조체 (Nested struct)
	DBConn *sql.DB // DB 관련 선언
}

func NewDBProc(dp DBProc) {
	dp.readConf()
	dp.DBConn = GetConnector(dp.DBInfo)
}

func (dp *DBProc) readConf() {
	file, _ := os.Open("./dbinfo.json")
	defer file.Close()
	decoder := json.NewDecoder(file)  //마샬링,언마샬링(정수형이나 구조체를 바이트 슬라이스로 변경) 말고 많은 데이터를 처리할때 json 문자열을 go 밸류로 바꾸는 것 (디코딩)
	err := decoder.Decode(&dp.DBInfo) //dp 하나만 세팅한게 아니라 dp.DBInfo 를 한 이유 > 선언 한게 두가지가 있기 때문

	if err != nil { //에러시
		fmt.Println("error: ", err)
	}

	fmt.Println("DBUser : ", dp.DBInfo.DBUser) //출력
	fmt.Println("DBPW : ", dp.DBInfo.DBPW)
	fmt.Println("DBIP : ", dp.DBInfo.DBIP)
	fmt.Println("DBPort : ", dp.DBInfo.DBPort)
	fmt.Println("DBName : ", dp.DBInfo.DBName)
}

func GetConnector(dbinfo DBInfo) *sql.DB {
	cfg := mysql.Config{
		User:   dbinfo.DBUser, // "279"
		Passwd: dbinfo.DBPW,   // "279developer",
		Net:    "tcp",
		//Addr:                 "127.0.0.1:3306",
		Addr:                 dbinfo.DBIP + ":" + strconv.Itoa(dbinfo.DBPort),
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               dbinfo.DBName, // "applesensors",
		ParseTime:            true,
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		panic(err)
	}
	db := sql.OpenDB(connector)
	return db
}

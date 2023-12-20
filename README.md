# mqtt2

### golnag 에서 mqtt borker를 사용하여 정보를 전달하는 방법
---
### 1. main.go
golang 폴더에선 기본적으로 하나의 main 함수만 취급한다. 

func main() { } 부터 읽어 들이기 시작한다.

main.go 파일에서
```go
package main

import (
	"sensor_server/gos"
)

func main() {
	dp := gos.DBProc{}
	gos.NewDBProc(dp)
	gos.NewMqttClient(dp)
}
```
gos 라는 폴더를 사용하기위해 import 해준다.("sensor_server"는 go.mod 파일 생성시 module을 sensor_server로 작성해서 그러하다.)

main 함수를 보면 dp 라는 생성자를 통해 새로운 구조체를 선언하고 형식은 dbproc.go 파일에 있는 DBProc 구조체 형식이다.

gos 파일 안에 있는 NewDBProc 함수를 통해 초기화를 해준다. 

(초기화 내용에는 DB 정보를 읽는 readConf() 함수와 연결하는 부분인 GetConnector() 함수가 있다.)


다음으로는 gos 폴더 안에 있는 client.go 파일에 NewMqttClient() 함수를 실행하여 초기화하는데 초기화 안에는 생성자 와 readConf(), init(), sub() 함수로 이뤄져 있다.

### 2. dbproc.go
main.go 함수에서 호술하는 파일이다.

```go
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
	DBInfo DBInfo  //구조체 안 구조체 (Nested struct?)
	DBConn *sql.DB // DB 관련 선언
}

func NewDBProc(dp DBProc) {
	dp.readConf()
	dp.DBConn = GetConnector(dp)
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

func GetConnector(dp DBProc) *sql.DB {
	cfg := mysql.Config{
		User:   dp.DBInfo.DBUser, // "279"
		Passwd: dp.DBInfo.DBPW,   // "279developer",
		Net:    "tcp",
		//Addr:                 "127.0.0.1:3306",
		Addr:                 dp.DBInfo.DBIP + ":" + strconv.Itoa(dp.DBInfo.DBPort),
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               dp.DBInfo.DBName, // "applesensors",
		ParseTime:            true,
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		panic(err)
	}
	db := sql.OpenDB(connector)
	return db
}
```
gos 폴더 안에 있는 파일이라 package 는 gos로 되어 있고 

사용할 패키지를 임포트를 통해 불러와 사용한다.

구조체는 두개를 선언한다.(구조체는 일종의 클래스라 보면 된다. 다만 필드만 존재하는 클래스이다.)

DB 정보를 불러와 저장하는 DBInfo 구조체, DB정보가 담긴 구조체와 연결부분을 저장하는 DBProc 구조체 두개가 있다.

다음으로는 앞서 말한거 처럼 NewDBProc()함수가 있는데 이 함수는 DBProc 값을 대입받는 함수이다.

readConf() 함수는 함수 앞 dp *DBProc 를 지정하여 DBProc 를 위한 메소드 임을 표시한다.
```go
file, _ := os.Open("./dbinfo.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&dp.DBInfo) 
```
os.Open 을 통해 dbinfo 파일을 열어 읽어온다. 다음 json.NewDecoder 함수로 디코더를 만든 후 json.Decode를 통해 JSON 문자열을 GO 밸류로 변경하게 된다. (JSON 문자열을 GO 밸류로 바꾸는 이작업을 디코딩이라 한다.)

그 다음 dp 라는 이름의 객체에 DBInfo 안 필드에 데이터가 잘 들어 갔는지 출력을 통해 확인한다.

함수 GetConnector(dp DBProc) *sql.DB 은 연결 부분이다. 파라미터로 DBProc 구조체가 들어간다는 말이고 리턴 값이 sql.DB 임을 뜻한다.

cfg  안 내용은 연결 부분이다. 필요한 구조체 내용으로 대체하면 된다.
```go
connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		panic(err)
	}
	db := sql.OpenDB(connector)
	return db
```
mysql.NewConnector 를 통해 connector를 생성하고 OpenDB에 인자로 넣어주면 된다.
위 값이 리턴 값이 되고 리턴 값은 같은 형태인 DBProc 필드 DBConn으로 들어가며 연결이 된다.

### 3. client.go

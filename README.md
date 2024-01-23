# mqtt2

### golnag 에서 mqtt borker를 사용하여 정보를 전달하는 방법
---
### 0. mqtt관련 라이브러리 설치
golang에서 mqtt broker 서버를 만들기 위해 필요한 라이브러리가 있다.

github.com/eclipse/paho.mqtt.golang 라는 라이브러리 이다.

위 라이브러리를 사용하기 위해선 터미널에서 go get 명령으로 다운받을 수 있다.
```
go get github.com/eclipse/paho.mqtt.golang 
```
go.mod 에 라이브러리가 입력되면 go.sum에 필요한 것들이 전부 다운받아진다. (go get없이도 mod차원에서 자동으로 될 수 있으나 명령으로 치는게 깔끔해 보인다.)

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

main 함수를 보면 dp 라는 객체를 선언하고 형식은 dbproc.go 파일에 있는 DBProc 구조체 형식이다.

gos 파일 안에 있는 NewDBProc 함수(생성자함수)를 통해 초기화를 해준다. 

(초기화 내용에는 DB 정보를 읽는 readConf() 함수와 연결하는 부분인 GetConnector() 함수가 있다.)


다음으로는 gos 폴더 안에 있는 client.go 파일에 NewMqttClient()함수(생성자함수)를 실행하여 초기화하는데 초기화 안에는 생성자 와 readConf(), init(), sub() 함수로 이뤄져 있다.

### 2. dbproc.go
main.go 함수에서 호출한 함수가 있는 파일이다.

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

다음으로는 앞서 말한거 처럼 생성자함수 NewDBProc()함수가 있는데 이 함수는 DBProc 값을 대입받는 함수이다.

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
main 함수에서 dp를 생성하고 초기화를 한 다음 gos 폴더 안에 있는 생성자함수 NewMqttClinet()를 호출하고 값을 넣었다.

clinet.go 파일을 보자
```go
package gos

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttInfo struct {
	BrokerIP   string
	BrokerPort int
	Topic      string
	ClientID   string
	SiteID     int
}

type MqttClient struct {
	MqtInfo         MqttInfo
	MsgHandler      mqtt.MessageHandler
	ProcConnect     mqtt.OnConnectHandler
	ProcLostConnect mqtt.ConnectionLostHandler
	MqttClient      mqtt.Client
	isConnected     bool
}

func NewMqttClient(dp DBProc) {
	mq := MqttClient{}
	mq.isConnected = false

	mq.readConf()
	mq.init()

	mq.MsgHandler = func(client mqtt.Client, msg mqtt.Message) {
		str := string(msg.Payload())
		name := str[0:5]
		value := str[6:]
		fmt.Print(name + ":" + value)
		_, err := dp.DBConn.Exec("INSERT INTO sensorlogs(name, temp) value(?,?)", name, value)
		if err == nil {
		} else {
			log.Println(err)
			print("에러1")
		}
	}
	mq.ProcConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	mq.ProcLostConnect = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
	mq.sub()
}

func (mq *MqttClient) readConf() { //json 파일 읽기
	file, _ := os.Open("./mqtt.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&mq.MqtInfo) //DB쪽이라 생각하면 될듯

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("BrokerIP : ", mq.MqtInfo.BrokerIP)
	fmt.Println("BrokerPort : ", mq.MqtInfo.BrokerPort)
	fmt.Println("Topic : ", mq.MqtInfo.Topic)
	fmt.Println("ClientID : ", mq.MqtInfo.ClientID)
	fmt.Println("SiteID : ", mq.MqtInfo.SiteID)
}

func (mq *MqttClient) init() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mq.MqtInfo.BrokerIP, mq.MqtInfo.BrokerPort)) //로컬브로커
	opts.SetClientID(mq.MqtInfo.ClientID)

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(mq.MsgHandler)
	opts.OnConnect = mq.ProcConnect
	opts.OnConnectionLost = mq.ProcLostConnect

	mq.MqttClient = mqtt.NewClient(opts) //여기서부터 시작 고치기 완료
	if token := mq.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		mq.isConnected = true
	}
}

func (mq *MqttClient) sub() { //topic 방 만들기
	token := mq.MqttClient.Subscribe(mq.MqtInfo.Topic, 1, nil) //func (mqtt.Client).Subscribe(topic string, qos byte, callback mqtt.MessageHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", mq.MqtInfo.Topic)
	for i := 1; true; i++ {
		time.Sleep(60 * time.Second)
	}
}
```
이 파일은 두개의 구조체가 있다.브로커IP와 port, topic 같은 정보를 저장하는 MqttInfo 구조체와

클라인언트 정보와 서버연결 정보를 저장하는 MqttClient 구조체가 있다.

바로 밑에는 main 함수에서 실행한 생성자함수(NewMqttClient)를 볼 수 있다.
```go
mq := MqttClient{}
	mq.isConnected = false

	mq.readConf()
	mq.init()
```
MqttClient 객체(mq)를 만들고 필드값(isConnected)을 설정(false)

MqttClient를 위한 메소드 함수 readConf(),init() 또한 실행시킨다. 

readConf()함수는 dbporc에서도 나왔다 시피 mqtt.json 파일을 읽고 저장한다.

init() 함수는 mqtt 브로커 초기설정을 위한 함수이다.(어쩌면 생성자함수와 같은 결일 수 있다.)

import 되어있는 mqtt "github.com/eclipse/paho.mqtt.golang" 패키지를 사용하여 클라이언트 옵션설정을 한다.

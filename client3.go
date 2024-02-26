package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-sql-driver/mysql"
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

func NewMqttClient(dbpro DBProc) {
	mq := MqttClient{}
	mq.isConnected = false

	mq.readConf()

	mq.MsgHandler = func(client mqtt.Client, msg mqtt.Message) {
		str := string(msg.Payload())
		name := str[0:4]
		value := str[5:]
		fmt.Println(name + ":" + value)
		_, err := dbpro.DBConn.Exec("INSERT INTO sensorlogs(name, temp) value(?,?)", name, value)
		if err == nil {
		} else {
			log.Println(err)
		}
	}
	mq.ProcConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	mq.ProcLostConnect = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
	mq.init()
	mq.sub()
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

func main() {
	dp := DBProc{}
	dp.readConf()
	dp.DBConn = GetConnector(dp.DBInfo)
	NewMqttClient(dp)
}

//client3은 건들지 말것 .. 완성합체본

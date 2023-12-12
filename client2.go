// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"time"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// 	"github.com/go-sql-driver/mysql"
// )

// type MqttInfo struct {
// 	BrokerIP   string
// 	BrokerPort int
// 	Topic      string
// 	ClientID   string
// 	SiteID     int
// }

// type DBInfo struct { //ADBInfo구조체 선언
// 	DBUser string
// 	DBPW   string
// 	DBIP   string
// 	DBPort int
// 	DBName string
// }

// type DBProc struct { //ADBProc 구조체 선언
// 	DBInfo DBInfo  //구조체 안 구조체 (Nested struct?)
// 	DBConn *sql.DB // DB 관련 선언
// }

// func (dp *DBProc) readConf() {
// 	file, _ := os.Open("./dbinfo.json")
// 	defer file.Close()
// 	decoder := json.NewDecoder(file)  //마샬링,언마샬링(정수형이나 구조체를 바이트 슬라이스로 변경) 말고 많은 데이터를 처리할때 json 문자열을 go 밸류로 바꾸는 것 (디코딩)
// 	err := decoder.Decode(&dp.DBInfo) //dp 하나만 세팅한게 아니라 dp.DBInfo 를 한 이유 > 선언 한게 두가지가 있기 때문

// 	if err != nil { //에러시
// 		fmt.Println("error: ", err)
// 	}
// }

// func GetConnector(dbinfo DBInfo) *sql.DB {
// 	cfg := mysql.Config{
// 		User:   dbinfo.DBUser, // "279"
// 		Passwd: dbinfo.DBPW,   // "279developer",
// 		Net:    "tcp",
// 		//Addr:                 "127.0.0.1:3306",
// 		Addr:                 dbinfo.DBIP + ":" + strconv.Itoa(dbinfo.DBPort),
// 		Collation:            "utf8mb4_general_ci",
// 		Loc:                  time.UTC,
// 		MaxAllowedPacket:     4 << 20.,
// 		AllowNativePasswords: true,
// 		CheckConnLiveness:    true,
// 		DBName:               dbinfo.DBName, // "applesensors",
// 		ParseTime:            true,
// 	}
// 	connector, err := mysql.NewConnector(&cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	db := sql.OpenDB(connector)
// 	return db
// }

// var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	// fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// 	str := string(msg.Payload())
// 	name := str[0:5]
// 	value := str[6:]
// 	fmt.Print(name + ":" + value)

// 	dp := DBProc{}
// 	dp.readConf()
// 	dp.DBConn = GetConnector(dp.DBInfo)

// 	/*result*/
// 	_, err := dp.DBConn.Exec("INSERT INTO sensorlogs(name, temp) value(?,?)", name, value)
// 	if err == nil {
// 	} else {
// 		log.Println(err)
// 	}
// }

// var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
// 	fmt.Println("Connected")
// }

// var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
// 	fmt.Printf("Connect lost: %v", err)
// }

// func (mq *MqttInfo) readConf() { //json 파일 읽기
// 	file, _ := os.Open("./mqtt.json")
// 	defer file.Close()
// 	decoder := json.NewDecoder(file)
// 	err := decoder.Decode(&mq) //DB쪽이라 생각하면 될듯

// 	if err != nil {
// 		fmt.Println("err: ", err)
// 	}

// 	fmt.Println("BrokerIP : ", mq.BrokerIP)
// 	fmt.Println("BrokerPort : ", mq.BrokerPort)
// 	fmt.Println("Topic : ", mq.Topic)
// 	fmt.Println("ClientID : ", mq.ClientID)
// 	fmt.Println("SiteID : ", mq.SiteID)
// }

// func (mq *MqttInfo) init() {

// 	opts := mqtt.NewClientOptions()
// 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mq.BrokerIP, mq.BrokerPort)) //로컬브로커
// 	opts.SetClientID("go_mqtt_client")
// 	opts.SetUsername("emqx")
// 	opts.SetPassword("public")
// 	opts.SetDefaultPublishHandler(messagePubHandler)
// 	opts.OnConnect = connectHandler
// 	opts.OnConnectionLost = connectLostHandler
// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}

// 	mq.sub(client)

// 	client.Disconnect(250)
// }

// func (mq *MqttInfo) sub(client mqtt.Client) { //topic 방 만들기
// 	token := client.Subscribe(mq.Topic, 1, nil) //func (mqtt.Client).Subscribe(topic string, qos byte, callback mqtt.MessageHandler)
// 	token.Wait()
// 	fmt.Printf("Subscribed to topic: %s\n", mq.Topic)
// 	for i := 0; true; i++ {
// 		time.Sleep(60 * time.Second)
// 	}
// }

// func main() {
// 	mq := MqttInfo{}
// 	mq.readConf()
// 	mq.init()
// }

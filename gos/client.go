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

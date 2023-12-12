package gos

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

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

func (mq *MqttInfo) init1() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mq.BrokerIP, mq.BrokerPort)) //로컬브로커
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mq.pub(client)

	client.Disconnect(250)
}

func (mq *MqttInfo) pub(client mqtt.Client) {

	num := 71
	for i := 65; i < num; i++ {
		text := fmt.Sprintf("%d", i)
		token := client.Publish(mq.Topic, 0, false, text) //Publish(topic string, qos byte, retained bool, payload interface{})
		token.Wait()
		time.Sleep(5 * time.Second)
	}
}

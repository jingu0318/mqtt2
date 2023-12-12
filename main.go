package main

import (
	"sensor_server/gos"
)

func main() {
	dp := gos.DBProc{}
	gos.NewDBProc(dp)
	gos.NewMqttClient(dp)
}

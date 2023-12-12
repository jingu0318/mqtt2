# mqtt2

### golnag 에서 mqtt borker를 사용하여 정보를 전달하는 방법
---
### 1. main.go
golang 폴더에선 기본적으로 하나의 main 함수만 취급한다. 

func main() { } 부터 읽어 들이기 시작한다.

main.go 파일에서
```
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


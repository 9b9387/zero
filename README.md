# zero
A simple TCP Server with heartbeat

## Requirements

```
$ go get github.com/satori/go.uuid
```

## Usage

```
func main() {
 	host := "127.0.0.1:18787"

 	ss, err := zero.NewSocketService(host)
	if err != nil {
		return
	}

	// set Heartbeat
	//ss.SetHeartBeat(5*time.Second, 30*time.Second)

	// net event
	//ss.RegOnMessageHandler(HandleMessage)
	//ss.RegOnConnectHandler(HandleConnect)
	//ss.RegOnDisconnectHandler(HandleDisconnect)

	ss.Serv()
}

```

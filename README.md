# zero
A simple TCP Server with heartbeat

## Usage

```
func main() {
 	host := "127.0.0.1:18787"

 	ss, err := zero.NewSocketService(host)
	if err != nil {
		return
	}

	ss.SetHeartBeat(5*time.Second, 30*time.Second)

	ss.RegOnMessageHandler(HandleMessage)
	ss.RegOnConnectHandler(HandleConnect)
	ss.RegOnDisconnectHandler(HandleDisconnect)

	ss.Serv()
}

```

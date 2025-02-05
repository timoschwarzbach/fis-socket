package ibis

import (
	"encoding/xml"
	"fis/socket/ibis/definitions"
	"fis/socket/socket"
	"fmt"
	"io"
	"net/http"
)

/*
The Ibis controller accesses data from the Ibis interfaces
and provides updates to the client via websocket
Requested services:

  - CustomerInformationService
  - GNSSLocationService
*/
type IbisController struct {
	socketServer *socket.Server
}

func CreateController(server *socket.Server) *IbisController {
	c := &IbisController{
		socketServer: server,
	}

	c.socketServer.On("connection", func(clients ...any) {
		c.fullPullPush()
	})
	c.fullPullPush()

	return c
}

func (c *IbisController) send(data any) {
	c.socketServer.Emit("ibis", data)
}

func (c *IbisController) fullPullPush() {
	url := "http://localhost:2092/CustomerInformationService"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var responseData definitions.GetAllDataResponse
	err = xml.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return
	}

	// jsonData, err := json.Marshal(responseData)
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// }

	c.send(responseData.AllData)
}

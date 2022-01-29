package server

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rbren/midi/pkg/logger"
)

var upgrader = websocket.Upgrader{} // use default options

func StartServer() {
	http.HandleFunc("/echo", echo)
	http.Handle("/", http.FileServer(http.Dir("web")))
	logger.ForceLog(http.ListenAndServe(":8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ForceLog("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.ForceLog("read:", err)
			break
		}
		logger.ForceLog("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.ForceLog("write:", err)
			break
		}
	}
}

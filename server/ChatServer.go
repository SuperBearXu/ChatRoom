package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var Up = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	conns = make(map[string]*websocket.Conn)
	mutex sync.Mutex
)

func handleWebSocket(res http.ResponseWriter, req *http.Request) {
	conn, err := Up.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	username := req.URL.Query().Get("username")
	fmt.Println(username + " connected")
	mutex.Lock()
	conns[username] = conn
	mutex.Unlock()

	// 发送聊天室成员列表
	//conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the chat room!"))
	historyMsg, err := json.Marshal(map[string]string{
		"type":  "users",
		"users": conns[username].RemoteAddr().String(),
	})
	if err != nil {
		log.Println("json 序列化失败：", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, historyMsg)
	if err != nil {
		return
	}

	// 心跳设置
	err = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return
	}
	conn.SetPongHandler(func(string) error {
		err = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	defer func() {
		mutex.Lock()
		for user, c := range conns {
			if c == conn {
				delete(conns, user)
				break
			}
		}
		mutex.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(conn.RemoteAddr(), string(msg))

		mutex.Lock()
		tmpConns := make([]*websocket.Conn, 0, len(conns))
		for _, c := range conns {
			tmpConns = append(tmpConns, c)
		}
		mutex.Unlock()

		for _, c := range tmpConns {
			if c != conn && c != nil {
				err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s: %s", username, msg)))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws_test", handleWebSocket)
	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Println(err)
	}
}

//func SendUserList(conns map[string]*Client) {
//	//fmt.Println("执行广播功能： ", len(conns))
//	//for name, _ := range conns {
//	//	fmt.Println("查看广播对象： ", name)
//	//}
//	// 发送聊天室成员列表
//	users := GetAlKeys(conns)
//	usersMsg, err := json.Marshal(map[string]interface{}{
//		"type":  "users",
//		"users": users,
//	})
//	if err != nil {
//		log.Println("json 序列化失败：", err)
//		return
//	}
//
//	tmpConns := make([]*websocket.Conn, 0, len(conns))
//	for _, c := range conns {
//		tmpConns = append(tmpConns, c.conn)
//	}
//
//	for _, c := range tmpConns {
//		if c != nil {
//			err = c.WriteMessage(websocket.TextMessage, usersMsg)
//			if err != nil {
//				log.Println(err)
//			}
//		}
//	}
//}
//
//func GetAlKeys(m map[string]*Client) []string {
//	keys := make([]string, 0, len(m))
//	for key := range m {
//		keys = append(keys, key)
//	}
//	return keys
//}

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
	conns = make(map[string]*Client)
	mutex sync.Mutex
)

type Client struct {
	conn       *websocket.Conn
	lastActive time.Time
	username   string
	done       chan struct{}
}

func handleWebSocket(res http.ResponseWriter, req *http.Request) {
	conn, err := Up.Upgrade(res, req, nil)
	if err != nil {
		log.Println("升级WebSocket失败:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("关闭WebSocket连接失败:", err)
		}
	}(conn)

	username := req.URL.Query().Get("username")
	if username == "" {
		err = conn.WriteMessage(websocket.CloseMessage, []byte("需要用户名"))
		if err != nil {
			return
		}
		return
	}

	client := &Client{
		conn:       conn,
		lastActive: time.Now(),
		username:   username,
		done:       make(chan struct{}),
	}

	fmt.Println(username + " connected")

	mutex.Lock()
	conns[username] = client
	mutex.Unlock()

	SendUserList(conns)

	// 配置ping/pong处理器
	heartbeatInterval := 25 * time.Second
	pongWait := 60 * time.Second

	client.conn.SetPingHandler(func(string) error {
		err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Println("SetPingHandler 设置超时退出时间：" + err.Error())
			return err
		}
		return client.conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
	})

	client.conn.SetPongHandler(func(string) error {
		client.lastActive = time.Now()
		err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Println("SetPongHandler 设置超时退出时间：" + err.Error())
			return err
		}
		return nil
	})

	// 启动定时ping
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := client.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					log.Println("发送ping失败：", err)
					return
				}
			case <-client.done:
				log.Println("心跳检测已停止")
				return
			}
		}
	}()

	defer func() {
	mutex.Lock()
	defer mutex.Unlock()
	for user, c := range conns {
		if c.conn == conn {
			delete(conns, user)
			log.Printf("连接已关闭：%s (%s)", user, conn.RemoteAddr())
			SendUserList(conns)
			break
		}
	}
}()

	for {
		_, msgJson, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取信息error: %v", err)
			break
		}

		fmt.Println(conn.RemoteAddr(), string(msgJson))

		// 更新最后活跃时间
		client.lastActive = time.Now()

		mutex.Lock()
		tmpConns := make([]*websocket.Conn, 0, len(conns))
		for _, c := range conns {
			tmpConns = append(tmpConns, c.conn)
		}
		mutex.Unlock()

		for _, c := range tmpConns {
			if c != nil {
				err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s: %s", username, msgJson)))
				if err != nil {
					log.Println("广播信息失败：", err)
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

func SendUserList(conns map[string]*Client) {
	// 发送聊天室成员列表
	users := GetAlKeys(conns)
	usersMsg, err := json.Marshal(map[string]interface{}{
		"type":  "users",
		"users": users,
	})
	if err != nil {
		log.Println("json 序列化失败：", err)
		return
	}

	mutex.Lock()
	tmpConns := make([]*websocket.Conn, 0, len(conns))
	for _, c := range conns {
		tmpConns = append(tmpConns, c.conn)
	}
	mutex.Unlock()

	for _, c := range tmpConns {
		if c != nil {
			err = c.WriteMessage(websocket.TextMessage, usersMsg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func GetAlKeys(m map[string]*Client) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

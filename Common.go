package main

//func SendUserList(conns map[string]*Client) {
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
//	mutex.Lock()
//	tmpConns := make([]*websocket.Conn, 0, len(conns))
//	for _, c := range conns {
//		tmpConns = append(tmpConns, c.conn)
//	}
//	mutex.Unlock()
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

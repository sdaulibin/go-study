package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int64 // 发送者
	TargetId int64 // 接收者
	Type     int   // 1-私聊，2-群聊，3-广播
	Media    int   // 1-文字，2-表情包，3-图片，4-音频
	Content  string
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	//token := query.Get("token")
	//targetId := strconv.ParseInt(query.Get("targetId"), 10, 64)
	//context := query.Get("context")
	//msgType := query.Get("type")

	isvalidate := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalidate
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	fmt.Println("message start")
	go sendProc(node)
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室。。。。。。"))
}

func sendProc(node *Node) {
	fmt.Println("[ws] sendProc >>>>> start")
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	fmt.Println("[ws] recvProc >>>>> start ")
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		fmt.Println("[ws] recvMsg <<<<<< msg: ", string(data))
	}
}

var upSendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	upSendChan <- data
}

func init() {
	go udpSendProc()
	go updRecvProc()
	fmt.Println("init goroutine ......")
}

func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		select {
		case data := <-upSendChan:
			fmt.Println("udpSendProc data: ", string(data))
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func updRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("dispatch msg type: ", msg.Type)
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch data: ", string(data))
		sendMsg(msg.TargetId, data)
		//case 2: //群发
		//	sendGroupMsg()
		//case 3: //广播
		//	sendAllMsg()
	}
}

func sendMsg(userId int64, data []byte) {
	fmt.Println("sendMsg userId>>>>> ", userId)
	for k, v := range clientMap {
		fmt.Println(k, "<<<<>>>>", v)
	}
	rwLocker.RLock()
	node, ok := clientMap[userId]
	fmt.Println("sendMsg OK:", ok)
	if ok {
		node.DataQueue <- data
	}
	rwLocker.RUnlock()
}

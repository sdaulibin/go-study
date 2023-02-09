package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId     int64 // 发送者
	TargetId   int64 // 接收者
	Type       int   // 1-私聊，2-群聊，3-广播
	Media      int   // 1-文字，2-表情包，3-图片，4-音频
	Content    string
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn          *websocket.Conn
	Addr          string //客户端地址
	DataQueue     chan []byte
	GroupSets     set.Interface
	FirstTime     uint64 //首次连接时间
	HeartbeatTime uint64
	LoginTime     uint64
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.  获取参数 并 检验 token 等合法性
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
	//2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
		FirstTime:     currentTime,
		HeartbeatTime: currentTime,
	}
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//发送消息协程
	go sendProc(node)
	//接收消息协程
	go recvProc(node)
	SetUserOnlineInfo("online_"+id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
	sendMsg(userId, []byte("欢迎进入聊天室。。。。。。"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err)
		}
		//心跳检测 msg.Media == -1 || msg.Type == 3
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			dispatch(data)
			//broadMsg(data)
			fmt.Println("[ws] recvMsg <<<<<< msg: ", string(data))
		}
	}
}

// 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
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
		IP:   net.IPv4(172, 168, 255, 255),
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
		Port: viper.GetInt("port.udp"),
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
		fmt.Println("udpSendProc data: ", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
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

func sendMsg(userId int64, msg []byte) {
	// for k, v := range clientMap {
	// 	fmt.Println(k, "<<<<>>>>", v)
	// }
	rwLocker.RLock()
	node, ok := clientMap[userId]
	//fmt.Println("sendMsg OK:", ok)
	rwLocker.RUnlock()
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.FromId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	r, err := utils.RedisClient.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}
	if r != "" {
		if ok {
			fmt.Println("sendMsg >>>>> userId: ", userId, " msg: ", string(msg))
			node.DataQueue <- msg
		}
	}
	var key string
	if userId > jsonMsg.FromId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.RedisClient.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	result, err := utils.RedisClient.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	var rels []string
	var err error
	if isRev {
		rels, err = utils.RedisClient.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.RedisClient.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	return rels
}

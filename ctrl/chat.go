package ctrl

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"` //消息ID
	//谁发的
	Userid  int64  `json:"userid,omitempty" form:"userid"` //谁发的
	//什么业务
	Cmd     int    `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	//发给谁
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`//对端用户ID/群ID
	//怎么展示
	Media   int    `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	//内容是什么
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	//图片是什么
	Pic     string `json:"pic,omitempty" form:"pic"` //预览图片
	//连接是什么
	Url     string `json:"url,omitempty" form:"url"` //服务的URL
	//简单描述
	Memo    string `json:"memo,omitempty" form:"memo"` //简单描述
	//其他的附加数据，语音长度/红包金额
	Amount  int    `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}
const (
	//点对点单聊,dstid是用户ID
	CMD_SINGLE_MSG = 10
	//群聊消息,dstid是群id
	CMD_ROOM_MSG   = 11
	//心跳消息,不处理
	CMD_HEART      = 0

)
const (
	//文本样式
	MEDIA_TYPE_TEXT=1
	//新闻样式,类比图文消息
	MEDIA_TYPE_News=2
	//语音样式
	MEDIA_TYPE_VOICE=3
	//图片样式
	MEDIA_TYPE_IMG=4

	//红包样式
	MEDIA_TYPE_REDPACKAGR=5
	//emoj表情样式
	MEDIA_TYPE_EMOJ=6
	//超链接样式
	MEDIA_TYPE_LINK=7
	//视频样式
	MEDIA_TYPE_VIDEO=8
	//名片样式
	MEDIA_TYPE_CONCAT=9
	//其他自己定义,前端做相应解析即可
	MEDIA_TYPE_UDEF=100
)
/**
消息发送结构体,点对点单聊为例
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,
content:"hello"}

3、MEDIA_TYPE_VOICE,amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,
url:"http://www.a,com/dsturl.mp3",
amount:40}

4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,
url:"http://www.baidu.com/a/log.jpg"}


2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,
content:"标题",
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/dsturl",
"memo":"这是描述"}


5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}

7、MEDIA_TYPE_Link 7
{id:1,userid:2,dstid:3,cmd:10,media:7,
"url":"http://www.a.com/dsturl.html"
}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/a.mp4"
}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,
"content":"10086",
"pic":"http://www.baidu.com/a/avatar,jpg",
"memo":"胡大力"}
*/


//核心在于形成userid和Node的映射关系
type Node struct{
	Conn *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node)

var rwlocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request){
	//检验接入是否合法 CheckToken
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isValid := checkToken(userId, token)
	//如果isValid为true
	//todo 如果isValid为False
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(w, r, nil)

	if err!=nil{
		log.Println(err.Error())
		return
	}
	//获得conn
	node := &Node{
		Conn:conn,
		//并行数据转串行数据
		DataQueue:make(chan []byte,50),
		GroupSets:set.New(set.ThreadSafe),
	}

	//todo 获取用户全部群Id
	comIds := contactService.SearchComunityIds(userId)
	for _,v:=range comIds{
		node.GroupSets.Add(v)
	}

	//userid和node形成绑定关系，读写锁，并发不出错
	rwlocker.Lock()
	clientMap[userId]=node
	rwlocker.Unlock()
	//完成发送逻辑
	go sendproc(node)
	//完成接受逻辑
	go recvproc(node)

	sendMsg(userId, []byte("hello,world"))
}


func checkToken(userId int64, token string) bool{
	//从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token==token
}

func sendproc(node *Node){
	for{
		select{
			case data := <-node.DataQueue:
				err := node.Conn.WriteMessage(websocket.TextMessage,data)
				if err != nil{
					log.Println(err.Error())
					return
				}
		}
	}
}

func recvproc(node *Node){
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//对data做进一步的处理
		dispatch(data)
		//fmt.Printf("recv<=%s", data)
	}
}

//后端调度处理
func dispatch(data []byte){
	//解析data为message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err!=nil{
		log.Println(err.Error())
		return
	}
	//解析成功，根据cmd对逻辑进行处理
	switch msg.Cmd{
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid,data)
	case CMD_ROOM_MSG:
		//群聊
		for _,v:= range clientMap{
			if v.GroupSets.Has(msg.Dstid){
				v.DataQueue<-data
			}
		}
	case CMD_HEART:
			//心跳为了保证网络的持久性，因为有些服务器会把在一定时间内没有数据传输的网络进行关闭
			//什么都不做，只要接收到数据 链路就是正常的
	}
}

//发送消息
func sendMsg(dstId int64, msg []byte){
	rwlocker.RLock()
	node,ok := clientMap[dstId]
	//保证map的并发安全性
	rwlocker.RUnlock()
	if ok{
		node.DataQueue<-msg
	}
}

//todo 添加新的群ID到用户的groupset中
func AddGroupId(userId,gid int64){
	//取得node
	rwlocker.Lock()
	node ,ok := clientMap[userId]
	if ok{
		node.GroupSets.Add(gid)
	}
	rwlocker.Unlock()
	//添加gid到set
}

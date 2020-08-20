package gonet

import (
	"container/list"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type IWebSocketTask interface {
	ParseMsg(data []byte, flag byte) bool
	OnClose()
}

type WebSocketTask struct {
	closed      int32
	verified    bool
	stopedChan  chan bool
	sendMsgList *list.List
	sendMutex   sync.Mutex
	Conn        *websocket.Conn
	Derived     IWebSocketTask
	msgChan     chan []byte
	signal      chan int
}

func NewWebSocketTask(conn *websocket.Conn) *WebSocketTask {
	return &WebSocketTask{
		closed:      -1,
		verified:    false,
		Conn:        conn,
		stopedChan:  make(chan bool, 1),
		sendMsgList: list.New(),
		msgChan:     make(chan []byte, 1024),
		signal:      make(chan int, 1),
	}
}

func (this *WebSocketTask) Signal() {
	select {
	case this.signal <- 1:
	default:
	}
}

func (this *WebSocketTask) Start() {
	if atomic.CompareAndSwapInt32(&this.closed, -1, 0) {
		glog.Info("[WS Connect] Got Connect, ", this.Conn.RemoteAddr())
		go this.sendloop()
		go this.recvloop()
	}
}

func (this *WebSocketTask) Stop() bool {
	if !this.IsClosed() && len(this.stopedChan) == 0 {
		this.stopedChan <- true
	} else {
		glog.Info("[WS Connect] Stop Connect Fail ", len(this.stopedChan))
		return false
	}
	glog.Info("[WS Connect] Stop Connect Success")
	return true
}

func (this *WebSocketTask) Close() {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		this.Conn.Close()
		this.Derived.OnClose()
		close(this.stopedChan)

		glog.Info("[WS Connect] Connect Close ", this.Conn.RemoteAddr())
	}
}

func (this *WebSocketTask) Reset() {
	if atomic.LoadInt32(&this.closed) == 1 {
		glog.Info("[WS Connect] Connect Reset ", this.Conn.RemoteAddr())
		this.closed = -1
		this.verified = false
		this.stopedChan = make(chan bool)
	}
}

func (this *WebSocketTask) AsyncSend(buffer []byte, flag byte) bool {
	if this.IsClosed() {
		return false
	}

	bsize := len(buffer)
	totalsize := bsize + 4

	sendbuffer := make([]byte, 0, totalsize)
	sendbuffer = append(sendbuffer, byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	sendbuffer = append(sendbuffer, buffer...)
	this.msgChan <- sendbuffer

	return true
}

func (this *WebSocketTask) recvloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[Unexpeted] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var datasize int

	for !this.IsClosed() {
		_, bytemsg, err := this.Conn.ReadMessage()
		if nil != err {
			glog.Error("[WS] Recv Fail ", this.Conn.RemoteAddr(), ",", err)
			return
		}

		datasize = int(bytemsg[0]) | int(bytemsg[1])<<8 | int(bytemsg[2])<<16
		if datasize > cmd_max_size {
			glog.Error("[WS] Packet Too Large ", this.Conn.RemoteAddr(), ",", datasize)
			return
		}

		this.Derived.ParseMsg(bytemsg[cmd_header_size:], bytemsg[3])
	}
}

func (this *WebSocketTask) sendloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[Unexpeted] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var (
		timeout = time.NewTimer(time.Second * cmd_verify_time)
	)

	for {
		select {
		case bytemsg := <-this.msgChan:
			if nil != bytemsg && len(bytemsg) > 0 {
				err := this.Conn.WriteMessage(websocket.BinaryMessage, bytemsg)
				if nil != err {
					glog.Error("[WS] Send Fail ", this.Conn.RemoteAddr(), ",", err)
					return
				}
			} else {
				glog.Error("[WS] Wrong Msg! ", bytemsg)
				return
			}
		case <-this.stopedChan:
			return
		case <-timeout.C:
			// 使用验证通常是为了防止用户连接而不使用，占据服务器资源
			// 验证客户端是否合法后可以减少这种情况
			if !this.IsVerifed() {
				glog.Error("[WS] Verify Fail ", this.Conn.RemoteAddr())
				return
			}
		}
	}
}

func (this *WebSocketTask) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) != 0
}

func (this *WebSocketTask) Verify() {
	this.verified = true
}

func (this *WebSocketTask) IsVerifed() bool {
	return this.verified
}

func (this *WebSocketTask) Terminate() {
	this.Close()
}

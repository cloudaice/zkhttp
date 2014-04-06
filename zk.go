package zkhttp

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Zh struct {
	addrs  []string
	zkConn *zk.Conn
}

func NewZh(servers []string) *Zh {
	return &Zh{
		addrs: servers,
	}
}

func (zh *Zh) connect() (*zk.Conn, <-chan zk.Event, error) {
	zkConn, session, err := zk.Connect(
		zh.addrs,
		time.Duration(5)*time.Second)
	if err != nil {
		return nil, session, err
	}
	return zkConn, session, nil
}

// deal with session and return when state is disconnect
func (zh *Zh) dealSession(session <-chan zk.Event) {
	for e := range session {
		log.Printf("recieve zookeeper event %v\n", e.State)
		if e.State == zk.StateDisconnected {
			return
		}
	}
}

func (zh *Zh) ConnForever() {
	for {
		zkConn, session, err := zh.connect()
		if err != nil {
			log.Printf("connect zookeeper cause by %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		zh.zkConn = zkConn
		go zh.dealSession(session)

		checker := time.NewTicker(2 * time.Millisecond)
	CHECKSTATE:
		for _ = range checker.C {
			switch zh.zkConn.State() {
			case zk.StateConnected, zk.StateConnecting, zk.StateHasSession:
				continue
			default:
				log.Println("reconnect zookeeper")
				break CHECKSTATE
			}
		}
	}
}

func (zh *Zh) state(addrport string) (string, error) {
	conn, err := net.Dial("tcp", addrport)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	n, err := conn.Write([]byte("stat"))
	if err != nil {
		return "", err
	}
	if n != 4 {
		return "", errors.New("cmd stat return is not 4")
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (zh *Zh) State() {

}

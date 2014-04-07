package zkop

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
	zh := &Zh{
		addrs: servers,
	}
	if err := zh.connect(); err != nil {
		return nil
	}

	return zh
}

func (zh *Zh) connect() error {
	zkConn, session, err := zk.Connect(
		zh.addrs,
		time.Duration(5)*time.Second)
	if err != nil {
		return err
	}

	go zh.dealSession(session)
	zh.zkConn = zkConn
	return nil
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

func (zh *Zh) Close() {
	zh.zkConn.Close()
}

//Fetch state with single zookeeper addr and port
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

func (zh *Zh) State() []*ZStat {
	clusterState := []*ZStat{}
	for _, addrport := range zh.addrs {
		stat, err := zh.state(addrport)
		if err != nil {
			log.Println(err)
			return clusterState
		}
		zstat := NewStat(stat)
		zstat.Name = addrport
		clusterState = append(clusterState, zstat)
	}

	return clusterState
}

// list children of given znode
func (zh *Zh) Ls(znode string) ([]string, error) {
	var nodes []string
	children, _, err := zh.zkConn.Children(znode)
	if err != nil {
		return nodes, err
	}

	return children, nil
}

func (zh *Zh) Get(znode string) (string, error) {
	data, _, err := zh.zkConn.Get(znode)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (zh *Zh) Delete(znode string) error {
	return zh.zkConn.Delete(znode, -1)
}

func (zh *Zh) Create(znode, zdata string) error {
	_, err := zh.zkConn.Create(znode, []byte(zdata), 0, zk.WorldACL(zk.PermAll))
	return err
}

func (zh *Zh) Set(znode, zdata string) error {
	_, err := zh.zkConn.Set(znode, []byte(zdata), -1)
	return err
}

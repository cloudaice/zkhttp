package zkop

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type StatError struct {
	e string
}

func NewStateError(text string) *StatError {
	return &StatError{text}
}

func (err *StatError) Error() string {
	return "State Error: " + err.e
}

type ZStat struct {
	Name        string   // zookeeper 地址
	Conns       []string //各个连接详细信息
	Latency     string   //
	Received    int64    // 接收的字节数
	Sent        int64    // 发送的字节数
	Connections int64    // 连接数
	Outstanding int64    // 已经提交，未commit的事物
	Zxid        string   // 事物ID
	Mode        string   // 角色
	NodeCount   int64    // 节点个数
}

func NewStat(data string) *ZStat {
	stat := new(ZStat)
	s := strings.Split(data, "\n\n")
	index := strings.Index(s[0], "Clients:")
	if index != -1 {
		conns := strings.Split(strings.
			TrimSpace(s[0][index+9:]), "\n")
		stat.Conns = conns
	} else {
		log.Println(NewStateError("Parse Clients"))
	}

	lines := strings.Split(strings.TrimSpace(s[1]), "\n")
	if len(lines) != 8 {
		log.Println(NewStateError("Length is beyond 8"))
		return stat
	}

	stat.Latency = strings.TrimSpace(strings.Split(lines[0], ":")[1])

	if received, err := strconv.ParseInt(
		strings.TrimSpace(strings.Split(lines[1], ":")[1]), 10, 64); err == nil {
		stat.Received = received
	} else {
		log.Println(NewStateError(
			fmt.Sprintf("Parse Int64 Received %s", lines[1])))
	}

	if sents, err := strconv.ParseInt(
		strings.TrimSpace(strings.Split(lines[2], ":")[1]), 10, 64); err == nil {
		stat.Sent = sents
	} else {
		log.Println(NewStateError(
			fmt.Sprintf("Parse Int64 Sents", lines[2])))
	}

	if conns, err := strconv.ParseInt(
		strings.TrimSpace(strings.Split(lines[3], ":")[1]), 10, 64); err == nil {
		stat.Connections = conns
	} else {
		log.Println(NewStateError(
			fmt.Sprintf("Parse Int64 Connections", lines[3])))
	}

	if outs, err := strconv.ParseInt(
		strings.TrimSpace(strings.Split(lines[4], ":")[1]), 10, 64); err == nil {
		stat.Outstanding = outs
	} else {
		log.Println(NewStateError(
			fmt.Sprintf("Parse Int Outstanding", lines[4])))
	}

	stat.Zxid = strings.TrimSpace(strings.Split(lines[5], ":")[1])

	stat.Mode = strings.TrimSpace(strings.Split(lines[6], ":")[1])

	if nodes, err := strconv.ParseInt(
		strings.TrimSpace(strings.Split(lines[7], ":")[1]), 10, 64); err == nil {
		stat.NodeCount = nodes
	} else {
		log.Println(NewStateError(
			fmt.Sprintf("Parse Int64 NodeCount", lines[7])))
	}

	return stat
}

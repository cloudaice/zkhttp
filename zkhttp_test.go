package zkhttp

import (
	"strings"
	"testing"
)

func TestState(t *testing.T) {
	zkServers := []string{"10.237.36.153:2181", "10.237.36.154:2181", "10.237.36.155:2181"}
	zh := NewZh(zkServers)
	go zh.ConnForever()

	for _, zkAddr := range zkServers {
		result, err := zh.state(zkAddr)
		if err != nil {
			t.Error("Test State error ", err)
		}
		s := strings.Split(result, "\n\n")
		t.Log(len(s))
		t.Log(s[0])
		t.Log("***")
		t.Log(s[1])
	}
}

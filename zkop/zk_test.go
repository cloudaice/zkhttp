package zkop

import (
	"testing"
)

func TestState(t *testing.T) {
	zkServers := []string{"10.237.36.153:2181", "10.237.36.154:2181", "10.237.36.155:2181"}
	zh := NewZh(zkServers)
	go zh.ConnForever()

	for _, zkAddr := range zkServers {
		_, err := zh.state(zkAddr)
		if err != nil {
			t.Error("Test State error ", err)
		}
		//s := strings.Split(result, "\n\n")
		//t.Log(len(s))
		//t.Log(s[0])
		//t.Log("***")
		//t.Log(s[1])
		//t.Log(len(strings.Split(strings.TrimSpace(s[1]), "\n")))
	}
	t.Log("Test State Pass")

	for _, st := range zh.State() {
		t.Log("Name: ", st.Name)
		t.Log("Mode: ", st.Mode)
		t.Log("Connections: ", st.Connections)
		t.Log("Conn List: ", st.Conns)
		t.Log("Latency: ", st.Latency)
		t.Log("NodeCount: ", st.NodeCount)
		t.Log("Outstanding: ", st.Outstanding)
		t.Log("Received: ", st.Received)
		t.Log("Sent: ", st.Sent)
		t.Log("Zxid: ", st.Zxid)
	}

}

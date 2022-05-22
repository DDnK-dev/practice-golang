package main

import "time"

type NetstatInfo struct {
	Cluster    string
	Node       string
	Timestamp  int64
	NetstatStr string
	Containers []Container
}

type Container struct {
	ID         string
	PID        int32
	FullName   string
	NetstatStr string
}

func main() {
	netstatInfo := NetstatInfo{
		Cluster:    "cloudmoa",
		Node:       "imxc-master",
		Timestamp:  time.Now().Unix(),
		NetstatStr: "dummy dummy dummy",
		Containers: nil,
	}

	nodeMap = make(map[string]*NetstatInfo)
	if nodeMap == nil {
		nodeMap = make(map[string]*NetstatInfo)
		m.infoMap[cluster] = nodeMap
	}
	print(netstatInfo)
}

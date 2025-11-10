package util

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func Init(nodeId int64) (err error) {
	node, err = snowflake.NewNode(nodeId)
	return
}

func GetID() string {
	return node.Generate().String()
}

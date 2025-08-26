package idgen

import (
	"github.com/bwmarrin/snowflake"
)

type IDGenerator interface {
	NextId() int64
}

var global *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	global = node
}

func NextId() int64 {
	return global.Generate().Int64()
}

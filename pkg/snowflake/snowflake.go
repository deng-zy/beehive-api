package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node
var once sync.Once

func NewNode(node int64) (*snowflake.Node, error) {
	return snowflake.NewNode(node)
}

func Generate() uint64 {
	if node == nil {
		once.Do(func() {
			node, _ = NewNode(0)
		})
	}
	return uint64(node.Generate().Int64())
}

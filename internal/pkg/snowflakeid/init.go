package snowflakeid

import (
	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func InitSnowflake() (*snowflake.Node, error) {
	var err error
	Node, err = snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	return Node, nil
}

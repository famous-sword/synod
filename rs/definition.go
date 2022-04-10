package rs

import "fmt"

const (
	NumParityShard = 2
	NumDataShard   = 4
	TotalShards    = NumDataShard + NumParityShard

	BlockPerShard = 8000
	BlockSize     = BlockPerShard * NumDataShard
)

func getShardSeq(hash string, seq int) string {
	return fmt.Sprintf("%s.%d", hash, seq)
}

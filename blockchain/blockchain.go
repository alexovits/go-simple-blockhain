package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const difficulty = 1

type Block struct {
	Index      int
	Timestamp  string
	Data       int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func GenerateBlock(oldBlock Block, Data int) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(CalculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(CalculateHash(newBlock), " work done!")
			newBlock.Hash = CalculateHash(newBlock)
			break
		}

	}
	return newBlock
}

func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

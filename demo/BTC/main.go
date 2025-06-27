package BTC

import (
	"fmt"
	"github.com/qujing226/cases/demo/BTC/blockchain"
)

func main() {
	// 创建一个区块链
	chain := blockchain.NewChain()

	//创建新区块并添加到区块链
	chain.Add("转账十元")
	chain.Add("转账二十元")

	// 打印区块链
	fmt.Println(chain)
}

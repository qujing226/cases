package BTC

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ripemd160"
	"testing"
)

func TestGenerateBitcoinAddress(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "test",
			input:  "0fe57f1532728ea9f1891b0bce90ba3f9c3c64f0cda0e439e9c2fa56553014b9",
			output: "15gQHFAHyvSAMZZ7tZLfbNyngtJ6fQe3Hm",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := GenerateBitcoinAddress(tc.input)
			assert.Equal(t, tc.output, res)
		})
	}
}

func GenerateBitcoinAddress(privateKeyHex string) string {
	// 1. 解码私钥
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// 2. 生成压缩公钥
	pubKey := privateKey.PubKey()
	x := pubKey.X()
	y := pubKey.Y()

	var prefix byte
	if y.Bit(0) == 0 { // 根据 y 的奇偶性选择前缀
		prefix = 0x02
	} else {
		prefix = 0x03
	}

	// 将 x 填充到 32 字节
	xBytes := x.Bytes()
	paddedX := make([]byte, 32)
	copy(paddedX[32-len(xBytes):], xBytes) // ✅ 高位补零
	compressedPublicKey := append([]byte{prefix}, paddedX...)

	// 3. 计算 SHA-256 + RIPEMD-160
	sha256Hash := sha256.Sum256(compressedPublicKey)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// 4. 添加版本字节 0x00
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// 5. 计算校验和
	firstSHA := sha256.Sum256(versionedPayload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	// 6. Base58 编码
	finalPayload := append(versionedPayload, checksum...)
	bitcoinAddress := base58.Encode(finalPayload)
	return bitcoinAddress
}

// GenerateBitcoinAddress generates a compressed Bitcoin P2PKH address from a private key.
func GenerateBitcoinAddressV1(privateKeyHex string) string {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// 2. 生成压缩公钥
	compressedPublicKey := privateKey.PubKey().SerializeCompressed()

	// 3. 生成 P2PKH 地址（Legacy 格式）
	addressPubKeyHash, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(compressedPublicKey), &chaincfg.MainNetParams)
	return addressPubKeyHash.EncodeAddress()
}

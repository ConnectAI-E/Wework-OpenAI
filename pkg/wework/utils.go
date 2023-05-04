package wework

import (
	"crypto/sha1"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"sort"
)

func CheckSignature(token, signature string, timestamp int64, nonce string, echoStr string) bool {
	// 将 token、timestamp、nonce、echostr 按字典序排序
	strs := []string{token, strconv.FormatInt(timestamp, 10), nonce, echoStr}
	sort.Strings(strs)

	// 将排序后的四个参数拼接成一个字符串
	str := strings.Join(strs, "")

	// 使用 sha1 算法计算字符串的哈希值
	hash := sha1.New()
	hash.Write([]byte(str))
	sum := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hexStr := hex.EncodeToString(sum)

	// 将十六进制字符串和签名比对，判断消息的合法性
	return signature == hexStr
}

func EncrpytMsg(msg string) (string, error) {
	return "", nil
}

func DecryptMsg(encodingAESKey, msg string) (result *EncrpytMsgData, err error) {
    // 解密消息体
	aesKey := encodingAESKey
    aesKeyBytes, err := base64.StdEncoding.DecodeString(aesKey + "=")
    encryptedData, err := base64.StdEncoding.DecodeString(msg)

    block, err := aes.NewCipher(aesKeyBytes)
    iv := aesKeyBytes[:aes.BlockSize]
    mode := cipher.NewCBCDecrypter(block, iv)

    decryptedData := make([]byte, len(encryptedData))
    mode.CryptBlocks(decryptedData, encryptedData)

	// PKCS#7填充处理
    pad := int(decryptedData[len(decryptedData)-1])
    if pad < 1 || pad > 32 {
        pad = 0
    }

	data, receiveid := parseContent(decryptedData)

    return &EncrpytMsgData{
		Msg: string(data),
		ReceiveId: receiveid,
	}, nil
}

func parseContent(content []byte) (msg []byte, receiveid string) {
    // 去掉前16个字节
    content = content[16:]

    // 取出前4个字节表示消息长度
    msgLen := binary.BigEndian.Uint32(content[:4])
    msg = content[4:4+msgLen]

    // 剩余字节为 receiveid
    receiveid = string(content[4+msgLen:])
    return
}
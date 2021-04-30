package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// func main() {
// 	origData := []byte("Hello World") // 待加密的数据
// 	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
// 	log.Println("原文：", string(origData))
//
// 	log.Println("------------------ CBC模式 --------------------")
// 	encrypted := AesEncryptCBC(origData, key)
// 	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
// 	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
// 	decrypted := AesDecryptCBC(encrypted, key)
// 	log.Println("解密结果：", string(decrypted))
//
// 	log.Println("------------------ ECB模式 --------------------")
// 	encrypted = AesEncryptECB(origData, key)
// 	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
// 	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
// 	decrypted = AesDecryptECB(encrypted, key)
// 	log.Println("解密结果：", string(decrypted))
//
// 	log.Println("------------------ CFB模式 --------------------")
// 	encrypted = AesEncryptCFB(origData, key)
// 	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
// 	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
// 	decrypted = AesDecryptCFB(encrypted, key)
// 	log.Println("解密结果：", string(decrypted))
// }

type Aes struct {
	Key []byte
}

// AesEncryptCBC AES，CBC模式加密
func (a Aes) AesEncryptCBC(origData []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(a.Key)
	blockSize := block.BlockSize()                                // 获取秘钥块的长度
	origData = a.pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, a.Key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                       // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                    // 加密
	return encrypted
}

// AesDecryptCBC AES，CBC模式解密
func (a Aes) AesDecryptCBC(encrypted []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(a.Key)                              // 分组秘钥
	blockSize := block.BlockSize()                                // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, a.Key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                      // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                   // 解密
	decrypted = a.pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}

// pkcs5Padding
func (a Aes) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs5UnPadding
func (a Aes) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncryptECB AES，ECB模式加密
func (a Aes) AesEncryptECB(origData []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(a.generateKey())
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// AesDecryptECB AES，ECB模式解密
func (a Aes) AesDecryptECB(encrypted []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(a.generateKey())
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

// generateKey 创建一个密钥
func (a Aes) generateKey() (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, a.Key)
	for i := 16; i < len(a.Key); {
		for j := 0; j < 16 && i < len(a.Key); j, i = j+1, i+1 {
			genKey[j] ^= a.Key[i]
		}
	}
	return genKey
}

// AesEncryptCFB AES，CFB模式加密
func (a Aes) AesEncryptCFB(origData []byte) (encrypted []byte) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}

// AesDecryptCFB AES，CFB模式解密
func (a Aes) AesDecryptCFB(encrypted []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(a.Key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// 高级加密标准（Advanced Encryption Standard，AES) 又称Rijndael加密法
// AES中常见的有三种解决方案，分别为AES-128、AES-192和AES-256。如果采用真正的128位加密技术甚至256位加密技术，蛮力攻击要取得成功需要耗费相当长的时间
// AES 有五种加密模式：
// 电码本模式（Electronic Codebook Book (ECB)）、
// 密码分组链接模式（Cipher Block Chaining (CBC)）、
// 计算器模式（Counter (CTR)）、
// 密码反馈模式（Cipher FeedBack (CFB)）
// 输出反馈模式（Output FeedBack (OFB)）

//加密解密接口
type Encryption interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte, []byte) ([]byte, error)
}

//电码本模式
type ECB struct {
	Encryption
}

//ECB加密
func (e *ECB) Encrypt(src, key []byte) ([]byte, error) {
	cipher, _ := aes.NewCipher(e.generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}

//ECB解密
func (e *ECB) Decrypt(encrypted, key []byte) ([]byte, error) {

	cipher, _ := aes.NewCipher(e.generateKey(key))
	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], nil
}

func (e *ECB) generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

//密码分组链接模式
type CBC struct {
	Encryption
}

func (c *CBC) Encrypt(orig, key []byte) ([]byte, error) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	orig = c.PKCS7Padding(orig, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 创建数组
	cryted := make([]byte, len(orig))
	// 加密
	blockMode.CryptBlocks(cryted, orig)
	return cryted, nil
}

func (c *CBC) Decrypt(crytedByte, k []byte) ([]byte, error) {

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = c.PKCS7UnPadding(orig)
	return orig, nil
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func (c *CBC) PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func (c *CBC) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//计算器模式
type CTR struct {
	Encryption
}

//加密
func (c *CTR) Encrypt(plainText, key []byte) ([]byte, error) {
	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)
	return dst, nil
}

//解密
func (c *CTR) Decrypt(decrypt, key []byte) ([]byte, error) {
	return c.Encrypt(decrypt, key)
}

//密码反馈模式
type CFB struct {
	Encryption
}

func (c *CFB) Encrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted, nil
}

func (c *CFB) Decrypt(encrypted, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}

//密码反馈模式
type OFB struct {
	Encryption
}

func (c *OFB) Encrypt(data, key []byte) ([]byte, error) {
	data = c.PKCS7Padding(data, aes.BlockSize)
	block, _ := aes.NewCipher([]byte(key))
	out := make([]byte, aes.BlockSize+len(data))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], data)
	return out, nil
}

func (c *OFB) Decrypt(data, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(key))
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("data is not a multiple of the block size")
	}
	out := make([]byte, len(data))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(out, data)
	out = c.PKCS7UnPadding(out)
	return out, nil
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func (c *OFB) PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func (c *OFB) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

const (
	Ecb = 1 //电码本模式
	Cbc = 2 //密码分组链接模式
	Ctr = 3 //计算器模式
	Cfb = 4 //密码反馈模式
	Ofb = 5 //输出反馈模式
)

type AES struct {
	mode int //模式
	obj  Encryption
}

func NewAES() *AES {
	aec := &AES{}
	aec.obj = aec.getType()
	return aec
}

//加密
func (a *AES) Encrypt(cryted, key []byte) ([]byte, error) {
	return a.obj.Encrypt(cryted, key)
}

//解密
func (a *AES) Decrypt(cryted, key []byte) ([]byte, error) {
	return a.obj.Decrypt(cryted, key)
}

func (a *AES) SetMode(m int) {
	a.mode = m
	a.obj = a.getType()
}

func (a *AES) getType() Encryption {
	var obj Encryption
	switch a.mode {
	case 1:
		obj = &ECB{}
	case 2:
		obj = &CBC{}
	case 3:
		obj = &CTR{}
	case 4:
		obj = &CFB{}
	case 5:
		obj = &OFB{}
	default:
		obj = &CBC{}
	}
	a.obj = obj
	return obj
}

//AES加密调用
func Aes_en(src, key []byte, mod int) ([]byte, *AES) {
	aes := NewAES() // 默认是Cbc 模式
	aes.SetMode(mod)
	got, _ := aes.Encrypt(src, key)
	return got, aes

}

//AES解密调用
func Aes_de(src, key []byte, aes *AES) []byte {
	want, _ := aes.Decrypt(src, key)
	return want

}

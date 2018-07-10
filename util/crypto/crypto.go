package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const (
	byteLen      = 8
	AES128KeyLen = 128 / byteLen
	AES192KeyLen = 192 / byteLen
	AES256KeyLen = 256 / byteLen
)

var (
	// AesType : used to let the Encrypter/Decrypter know which type of AES to use
	AesType int = AES256KeyLen
)

func CFB_Enc(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(UnifyKeyLen(key, AesType))
	if err != nil {
		return nil, err
	}

	retData := make([]byte, len(data)+aes.BlockSize)
	iv := retData[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	encStream := cipher.NewCFBEncrypter(block, iv)
	encStream.XORKeyStream(retData[aes.BlockSize:], data)
	return retData, nil
}

func CFB_Dec(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(UnifyKeyLen(key, AES256KeyLen))
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	decStream := cipher.NewCFBDecrypter(block, iv)

	retData := make([]byte, len(data)-aes.BlockSize)
	decStream.XORKeyStream(retData, data[aes.BlockSize:])

	return retData, nil
}

func UnifyKeyLen(key []byte, targetLen int) []byte {
	ret := key
	keyLen := len(key)
	if keyLen < targetLen {
		pad := make([]byte, targetLen-keyLen)
		ret = append(key, pad...)
	} else if keyLen > targetLen {
		ret = key[:targetLen]
	}
	return ret
}

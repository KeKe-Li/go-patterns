package chapter04

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"github.com/KeKe-Li/log"
	"github.com/pkg/errors"
	"io"
	"time"
)

// ICrypto for crypto
type ICrypto interface {
	ParseParamEncrypto(ctx context.Context, encodingAESKey string, encryptMsg string) (*ResquestBody, error)
}

// ResquestBody for struct
type ResquestBody struct {
	XMLName      xml.Name      `xml:"xml"`
	ToUserName   string        `xml:"ToUserName"`
	FromUserName string        `xml:"FromUserName"`
	CreateTime   time.Duration `xml:"CreateTime"`
	MsgType      string        `xml:"MsgType"`
	Content      string        `xml:"Content"`
	MsgID        int           `xml:"MsgId"`
	Event        string        `xml:"Event"`
	EventKey     string        `xml:"EventKey"`
}

// NewCrypto for invoke
func NewCrypto(encodingAESKey string, encryptMsg string) ICrypto {
	return &crypto{
		EncodingAESKey: encodingAESKey,
		EncryptMsg:     encryptMsg,
		err:            nil,
	}
}

type crypto struct {
	err            error
	EncodingAESKey string `json:"encoding_aes_key"`
	EncryptMsg     string `json:"encrypt_msg"`
}

func (impl *crypto) resoleError(ctx context.Context) error {
	log.ErrorContext(ctx, "resoleError")
	return impl.err
}

func (impl *crypto) ParseParamEncrypto(ctx context.Context, encodingAESKey string, encryptMsg string) (*ResquestBody, error) {
	if encodingAESKey == "" {
		log.WarnContext(ctx, "encodingAESKey-is-not-found")
		return nil, errors.New("encodingAESKey is not found")
	}
	if encryptMsg == "" {
		log.WarnContext(ctx, "encryptMsg-is-not-found")
		return nil, errors.New("encryptMsg is not found")
	}
	var err error
	var encryptedMsg []byte
	encryptedMsg, err = base64.StdEncoding.DecodeString(encryptMsg)
	if err != nil {
		log.ErrorContext(ctx, "base64-StdEncoding-DecodeString-encryptMsg-failed", "error", err.Error())
		return nil, err
	}
	encodingAESKeyParam, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		log.ErrorContext(ctx, "base64-StdEncoding-DecodeString-encodingAESKey-failed", "error", err.Error())
		return nil, err
	}
	// AES Decrypt
	plainData, err := impl.aesDecrypt(ctx, encryptedMsg, encodingAESKeyParam)
	if err != nil {
		log.ErrorContext(ctx, "aesDecrypt-failed", "error", err.Error())
		return nil, err
	}
	textRequestBody, err := impl.parseEncryptTextRequestBody(ctx, plainData)
	if err != nil {
		log.ErrorContext(ctx, "parseEncryptTextRequestBody-failed", "error", err.Error())
		return nil, err
	}
	return textRequestBody, nil
}

func (impl *crypto) parseEncryptTextRequestBody(ctx context.Context, plainText []byte) (*ResquestBody, error) {
	// Read Length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		log.ErrorContext(ctx, "binary-Read-failed", "error", err.Error())
		return nil, err
	}
	var textRequestBody ResquestBody
	err = xml.Unmarshal(plainText[20:20+length], &textRequestBody)
	if err != nil {
		log.ErrorContext(ctx, "xml-Unmarshal", "error", err.Error())
		return nil, err
	}
	return &textRequestBody, nil
}

func (impl *crypto) aesDecrypt(ctx context.Context, cipherData []byte, aesKey []byte) ([]byte, error) {
	//PKCS#7
	a := len(aesKey)
	if len(cipherData)%a != 0 {
		return nil, errors.New("ciphertext size is not multiple of aes key length")
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		log.ErrorContext(ctx, "aes-NewCipher-failed", "error", err.Error())
		return nil, err
	}
	abs := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, abs); err != nil {
		log.ErrorContext(ctx, "ioutil-ReadAll-failed", "error", err.Error())
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, abs)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

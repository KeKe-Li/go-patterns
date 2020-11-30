package chater04

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/KeKe-Li/log"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Token          string `json:"token"`
	AppID          string `json:"app_id"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

// this is important for ToUserName and FromUserName, for wechat return ,we should use  Merchant FromUserName to ToUserName return
// and FromUserName use wechat ToUserName return .
type WechatRequestData struct {
	ToUserName string `xml:"ToUserName"`
	//Encrypt    string        `xml:"Encrypt"`
	FromUserName string        `xml:"FromUserName"`
	CreateTime   time.Duration `xml:"CreateTime"`
	MsgType      string        `xml:"MsgType"`
	MsgID        int64         `xml:"MsgId"`
	Event        string        `xml:"Event"`
	EventKey     string        `xml:"EventKey"`
	MediaID      string        `xml:"MediaId"`
	Content      string        `xml:"Content"`
	PicURL       string        `xml:"PicUrl"`
	Nonce        string        `xml:"Nonce"`
}

type IDecrypt interface {
	EncryptResponseData(ctx context.Context, input *WechatRequestData) (string, error)
}

type decrypt struct {
	Config Config
}

func NewDecrypt(config Config) IDecrypt {
	return &decrypt{
		Config: config,
	}
}

type EncryptTextResponseData struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   string
	MsgType      CDATAText
	Content      CDATAText
}

type EncryptResponseData struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    string
	Nonce        CDATAText
}

// CDATAText for aes
type CDATAText struct {
	Value string `xml:",cdata"`
}

func (impl *decrypt) Value2CDATA(s string) CDATAText {
	return CDATAText{Value: "<![CDATA[" + s + "]]>"}
}

// Signature for aes
func (impl *decrypt) Signature(ctx context.Context, timestamp, nonce, msgEncrypt string) string {
	sl := []string{impl.Config.Token, timestamp, nonce, msgEncrypt}
	sort.Strings(sl)
	s := sha1.New()
	_, err := io.WriteString(s, strings.Join(sl, ""))
	if err != nil {
		log.ErrorContext(ctx, "Signature-failed", "error", "error", err.Error())
		return ""
	}
	return fmt.Sprintf("%x", s.Sum(nil))
}

func (impl *decrypt) EncryptResponseData(ctx context.Context, input *WechatRequestData) (string, error) {
	if input == nil {
		log.WarnContext(ctx, "EncryptResponseData-input-is-nil")
		return "", errors.New("WechatRequestData is not found")
	}
	aesKey, err := base64.StdEncoding.DecodeString(impl.Config.EncodingAESKey + "=")
	if err != nil {
		log.ErrorContext(ctx, "base64-StdEncoding-DecodeString-failed", "error", err.Error())
		return "", err
	}
	encrypt, err := impl.encryptXMLData(ctx, input, aesKey)
	if err != nil {
		log.ErrorContext(ctx, "encryptXMLData-DecodeString-failed", "error", err.Error())
		return "", err
	}
	timeStamp := strconv.FormatInt(int64(input.CreateTime), 10)
	data := EncryptResponseData{
		Encrypt:      impl.Value2CDATA(encrypt),
		MsgSignature: impl.Value2CDATA(impl.Signature(ctx, timeStamp, input.Nonce, encrypt)),
		TimeStamp:    timeStamp,
		Nonce:        impl.Value2CDATA(input.Nonce),
	}
	result, err := xml.MarshalIndent(data, " ", " ")
	if err != nil {
		log.ErrorContext(ctx, "xml-MarshalIndent-failed", "error", err.Error())
		return "", err
	}

	return string(result), nil
}

func (impl *decrypt) encryptXMLData(ctx context.Context, input *WechatRequestData, aesKey []byte) (string, error) {
	createTime := strconv.FormatInt(int64(input.CreateTime), 10)
	encryptData := &EncryptTextResponseData{
		ToUserName:   impl.Value2CDATA(input.ToUserName),
		FromUserName: impl.Value2CDATA(input.FromUserName),
		CreateTime:   createTime,
		MsgType:      impl.Value2CDATA(input.ToUserName),
		Content:      impl.Value2CDATA(input.ToUserName),
	}
	xm, err := xml.MarshalIndent(encryptData, " ", " ")
	if err != nil {
		log.ErrorContext(ctx, "xml.MarshalIndent-failed", "error", err.Error())
		return "", err
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(xm)))
	if err != nil {
		log.ErrorContext(ctx, "binary.Write-failed", "error", err.Error())
		return "", err
	}
	bodyLength := buf.Bytes()
	randomBytes := []byte("abcdefghijklmnop")
	plainData := bytes.Join([][]byte{randomBytes, bodyLength, xm, []byte(impl.Config.AppID)}, nil)
	cipherData, err := impl.AesEncrypt(ctx, plainData, aesKey)
	if err != nil {
		log.ErrorContext(ctx, "AesEncrypt-failed", "error", err.Error())
		return "", errors.New("AesEncrypt error")
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// AesEncrypt for aes
func (impl *decrypt) AesEncrypt(ctx context.Context, plainData []byte, aesKey []byte) ([]byte, error) {
	aky := len(aesKey)
	if len(plainData)%aky != 0 {
		plainData = impl.PKCS7Pad(plainData, aky)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		log.ErrorContext(ctx, "aes.NewCipher-failed", "error", err.Error())
		return nil, err
	}
	abs := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, abs); err != nil {
		log.ErrorContext(ctx, "io.ReadFull-failed", "error", err.Error())
		return nil, err
	}
	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, abs)
	blockMode.CryptBlocks(cipherData, plainData)
	return cipherData, nil
}

// PKCS7Pad for aes
func (impl *decrypt) PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := PadLength(len(message), blocksize)
		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
		// apply padding
		padded = append(message, padding...)
		return padded
	}
	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}

// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(sliceLength, blocksize int) (padlen int) {
	padlen = blocksize - sliceLength%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}

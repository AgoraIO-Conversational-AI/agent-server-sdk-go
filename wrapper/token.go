package wrapper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"sort"
	"time"
)

const (
	RolePublisher  = 1
	RoleSubscriber = 2

	DefaultExpirySeconds = 3600

	tokenVersion = "007"
)

type GenerateTokenOptions struct {
	AppID          string
	AppCertificate string
	Channel        string
	UID            uint32
	Role           int
	ExpirySeconds  int
}

func packUint16(val uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, val)
	return buf
}

func packUint32(val uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	return buf
}

func packString(val string) []byte {
	encoded := []byte(val)
	result := packUint16(uint16(len(encoded)))
	result = append(result, encoded...)
	return result
}

func packMap(m map[uint16]uint32) []byte {
	keys := make([]uint16, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	result := packUint16(uint16(len(m)))
	for _, k := range keys {
		result = append(result, packUint16(k)...)
		result = append(result, packUint32(m[k])...)
	}
	return result
}

func GenerateRtcToken(opts GenerateTokenOptions) (string, error) {
	if opts.AppID == "" {
		return "", fmt.Errorf("app_id is required")
	}
	if opts.AppCertificate == "" {
		return "", fmt.Errorf("app_certificate is required")
	}
	if opts.Channel == "" {
		return "", fmt.Errorf("channel is required")
	}
	if opts.ExpirySeconds <= 0 {
		opts.ExpirySeconds = DefaultExpirySeconds
	}
	if opts.Role == 0 {
		opts.Role = RolePublisher
	}

	now := uint32(time.Now().Unix())
	privilegeExpireTs := now + uint32(opts.ExpirySeconds)

	privileges := map[uint16]uint32{
		1: privilegeExpireTs,
		2: privilegeExpireTs,
	}

	message := packUint32(0)
	message = append(message, packUint32(now)...)
	message = append(message, packUint32(privilegeExpireTs)...)
	message = append(message, packMap(privileges)...)

	signingKey := hmacSHA256([]byte(opts.AppCertificate), []byte(opts.AppID))

	content := packString(opts.AppID)
	content = append(content, packUint32(now)...)
	content = append(content, message...)

	signature := hmacSHA256(signingKey, content)

	contentB64 := base64.StdEncoding.EncodeToString(content)
	signatureB64 := base64.StdEncoding.EncodeToString(signature)
	messageB64 := base64.StdEncoding.EncodeToString(message)

	return fmt.Sprintf("%s%s%s%s", tokenVersion, signatureB64, contentB64, messageB64), nil
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

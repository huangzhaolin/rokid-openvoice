package auth

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	pb "pb"

	"golang.org/x/net/context"
)

var (
	key            = "PLEASE_FIX_IT"
	device_type_id = "PLEASE_FIX_IT"
	device_id      = "PLEASE_FIX_IT"
	secret         = "PLEASE_FIX_IT"
)

type SpeechCredential struct {
	service string
	version string
}

func NewSpeechCredential(service, version string) *SpeechCredential {
	return &SpeechCredential{
		service: service, version: version,
	}
}

func (c SpeechCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return c.genAuthMap(), nil
}

func (c SpeechCredential) RequireTransportSecurity() bool {
	return false
}

func (c SpeechCredential) genAuthMap() map[string]string {
	keys := []string{
		"key",
		"device_type_id",
		"device_id",
		"service",
		"version",
		"time",
	}

	vals := []string{
		key,
		device_type_id,
		device_id,
		c.service,
		c.version,
		strconv.FormatInt(time.Now().Unix(), 10),
	}

	amap := make(map[string]string)
	str := ""
	for n, k := range keys {
		str = str + k + "=" + vals[n] + "&"
		amap[k] = vals[n]
	}
	str = str + "secret=" + secret
	amap["sign"] = fmt.Sprintf("%X", md5.Sum([]byte(str)))

	return amap
}

func (c SpeechCredential) Auth(client pb.SpeechClient, ctx context.Context) (int32, error) {
	var response *pb.AuthResponse
	var err error
	if response, err = client.Auth(ctx, &pb.AuthRequest{
		Key:          key,
		DeviceTypeId: device_type_id,
		DeviceId:     device_id,
		Service:      c.service,
		Version:      c.version,
		Timestamp:    strconv.FormatInt(time.Now().Unix(), 10),
		Sign:         c.genAuthMap()["sign"],
	}); err == nil {
		return response.Result, nil
	}

	return -1, err
}

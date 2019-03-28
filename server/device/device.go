package device

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

var ErrDeviceNotFound = errors.New("Device not found")

func newToken() string {
	var t [16]byte
	if _, err := rand.Read(t[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(t[:])
}

// Gets the current device from the request.
// May return ErrDeviceNotFound in case the device is not found.
func GetDevice(ctx context.Context, txn *database.DatabaseTxn) (*database.Device, error) {
	c := server.GetApiContext(ctx)
	k, err := c.GetCookie("SYZOJDEVICE")
	if err == http.ErrNoCookie {
		return nil, ErrDeviceNotFound
	} else if err != nil {
		return nil, err
	} else if len(k.Value) != 48 {
		return nil, ErrDeviceNotFound
	}
	dev, err := txn.GetDevice(ctx, database.DeviceRef(k.Value[:16]))
	if err != nil {
		return nil, err
	} else if dev == nil || dev.Info == nil {
		return nil, ErrDeviceNotFound
	}
	info := dev.Info
	token := k.Value[16:48]
	if info.Token == nil || *info.Token != token { // || info.UserAgent == nil || *info.UserAgent != c.GetHeader("User-Agent") {
		return nil, ErrDeviceNotFound
	}
	return dev, nil
}

// Gets the current device from the request.
// Creates a new one in case it doesn't exist.
func NewDevice(ctx context.Context, txn *database.DatabaseTxn) (*database.Device, error) {
	d, err := GetDevice(ctx, txn)
	if err == nil {
		return d, nil
	} else if err != ErrDeviceNotFound {
		return nil, err
	}
	c := server.GetApiContext(ctx)
	token := newToken()
	d = &database.Device{
		Info: &model.DeviceInfo{
			Token:      proto.String(token),
			UserAgent:  proto.String(c.GetHeader("User-Agent")),
			RemoteAddr: proto.String(c.GetRemoteAddr()),
		},
	}
	if err = txn.InsertDevice(ctx, d); err != nil {
		return nil, err
	}
	k := &http.Cookie{
		Name:    "SYZOJDEVICE",
		Value:   string(d.GetId()) + token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 365),
	}
	c.SetCookie(k)
	return d, nil
}

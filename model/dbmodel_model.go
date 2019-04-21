package model

import (
	"database/sql/driver"
	"errors"

	"github.com/golang/protobuf/proto"
)

var ErrInvalidType = errors.New("Can only scan []byte into protobuf message")

func (m *DeviceInfo) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return proto.Marshal(m)
}

func (m *DeviceInfo) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}

func (m *Problem) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return proto.Marshal(m)
}

func (m *Problem) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}

func (m *Problemset) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return proto.Marshal(m)
}

func (m *Problemset) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}

func (m *Submission) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return proto.Marshal(m)
}

func (m *Submission) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}

func (m *UserAuth) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return proto.Marshal(m)
}

func (m *UserAuth) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}

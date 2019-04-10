package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Database struct {
	db  *sql.DB
	ctx context.Context // context for all database operations

	m       sync.Mutex
	waiters map[interface{}][]chan struct{}
	cache   map[interface{}]cacheEntry
	wg      sync.WaitGroup
}

type cacheEntry struct {
	curData  interface{}
	prevData interface{}
}

var log = logrus.StandardLogger()

func Open(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	d := &Database{db: db}
	d.ctx = context.Background()
	d.waiters = make(map[interface{}][]chan struct{})
	d.cache = make(map[interface{}]cacheEntry)
	d.wg.Add(1)
	return d, nil
}

func (d *Database) Close() error {
	d.wg.Done()
	d.wg.Wait()
	return d.db.Close()
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *Database) getWaiter(key interface{}) chan struct{} {
	c := make(chan struct{}, 1)
	q := append(d.waiters[key], c)
	d.waiters[key] = q
	if len(q) == 1 {
		select {
		case c <- struct{}{}:
		default:
		}
	}
	return c
}

func (d *Database) done(key interface{}) {
	q := d.waiters[key]
	if len(q) == 0 {
		panic("done: empty waiter queue")
	}
	q = q[1:]
	if len(q) == 0 {
		delete(d.waiters, key)
	} else {
		select {
		case q[0] <- struct{}{}:
		default:
		}
		d.waiters[key] = q
	}
}

func (d *Database) wait(ctx context.Context, key interface{}, waiter chan struct{}) error {
	select {
	case <-ctx.Done():
		go func() {
			<-waiter
			d.m.Lock()
			d.done(key)
			d.m.Unlock()
		}()
		return ctx.Err()
	case <-waiter:
		return nil
	}
}

func (d *Database) optWait(ctx context.Context, key interface{}) error {
	w := d.getWaiter(key)
	if len(d.waiters[key]) == 1 {
		<-w
		return nil
	}
	d.m.Unlock()
	if err := d.wait(ctx, key, w); err != nil {
		return err
	}
	d.m.Lock()
	return nil
}

func newId() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b[:])
}

type TimestampType = timestamp.Timestamp

func TimestampNow() *TimestampType {
	return ptypes.TimestampNow()
}

func Timestamp(ts *TimestampType) (time.Time, error) {
	return ptypes.Timestamp(ts)
}

func TimestampProto(t time.Time) (*TimestampType, error) {
	return ptypes.TimestampProto(t)
}

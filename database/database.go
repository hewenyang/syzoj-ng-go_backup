package database

import (
	"context"
    "crypto/rand"
	"database/sql"
    "encoding/base64"
	"sync/atomic"
    "reflect"

	"github.com/sirupsen/logrus"
)

type Database struct {
	db *sql.DB
}

type DatabaseTxn struct {
	tx   *sql.Tx
	done int32
}

var log = logrus.StandardLogger()

func Open(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) OpenTxn(ctx context.Context) (*DatabaseTxn, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	tx2 := &DatabaseTxn{tx: tx}
	go func() {
		<-ctx.Done()
		if atomic.LoadInt32(&tx2.done) == 0 {
			log.Warning("Detected a transaction that wasn't closed before context done")
		}
	}()
	return tx2, nil
}

func (d *Database) OpenReadonlyTxn(ctx context.Context) (*DatabaseTxn, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	tx2 := &DatabaseTxn{tx: tx}
	go func() {
		<-ctx.Done()
		if atomic.LoadInt32(&tx2.done) == 0 {
			log.Warning("Detected a transaction that wasn't closed before context done")
		}
	}()
	return tx2, nil
}

func (t *DatabaseTxn) Commit(context.Context) error {
	atomic.StoreInt32(&t.done, 1)
	return t.tx.Commit()
}

func (t *DatabaseTxn) Rollback() {
    if atomic.CompareAndSwapInt32(&t.done, 0, 1) {
        go func() {
            err := t.tx.Rollback()
            if err != nil {
                log.WithError(err).Error("Failed to rollback transaction")
            }
        }()
    }
}

func (t *DatabaseTxn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *DatabaseTxn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func newId() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b[:])
}

func ScanAll(r *sql.Rows, v interface{}) error {
    val := reflect.ValueOf(v)
    if val.Type().Kind() != reflect.Ptr {
        panic("database: ScanAll: must be a pointer to a slice")
    }
    val = val.Elem()
    if val.Type().Kind() != reflect.Slice {
        panic("database: ScanAll: must be a pointer to a slice")
    }
    slice := reflect.Zero(val.Type())
    elType := val.Type().Elem()
    var i int
    for r.Next() {
        slice = reflect.Append(slice, reflect.Zero(elType))
        err := r.Scan(slice.Index(i).Addr().Interface())
        if err != nil {
            return err
        }
    }
    if err := r.Err(); err != nil {
        return err
    }
    val.Set(slice)
    return nil
}

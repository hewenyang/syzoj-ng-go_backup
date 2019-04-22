package database

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var _ = bytes.NewBuffer
var _ = time.Now
var _ = jsonpb.Unmarshal
var _ = proto.Marshal
var _ = ptypes.Duration
var _ = any.Any{}
var _ = structpb.Struct{}
var _ = timestamp.Timestamp{}

func convertTimestamp(t *timestamp.Timestamp) interface{} {
	if t == nil {
		return nil
	}
	t2, _ := ptypes.Timestamp(t)
	return t2
}

var jsonpbMarshaler = &jsonpb.Marshaler{}

func convertStruct(t *structpb.Struct) interface{} {
	if t == nil {
		return nil
	}
	var v bytes.Buffer
	err := jsonpbMarshaler.Marshal(&v, t)
	if err != nil {
		panic(err)
	}
	return v.Bytes()
}

func convertAny(t *any.Any) interface{} {
	if t == nil {
		return nil
	}
	v, err := proto.Marshal(t)
	if err != nil {
		panic(err)
	}
	return v
}

type UserRef string

func NewUserRef() UserRef {
	return UserRef(newId())
}

func CreateUserRef(ref UserRef) *UserRef {
	x := ref
	return &x
}

func (d *Database) getUser(ctx context.Context, ref UserRef) (*User, error) {
	v := new(User)
	var var3 *time.Time
	err := d.QueryRowContext(ctx, "SELECT id, user_name, auth, register_time FROM user WHERE id=?", ref).Scan(&v.Id, &v.UserName, &v.Auth, &var3)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if var3 != nil {
		v.RegisterTime, _ = ptypes.TimestampProto(*var3)
	} else {
		v.RegisterTime = nil
	}
	return v, nil
}

func (d *Database) updateUser(ctx context.Context, ref UserRef, v *User) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE user SET user_name=?, auth=?, register_time=? WHERE id=?", v.UserName, v.Auth, convertTimestamp(v.RegisterTime), v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update User")
	}
}

func (d *Database) insertUser(ctx context.Context, v *User) {
	_, err := d.ExecContext(ctx, "INSERT INTO user (id, user_name, auth, register_time) VALUES (?, ?, ?, ?)", v.Id, v.UserName, v.Auth, convertTimestamp(v.RegisterTime))
	if err != nil {
		log.WithError(err).Error("Failed to insert User")
	}
}

func (d *Database) deleteUser(ctx context.Context, ref UserRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM user WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete User")
	}
}

func (d *Database) GetUser(ctx context.Context, ref UserRef) (*User, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*User), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*User), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getUser(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*User), nil
}

func (d *Database) UpdateUser(ctx context.Context, ref UserRef, updater func(*User) *User) (*User, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getUser(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*User))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushUser(d.ctx, ref)
	})
	return entry.curData.(*User), nil
}

func (d *Database) FlushUser(ctx context.Context, ref UserRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*User)
	curData := entry.curData.(*User)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertUser(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteUser(d.ctx, ref)
		} else {
			d.updateUser(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertUser(ctx context.Context, v *User) error {
	if v.Id == nil {
		v.Id = CreateUserRef(NewUserRef())
	}
	_, err := d.UpdateUser(ctx, v.GetId(), func(p *User) *User {
		if p != nil {
			panic("database.InsertUser: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteUser(ctx context.Context, ref UserRef) error {
	_, err := d.UpdateUser(ctx, ref, func(p *User) *User {
		return nil
	})
	return err
}

type DeviceRef string

func NewDeviceRef() DeviceRef {
	return DeviceRef(newId())
}

func CreateDeviceRef(ref DeviceRef) *DeviceRef {
	x := ref
	return &x
}

func (d *Database) getDevice(ctx context.Context, ref DeviceRef) (*Device, error) {
	v := new(Device)

	err := d.QueryRowContext(ctx, "SELECT id, user, info FROM device WHERE id=?", ref).Scan(&v.Id, &v.User, &v.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (d *Database) updateDevice(ctx context.Context, ref DeviceRef, v *Device) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE device SET user=?, info=? WHERE id=?", v.User, v.Info, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Device")
	}
}

func (d *Database) insertDevice(ctx context.Context, v *Device) {
	_, err := d.ExecContext(ctx, "INSERT INTO device (id, user, info) VALUES (?, ?, ?)", v.Id, v.User, v.Info)
	if err != nil {
		log.WithError(err).Error("Failed to insert Device")
	}
}

func (d *Database) deleteDevice(ctx context.Context, ref DeviceRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM device WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete Device")
	}
}

func (d *Database) GetDevice(ctx context.Context, ref DeviceRef) (*Device, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*Device), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*Device), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getDevice(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*Device), nil
}

func (d *Database) UpdateDevice(ctx context.Context, ref DeviceRef, updater func(*Device) *Device) (*Device, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getDevice(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*Device))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushDevice(d.ctx, ref)
	})
	return entry.curData.(*Device), nil
}

func (d *Database) FlushDevice(ctx context.Context, ref DeviceRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*Device)
	curData := entry.curData.(*Device)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertDevice(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteDevice(d.ctx, ref)
		} else {
			d.updateDevice(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertDevice(ctx context.Context, v *Device) error {
	if v.Id == nil {
		v.Id = CreateDeviceRef(NewDeviceRef())
	}
	_, err := d.UpdateDevice(ctx, v.GetId(), func(p *Device) *Device {
		if p != nil {
			panic("database.InsertDevice: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteDevice(ctx context.Context, ref DeviceRef) error {
	_, err := d.UpdateDevice(ctx, ref, func(p *Device) *Device {
		return nil
	})
	return err
}

type ProblemRef string

func NewProblemRef() ProblemRef {
	return ProblemRef(newId())
}

func CreateProblemRef(ref ProblemRef) *ProblemRef {
	x := ref
	return &x
}

func (d *Database) getProblem(ctx context.Context, ref ProblemRef) (*Problem, error) {
	v := new(Problem)
	var var3 *time.Time
	err := d.QueryRowContext(ctx, "SELECT id, problemset, user, create_time, problem FROM problem WHERE id=?", ref).Scan(&v.Id, &v.Problemset, &v.User, &var3, &v.Problem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if var3 != nil {
		v.CreateTime, _ = ptypes.TimestampProto(*var3)
	} else {
		v.CreateTime = nil
	}
	return v, nil
}

func (d *Database) updateProblem(ctx context.Context, ref ProblemRef, v *Problem) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE problem SET problemset=?, user=?, create_time=?, problem=? WHERE id=?", v.Problemset, v.User, convertTimestamp(v.CreateTime), v.Problem, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Problem")
	}
}

func (d *Database) insertProblem(ctx context.Context, v *Problem) {
	_, err := d.ExecContext(ctx, "INSERT INTO problem (id, problemset, user, create_time, problem) VALUES (?, ?, ?, ?, ?)", v.Id, v.Problemset, v.User, convertTimestamp(v.CreateTime), v.Problem)
	if err != nil {
		log.WithError(err).Error("Failed to insert Problem")
	}
}

func (d *Database) deleteProblem(ctx context.Context, ref ProblemRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM problem WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete Problem")
	}
}

func (d *Database) GetProblem(ctx context.Context, ref ProblemRef) (*Problem, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*Problem), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*Problem), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getProblem(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*Problem), nil
}

func (d *Database) UpdateProblem(ctx context.Context, ref ProblemRef, updater func(*Problem) *Problem) (*Problem, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getProblem(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*Problem))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushProblem(d.ctx, ref)
	})
	return entry.curData.(*Problem), nil
}

func (d *Database) FlushProblem(ctx context.Context, ref ProblemRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*Problem)
	curData := entry.curData.(*Problem)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertProblem(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteProblem(d.ctx, ref)
		} else {
			d.updateProblem(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertProblem(ctx context.Context, v *Problem) error {
	if v.Id == nil {
		v.Id = CreateProblemRef(NewProblemRef())
	}
	_, err := d.UpdateProblem(ctx, v.GetId(), func(p *Problem) *Problem {
		if p != nil {
			panic("database.InsertProblem: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteProblem(ctx context.Context, ref ProblemRef) error {
	_, err := d.UpdateProblem(ctx, ref, func(p *Problem) *Problem {
		return nil
	})
	return err
}

type ProblemsetRef string

func NewProblemsetRef() ProblemsetRef {
	return ProblemsetRef(newId())
}

func CreateProblemsetRef(ref ProblemsetRef) *ProblemsetRef {
	x := ref
	return &x
}

func (d *Database) getProblemset(ctx context.Context, ref ProblemsetRef) (*Problemset, error) {
	v := new(Problemset)

	err := d.QueryRowContext(ctx, "SELECT id, user, problemset FROM problemset WHERE id=?", ref).Scan(&v.Id, &v.User, &v.Problemset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (d *Database) updateProblemset(ctx context.Context, ref ProblemsetRef, v *Problemset) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE problemset SET user=?, problemset=? WHERE id=?", v.User, v.Problemset, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Problemset")
	}
}

func (d *Database) insertProblemset(ctx context.Context, v *Problemset) {
	_, err := d.ExecContext(ctx, "INSERT INTO problemset (id, user, problemset) VALUES (?, ?, ?)", v.Id, v.User, v.Problemset)
	if err != nil {
		log.WithError(err).Error("Failed to insert Problemset")
	}
}

func (d *Database) deleteProblemset(ctx context.Context, ref ProblemsetRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM problemset WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete Problemset")
	}
}

func (d *Database) GetProblemset(ctx context.Context, ref ProblemsetRef) (*Problemset, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*Problemset), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*Problemset), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getProblemset(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*Problemset), nil
}

func (d *Database) UpdateProblemset(ctx context.Context, ref ProblemsetRef, updater func(*Problemset) *Problemset) (*Problemset, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getProblemset(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*Problemset))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushProblemset(d.ctx, ref)
	})
	return entry.curData.(*Problemset), nil
}

func (d *Database) FlushProblemset(ctx context.Context, ref ProblemsetRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*Problemset)
	curData := entry.curData.(*Problemset)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertProblemset(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteProblemset(d.ctx, ref)
		} else {
			d.updateProblemset(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertProblemset(ctx context.Context, v *Problemset) error {
	if v.Id == nil {
		v.Id = CreateProblemsetRef(NewProblemsetRef())
	}
	_, err := d.UpdateProblemset(ctx, v.GetId(), func(p *Problemset) *Problemset {
		if p != nil {
			panic("database.InsertProblemset: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteProblemset(ctx context.Context, ref ProblemsetRef) error {
	_, err := d.UpdateProblemset(ctx, ref, func(p *Problemset) *Problemset {
		return nil
	})
	return err
}

type SubmissionRef string

func NewSubmissionRef() SubmissionRef {
	return SubmissionRef(newId())
}

func CreateSubmissionRef(ref SubmissionRef) *SubmissionRef {
	x := ref
	return &x
}

func (d *Database) getSubmission(ctx context.Context, ref SubmissionRef) (*Submission, error) {
	v := new(Submission)

	err := d.QueryRowContext(ctx, "SELECT id, user, problemset, problem, submission FROM submission WHERE id=?", ref).Scan(&v.Id, &v.User, &v.Problemset, &v.Problem, &v.Submission)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (d *Database) updateSubmission(ctx context.Context, ref SubmissionRef, v *Submission) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE submission SET user=?, problemset=?, problem=?, submission=? WHERE id=?", v.User, v.Problemset, v.Problem, v.Submission, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Submission")
	}
}

func (d *Database) insertSubmission(ctx context.Context, v *Submission) {
	_, err := d.ExecContext(ctx, "INSERT INTO submission (id, user, problemset, problem, submission) VALUES (?, ?, ?, ?, ?)", v.Id, v.User, v.Problemset, v.Problem, v.Submission)
	if err != nil {
		log.WithError(err).Error("Failed to insert Submission")
	}
}

func (d *Database) deleteSubmission(ctx context.Context, ref SubmissionRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM submission WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete Submission")
	}
}

func (d *Database) GetSubmission(ctx context.Context, ref SubmissionRef) (*Submission, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*Submission), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*Submission), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getSubmission(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*Submission), nil
}

func (d *Database) UpdateSubmission(ctx context.Context, ref SubmissionRef, updater func(*Submission) *Submission) (*Submission, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getSubmission(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*Submission))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushSubmission(d.ctx, ref)
	})
	return entry.curData.(*Submission), nil
}

func (d *Database) FlushSubmission(ctx context.Context, ref SubmissionRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*Submission)
	curData := entry.curData.(*Submission)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertSubmission(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteSubmission(d.ctx, ref)
		} else {
			d.updateSubmission(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertSubmission(ctx context.Context, v *Submission) error {
	if v.Id == nil {
		v.Id = CreateSubmissionRef(NewSubmissionRef())
	}
	_, err := d.UpdateSubmission(ctx, v.GetId(), func(p *Submission) *Submission {
		if p != nil {
			panic("database.InsertSubmission: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteSubmission(ctx context.Context, ref SubmissionRef) error {
	_, err := d.UpdateSubmission(ctx, ref, func(p *Submission) *Submission {
		return nil
	})
	return err
}

type JudgerRef string

func NewJudgerRef() JudgerRef {
	return JudgerRef(newId())
}

func CreateJudgerRef(ref JudgerRef) *JudgerRef {
	x := ref
	return &x
}

func (d *Database) getJudger(ctx context.Context, ref JudgerRef) (*Judger, error) {
	v := new(Judger)

	err := d.QueryRowContext(ctx, "SELECT id, token FROM judger WHERE id=?", ref).Scan(&v.Id, &v.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (d *Database) updateJudger(ctx context.Context, ref JudgerRef, v *Judger) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE judger SET token=? WHERE id=?", v.Token, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Judger")
	}
}

func (d *Database) insertJudger(ctx context.Context, v *Judger) {
	_, err := d.ExecContext(ctx, "INSERT INTO judger (id, token) VALUES (?, ?)", v.Id, v.Token)
	if err != nil {
		log.WithError(err).Error("Failed to insert Judger")
	}
}

func (d *Database) deleteJudger(ctx context.Context, ref JudgerRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM judger WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete Judger")
	}
}

func (d *Database) GetJudger(ctx context.Context, ref JudgerRef) (*Judger, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*Judger), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*Judger), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getJudger(ctx, ref)
	d.m.Lock()
	if err != nil {
		d.done(ref)
		d.m.Unlock()
		return nil, err
	}
	entry.curData = entry.prevData
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	return entry.curData.(*Judger), nil
}

func (d *Database) UpdateJudger(ctx context.Context, ref JudgerRef, updater func(*Judger) *Judger) (*Judger, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getJudger(ctx, ref)
		if err != nil {
			d.m.Lock()
			d.done(ref)
			d.m.Unlock()
			return nil, err
		}
		entry.prevData = entry.curData
		d.m.Lock()
		d.cache[ref] = entry
	}
	entry.curData = updater(entry.curData.(*Judger))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushJudger(d.ctx, ref)
	})
	return entry.curData.(*Judger), nil
}

func (d *Database) FlushJudger(ctx context.Context, ref JudgerRef) {
	d.m.Lock()
	if err := d.optWait(d.ctx, ref); err != nil {
		return
	}
	entry, found := d.cache[ref]
	if !found || entry.prevData == entry.curData {
		d.done(ref)
		d.m.Unlock()
		return
	}
	prevData := entry.prevData.(*Judger)
	curData := entry.curData.(*Judger)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertJudger(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteJudger(d.ctx, ref)
		} else {
			d.updateJudger(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertJudger(ctx context.Context, v *Judger) error {
	if v.Id == nil {
		v.Id = CreateJudgerRef(NewJudgerRef())
	}
	_, err := d.UpdateJudger(ctx, v.GetId(), func(p *Judger) *Judger {
		if p != nil {
			panic("database.InsertJudger: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteJudger(ctx context.Context, ref JudgerRef) error {
	_, err := d.UpdateJudger(ctx, ref, func(p *Judger) *Judger {
		return nil
	})
	return err
}

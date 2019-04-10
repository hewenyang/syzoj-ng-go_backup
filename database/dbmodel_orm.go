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
	var var2 *time.Time
	err := d.QueryRowContext(ctx, "SELECT id, user, create_time, title FROM problem WHERE id=?", ref).Scan(&v.Id, &v.User, &var2, &v.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if var2 != nil {
		v.CreateTime, _ = ptypes.TimestampProto(*var2)
	} else {
		v.CreateTime = nil
	}
	return v, nil
}

func (d *Database) updateProblem(ctx context.Context, ref ProblemRef, v *Problem) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE problem SET user=?, create_time=?, title=? WHERE id=?", v.User, convertTimestamp(v.CreateTime), v.Title, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Problem")
	}
}

func (d *Database) insertProblem(ctx context.Context, v *Problem) {
	_, err := d.ExecContext(ctx, "INSERT INTO problem (id, user, create_time, title) VALUES (?, ?, ?, ?)", v.Id, v.User, convertTimestamp(v.CreateTime), v.Title)
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

type ProblemEntryRef string

func NewProblemEntryRef() ProblemEntryRef {
	return ProblemEntryRef(newId())
}

func CreateProblemEntryRef(ref ProblemEntryRef) *ProblemEntryRef {
	x := ref
	return &x
}

func (d *Database) getProblemEntry(ctx context.Context, ref ProblemEntryRef) (*ProblemEntry, error) {
	v := new(ProblemEntry)

	err := d.QueryRowContext(ctx, "SELECT id, title, problem FROM problem_entry WHERE id=?", ref).Scan(&v.Id, &v.Title, &v.Problem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return v, nil
}

func (d *Database) updateProblemEntry(ctx context.Context, ref ProblemEntryRef, v *ProblemEntry) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE problem_entry SET title=?, problem=? WHERE id=?", v.Title, v.Problem, v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update ProblemEntry")
	}
}

func (d *Database) insertProblemEntry(ctx context.Context, v *ProblemEntry) {
	_, err := d.ExecContext(ctx, "INSERT INTO problem_entry (id, title, problem) VALUES (?, ?, ?)", v.Id, v.Title, v.Problem)
	if err != nil {
		log.WithError(err).Error("Failed to insert ProblemEntry")
	}
}

func (d *Database) deleteProblemEntry(ctx context.Context, ref ProblemEntryRef) {
	_, err := d.ExecContext(ctx, "DELETE FROM problem_entry WHERE id=?", ref)
	if err != nil {
		log.WithError(err).Error("Failed to delete ProblemEntry")
	}
}

func (d *Database) GetProblemEntry(ctx context.Context, ref ProblemEntryRef) (*ProblemEntry, error) {
	d.m.Lock()
	entry, found := d.cache[ref]
	if found {
		d.m.Unlock()
		return entry.curData.(*ProblemEntry), nil
	}
	// slow path
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found = d.cache[ref]
	if found {
		d.done(ref)
		d.m.Unlock()
		return entry.curData.(*ProblemEntry), nil
	}
	d.m.Unlock()
	var err error
	entry.prevData, err = d.getProblemEntry(ctx, ref)
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
	return entry.curData.(*ProblemEntry), nil
}

func (d *Database) UpdateProblemEntry(ctx context.Context, ref ProblemEntryRef, updater func(*ProblemEntry) *ProblemEntry) (*ProblemEntry, error) {
	d.m.Lock()
	if err := d.optWait(ctx, ref); err != nil {
		return nil, err
	}
	entry, found := d.cache[ref]
	if !found {
		d.m.Unlock()
		var err error
		entry.curData, err = d.getProblemEntry(ctx, ref)
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
	entry.curData = updater(entry.curData.(*ProblemEntry))
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
	d.wg.Add(1)
	time.AfterFunc(time.Millisecond*5, func() {
		defer d.wg.Done()
		d.FlushProblemEntry(d.ctx, ref)
	})
	return entry.curData.(*ProblemEntry), nil
}

func (d *Database) FlushProblemEntry(ctx context.Context, ref ProblemEntryRef) {
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
	prevData := entry.prevData.(*ProblemEntry)
	curData := entry.curData.(*ProblemEntry)
	d.m.Unlock()
	if prevData == nil {
		if curData != nil {
			d.insertProblemEntry(d.ctx, curData)
		}
	} else {
		if curData == nil {
			d.deleteProblemEntry(d.ctx, ref)
		} else {
			d.updateProblemEntry(d.ctx, ref, curData)
		}
	}
	entry.prevData = entry.curData
	d.m.Lock()
	d.cache[ref] = entry
	d.done(ref)
	d.m.Unlock()
}

func (d *Database) InsertProblemEntry(ctx context.Context, v *ProblemEntry) error {
	if v.Id == nil {
		v.Id = CreateProblemEntryRef(NewProblemEntryRef())
	}
	_, err := d.UpdateProblemEntry(ctx, v.GetId(), func(p *ProblemEntry) *ProblemEntry {
		if p != nil {
			panic("database.InsertProblemEntry: Duplicate primary key")
		}
		return v
	})
	return err
}

func (d *Database) DeleteProblemEntry(ctx context.Context, ref ProblemEntryRef) error {
	_, err := d.UpdateProblemEntry(ctx, ref, func(p *ProblemEntry) *ProblemEntry {
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
	var var3 []byte
	err := d.QueryRowContext(ctx, "SELECT id, problem_judger, user, data FROM submission WHERE id=?", ref).Scan(&v.Id, &v.ProblemJudger, &v.User, &var3)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if var3 != nil {
		v.Data = &any.Any{}
		if err := proto.Unmarshal(var3, v.Data); err != nil {
			panic(err)
		}
	} else {
		v.Data = nil
	}
	return v, nil
}

func (d *Database) updateSubmission(ctx context.Context, ref SubmissionRef, v *Submission) {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := d.ExecContext(ctx, "UPDATE submission SET problem_judger=?, user=?, data=? WHERE id=?", v.ProblemJudger, v.User, convertAny(v.Data), v.Id)
	if err != nil {
		log.WithError(err).Error("Failed to update Submission")
	}
}

func (d *Database) insertSubmission(ctx context.Context, v *Submission) {
	_, err := d.ExecContext(ctx, "INSERT INTO submission (id, problem_judger, user, data) VALUES (?, ?, ?, ?)", v.Id, v.ProblemJudger, v.User, convertAny(v.Data))
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

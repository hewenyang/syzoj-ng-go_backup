package database

import (
	"context"
	"database/sql"
)

type UserRef string

func NewUserRef() UserRef {
	return UserRef(newId())
}

func CreateUserRef(ref UserRef) *UserRef {
	x := ref
	return &x
}

func (t *DatabaseTxn) GetUser(ctx context.Context, ref UserRef) (*User, error) {
	v := new(User)
	err := t.tx.QueryRowContext(ctx, "SELECT id, user_name, auth FROM user WHERE id=?", ref).Scan(&v.Id, &v.UserName, &v.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetUserForUpdate(ctx context.Context, ref UserRef) (*User, error) {
	v := new(User)
	err := t.tx.QueryRowContext(ctx, "SELECT id, user_name, auth FROM user WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.UserName, &v.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateUser(ctx context.Context, ref UserRef, v *User) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE user SET user_name=?, auth=? WHERE id=?", v.UserName, v.Auth, v.Id)
	return err
}

func (t *DatabaseTxn) InsertUser(ctx context.Context, v *User) error {
	if v.Id == nil {
		ref := NewUserRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO user (id, user_name, auth) VALUES (?, ?, ?)", v.Id, v.UserName, v.Auth)
	return err
}

func (t *DatabaseTxn) DeleteUser(ctx context.Context, ref UserRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM user WHERE id=?", ref)
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

func (t *DatabaseTxn) GetDevice(ctx context.Context, ref DeviceRef) (*Device, error) {
	v := new(Device)
	err := t.tx.QueryRowContext(ctx, "SELECT id, user, info FROM device WHERE id=?", ref).Scan(&v.Id, &v.User, &v.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetDeviceForUpdate(ctx context.Context, ref DeviceRef) (*Device, error) {
	v := new(Device)
	err := t.tx.QueryRowContext(ctx, "SELECT id, user, info FROM device WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.User, &v.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateDevice(ctx context.Context, ref DeviceRef, v *Device) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE device SET user=?, info=? WHERE id=?", v.User, v.Info, v.Id)
	return err
}

func (t *DatabaseTxn) InsertDevice(ctx context.Context, v *Device) error {
	if v.Id == nil {
		ref := NewDeviceRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO device (id, user, info) VALUES (?, ?, ?)", v.Id, v.User, v.Info)
	return err
}

func (t *DatabaseTxn) DeleteDevice(ctx context.Context, ref DeviceRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM device WHERE id=?", ref)
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

func (t *DatabaseTxn) GetProblem(ctx context.Context, ref ProblemRef) (*Problem, error) {
	v := new(Problem)
	err := t.tx.QueryRowContext(ctx, "SELECT id, title, user FROM problem WHERE id=?", ref).Scan(&v.Id, &v.Title, &v.User)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetProblemForUpdate(ctx context.Context, ref ProblemRef) (*Problem, error) {
	v := new(Problem)
	err := t.tx.QueryRowContext(ctx, "SELECT id, title, user FROM problem WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.Title, &v.User)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblem(ctx context.Context, ref ProblemRef, v *Problem) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem SET title=?, user=? WHERE id=?", v.Title, v.User, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblem(ctx context.Context, v *Problem) error {
	if v.Id == nil {
		ref := NewProblemRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem (id, title, user) VALUES (?, ?, ?)", v.Id, v.Title, v.User)
	return err
}

func (t *DatabaseTxn) DeleteProblem(ctx context.Context, ref ProblemRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem WHERE id=?", ref)
	return err
}

type ProblemSourceRef string

func NewProblemSourceRef() ProblemSourceRef {
	return ProblemSourceRef(newId())
}

func CreateProblemSourceRef(ref ProblemSourceRef) *ProblemSourceRef {
	x := ref
	return &x
}

func (t *DatabaseTxn) GetProblemSource(ctx context.Context, ref ProblemSourceRef) (*ProblemSource, error) {
	v := new(ProblemSource)
	err := t.tx.QueryRowContext(ctx, "SELECT id, source FROM problem_source WHERE id=?", ref).Scan(&v.Id, &v.Source)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetProblemSourceForUpdate(ctx context.Context, ref ProblemSourceRef) (*ProblemSource, error) {
	v := new(ProblemSource)
	err := t.tx.QueryRowContext(ctx, "SELECT id, source FROM problem_source WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.Source)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemSource(ctx context.Context, ref ProblemSourceRef, v *ProblemSource) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_source SET source=? WHERE id=?", v.Source, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemSource(ctx context.Context, v *ProblemSource) error {
	if v.Id == nil {
		ref := NewProblemSourceRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_source (id, source) VALUES (?, ?)", v.Id, v.Source)
	return err
}

func (t *DatabaseTxn) DeleteProblemSource(ctx context.Context, ref ProblemSourceRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_source WHERE id=?", ref)
	return err
}

type ProblemJudgerRef string

func NewProblemJudgerRef() ProblemJudgerRef {
	return ProblemJudgerRef(newId())
}

func CreateProblemJudgerRef(ref ProblemJudgerRef) *ProblemJudgerRef {
	x := ref
	return &x
}

func (t *DatabaseTxn) GetProblemJudger(ctx context.Context, ref ProblemJudgerRef) (*ProblemJudger, error) {
	v := new(ProblemJudger)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem, user, type, data FROM problem_judger WHERE id=?", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Type, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetProblemJudgerForUpdate(ctx context.Context, ref ProblemJudgerRef) (*ProblemJudger, error) {
	v := new(ProblemJudger)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem, user, type, data FROM problem_judger WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Type, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemJudger(ctx context.Context, ref ProblemJudgerRef, v *ProblemJudger) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_judger SET problem=?, user=?, type=?, data=? WHERE id=?", v.Problem, v.User, v.Type, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemJudger(ctx context.Context, v *ProblemJudger) error {
	if v.Id == nil {
		ref := NewProblemJudgerRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_judger (id, problem, user, type, data) VALUES (?, ?, ?, ?, ?)", v.Id, v.Problem, v.User, v.Type, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteProblemJudger(ctx context.Context, ref ProblemJudgerRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_judger WHERE id=?", ref)
	return err
}

type ProblemStatementRef string

func NewProblemStatementRef() ProblemStatementRef {
	return ProblemStatementRef(newId())
}

func CreateProblemStatementRef(ref ProblemStatementRef) *ProblemStatementRef {
	x := ref
	return &x
}

func (t *DatabaseTxn) GetProblemStatement(ctx context.Context, ref ProblemStatementRef) (*ProblemStatement, error) {
	v := new(ProblemStatement)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem, user, data FROM problem_statement WHERE id=?", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetProblemStatementForUpdate(ctx context.Context, ref ProblemStatementRef) (*ProblemStatement, error) {
	v := new(ProblemStatement)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem, user, data FROM problem_statement WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemStatement(ctx context.Context, ref ProblemStatementRef, v *ProblemStatement) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_statement SET problem=?, user=?, data=? WHERE id=?", v.Problem, v.User, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemStatement(ctx context.Context, v *ProblemStatement) error {
	if v.Id == nil {
		ref := NewProblemStatementRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_statement (id, problem, user, data) VALUES (?, ?, ?, ?)", v.Id, v.Problem, v.User, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteProblemStatement(ctx context.Context, ref ProblemStatementRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_statement WHERE id=?", ref)
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

func (t *DatabaseTxn) GetSubmission(ctx context.Context, ref SubmissionRef) (*Submission, error) {
	v := new(Submission)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem_judger, user, data FROM submission WHERE id=?", ref).Scan(&v.Id, &v.ProblemJudger, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) GetSubmissionForUpdate(ctx context.Context, ref SubmissionRef) (*Submission, error) {
	v := new(Submission)
	err := t.tx.QueryRowContext(ctx, "SELECT id, problem_judger, user, data FROM submission WHERE id=? FOR UPDATE", ref).Scan(&v.Id, &v.ProblemJudger, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateSubmission(ctx context.Context, ref SubmissionRef, v *Submission) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE submission SET problem_judger=?, user=?, data=? WHERE id=?", v.ProblemJudger, v.User, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertSubmission(ctx context.Context, v *Submission) error {
	if v.Id == nil {
		ref := NewSubmissionRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO submission (id, problem_judger, user, data) VALUES (?, ?, ?, ?)", v.Id, v.ProblemJudger, v.User, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteSubmission(ctx context.Context, ref SubmissionRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM submission WHERE id=?", ref)
	return err
}

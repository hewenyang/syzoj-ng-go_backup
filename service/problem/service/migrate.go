package service

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	"github.com/syzoj/syzoj-ng-go/service"
)

const sqlFresh = `
CREATE TABLE problems (id VARCHAR(255) PRIMARY KEY, title VARCHAR(255), statement BLOB);
`

func (s *serv) Migrate(ctx context.Context, c *service.ServiceContext, prevVersion string) error {
	s.log = c.GetLogger()
	// Alter config to support multi statements
	mscfg, err := mysql.ParseDSN(s.config.MySQL)
	if err != nil {
		return err
	}
	mscfg.MultiStatements = true
	db, err := sql.Open("mysql", mscfg.FormatDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	if prevVersion == "" {
		_, err := db.ExecContext(ctx, sqlFresh)
		return err
	}
	return nil
}

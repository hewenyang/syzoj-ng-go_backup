package service

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) Migrate(prevVersion string) error {
	ctx := context.Background()
	switch prevVersion {
	case "":
		return util.MigrateMySQL(ctx, s.config.MySQL, `CREATE TABLE user_profile (id VARCHAR(64) PRIMARY KEY, profile BLOB)`)
	}
	return nil
}

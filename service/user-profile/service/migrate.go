package service

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/service"
	"github.com/syzoj/syzoj-ng-go/util"
)

func (s *serv) Migrate(ctx context.Context, c *service.ServiceContext, prevVersion string) error {
	switch prevVersion {
	case "":
		return util.MigrateMySQL(ctx, s.config.MySQL, `CREATE TABLE user_profile (id VARCHAR(64) PRIMARY KEY, profile BLOB)`)
	}
	return nil
}

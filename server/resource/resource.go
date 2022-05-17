package resource

import (
	"database/sql"

	"github.com/go-redis/redis"
)

var MysqlClient *sql.DB

var RedisClient *redis.Client

var UserProfileStatement *sql.Stmt

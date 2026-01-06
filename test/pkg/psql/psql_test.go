package psql

import (
	"database/sql"
	"testing"

	"github.com/parxyws/nego-gin/pkg/database/psql"
	"github.com/parxyws/nego-gin/test"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	_, err := test.InitEnv()
	assert.Nil(t, err)
}

func TestDatabaseConn(t *testing.T) {
	cfg, err := test.InitEnv()
	assert.Nil(t, err)

	db, err := psql.NewDB(cfg)
	assert.Nil(t, err)
	conn, err := db.DB()
	assert.Nil(t, err)
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
}

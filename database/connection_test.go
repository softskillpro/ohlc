package database_test

import (
	"github.com/stretchr/testify/assert"
	"ohcl/config"
	"ohcl/database"
	"testing"
)

func TestConnection(t *testing.T) {
	conn, err := config.Get("DATABASE")
	assert.NoError(t, err)
	assert.NotEmpty(t, conn)

	err = database.Db().Ping()

	assert.NoError(t, err) // Check that there are no errors
}

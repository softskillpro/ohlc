package main

import (
	"ohcl/database"
	"ohcl/database/migrate"
)

// Initializer will initialize database and do auto migration on database.
func Initializer() error {

	if err := database.Db().Error(); err != nil {
		return err
	}

	if err := database.Db().Migration(migrate.Stmt).Error(); err != nil {
		return err
	}
	return nil
}

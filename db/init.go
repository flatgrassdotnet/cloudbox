/*
	cloudbox - the toybox server emulator
	Copyright (C) 2024-2025  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/go-sql-driver/mysql"
)

var (
	handle   *sql.DB
	s3client *s3.Client
)

func Init(username string, password string, protocol string, address string, database string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", username, password, protocol, address, database))
	if err != nil {
		return err
	}

	handle = db

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	s3client = s3.NewFromConfig(cfg)

	return nil
}

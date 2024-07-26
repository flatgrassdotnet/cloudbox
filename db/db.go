/*
	reboxed - the toybox server emulator
	Copyright (C) 2024  patapancakes <patapancakes@pagefault.games>

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
	"database/sql"
	"fmt"
	"reboxed/utils"

	_ "github.com/go-sql-driver/mysql"
)

var handle *sql.DB

func Init(username string, password string, address string, database string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, address, database))
	if err != nil {
		return err
	}

	handle = db

	return nil
}

func InsertLogin(version int, steamid int, vac string, ticket []byte) error {
	_, err := handle.Exec("INSERT INTO logins (version, steamid, vac, ticket) VALUES (?, ?, ?, ?)", version, steamid, vac, ticket)
	if err != nil {
		return err
	}

	return nil
}

func FetchSteamIDFromTicket(ticket []byte) (int64, error) {
	var steamid int64
	err := handle.QueryRow("SELECT steamid FROM logins WHERE ticket = ?", ticket).Scan(&steamid)
	if err != nil {
		return 0, err
	}

	return steamid, nil
}

func InsertMapLoad(version int, steamid int, duration float64, mapName string, platform string) error {
	_, err := handle.Exec("INSERT INTO maploads (version, steamid, duration, map, platform) VALUES (?, ?, ?, ?, ?)", version, steamid, duration, mapName, platform)
	if err != nil {
		return err
	}

	return nil
}

func InsertError(version int, steamid int, error string, content string, realm string, platform string) error {
	_, err := handle.Exec("INSERT INTO errors (version, steamid, error, content, realm, platform) VALUES (?, ?, ?, ?, ?, ?)", version, steamid, error, content, realm, platform)
	if err != nil {
		return err
	}

	return nil
}

func InsertPackage(packageType string, name string, dataname string, author int64, description string, data []byte) (int, error) {
	r, err := handle.Exec("INSERT INTO packages (type, name, dataname, author, description, data) VALUES (?, ?, ?, ?, ?, ?)", packageType, name, dataname, author, description, data)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func FetchPackage(scriptid int, rev int) (utils.Package, error) {
	var pkg utils.Package
	err := handle.QueryRow("SELECT id, rev, type, name, dataname, data FROM packages WHERE id = ? AND rev = ?", scriptid, rev).Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Data)
	if err != nil {
		return pkg, err
	}

	rows, err := handle.Query("SELECT f.id, f.rev, f.path, f.size, f.psize FROM files f JOIN content c ON f.id = c.fileid AND f.rev = c.filerev WHERE c.id = ? AND c.rev = ?", scriptid, rev)
	if err != nil {
		return pkg, err
	}

	for rows.Next() {
		var content utils.Content
		err := rows.Scan(&content.ID, &content.Revision, &content.Path, &content.Size, &content.PSize)
		if err != nil {
			return pkg, err
		}

		pkg.Content = append(pkg.Content, content)
	}

	rows, err = handle.Query("SELECT p.id, p.rev, p.type FROM packages p JOIN includes i ON p.id = i.includeid AND p.rev = i.includerev WHERE i.id = ? AND i.rev = ?", scriptid, rev)
	if err != nil {
		return pkg, err
	}

	for rows.Next() {
		var include utils.Include
		err := rows.Scan(&include.ID, &include.Revision, &include.Type)
		if err != nil {
			return pkg, err
		}

		pkg.Includes = append(pkg.Includes, include)
	}

	return pkg, nil
}

func FetchPackageListByType(category string) ([]utils.Package, error) {
	var list []utils.Package

	rows, err := handle.Query("SELECT p.id, p.rev, p.type, p.name, p.dataname FROM packages p WHERE p.type = ? AND p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id)", category)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var pkg utils.Package
		err := rows.Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname)
		if err != nil {
			return list, err
		}

		list = append(list, pkg)
	}

	return list, nil
}

func FetchPackageListByTypePaged(category string, offset int, count int) ([]utils.Package, error) {
	var list []utils.Package

	rows, err := handle.Query("SELECT p.id, p.rev, p.type, p.name, p.dataname FROM packages p WHERE p.type = ? AND p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id) LIMIT ?, ?", category, offset, count)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var pkg utils.Package
		err := rows.Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname)
		if err != nil {
			return list, err
		}

		list = append(list, pkg)
	}

	return list, nil
}

func InsertUpload(steamid int, upload utils.Upload) (int, error) {
	r, err := handle.Exec("INSERT INTO uploads (steamid, type, meta, inc, data) VALUES (?, ?, ?, ?, ?)", steamid, upload.Type, upload.Metadata, upload.Include, upload.Data)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func FetchUpload(id int) (utils.Upload, error) {
	var upload utils.Upload
	err := handle.QueryRow("SELECT type, meta, inc, data FROM uploads WHERE id = ?", id).Scan(&upload.Type, &upload.Metadata, &upload.Include, &upload.Data)
	if err != nil {
		return upload, err
	}

	return upload, nil
}

func DeleteUpload(id int) error {
	_, err := handle.Exec("DELETE FROM uploads WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

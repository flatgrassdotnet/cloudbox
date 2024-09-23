/*
	cloudbox - the toybox server emulator
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

import "github.com/flatgrassdotnet/cloudbox/common"

func InsertUpload(steamid int, upload common.Upload) (int, error) {
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

func FetchUpload(id int) (common.Upload, error) {
	var upload common.Upload
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

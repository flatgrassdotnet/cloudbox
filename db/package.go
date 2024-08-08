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

import "reboxed/common"

func InsertPackage(packageType string, name string, dataname string, author uint64, description string, data []byte) (int, error) {
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

func InsertPackageInclude(id int, rev int, iid int, irev int) error {
	_, err := handle.Exec("INSERT INTO includes (id, rev, includeid, includerev) VALUES (?, ?, ?, ?)", id, rev, iid, irev)
	if err != nil {
		return err
	}

	return nil
}

func FetchPackageLatestRevision(id int) (int, error) {
	var rev int
	err := handle.QueryRow("SELECT MAX(rev) FROM packages WHERE id = ?", id).Scan(&rev)
	if err != nil {
		return 0, err
	}

	return rev, nil
}

func FetchPackage(id int, rev int) (common.Package, error) {
	var pkg common.Package
	err := handle.QueryRow("SELECT id, rev, type, name, dataname, author, description, data FROM packages WHERE id = ? AND rev = ?", id, rev).Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.Description, &pkg.Data)
	if err != nil {
		return pkg, err
	}

	rows, err := handle.Query("SELECT f.id, f.rev, f.path, f.size, f.psize FROM files f JOIN content c ON f.id = c.fileid AND f.rev = c.filerev WHERE c.id = ? AND c.rev = ?", id, rev)
	if err != nil {
		return pkg, err
	}

	for rows.Next() {
		var content common.Content
		err := rows.Scan(&content.ID, &content.Revision, &content.Path, &content.Size, &content.PSize)
		if err != nil {
			return pkg, err
		}

		pkg.Content = append(pkg.Content, content)
	}

	rows, err = handle.Query("SELECT p.id, p.rev, p.type FROM packages p JOIN includes i ON p.id = i.includeid AND p.rev = i.includerev WHERE i.id = ? AND i.rev = ?", id, rev)
	if err != nil {
		return pkg, err
	}

	for rows.Next() {
		var include common.Include
		err := rows.Scan(&include.ID, &include.Revision, &include.Type)
		if err != nil {
			return pkg, err
		}

		pkg.Includes = append(pkg.Includes, include)
	}

	return pkg, nil
}

func FetchPackageList(category string, author uint64, search string, offset int, count int) ([]common.Package, error) {
	var args []any
	q := "SELECT p.id, p.rev, p.type, p.name, p.dataname, p.author, p.description FROM packages p WHERE p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id)"

	if category != "" {
		q += " AND p.type = ?"
		args = append(args, category)
	}

	if author != 0 {
		q += " AND p.author = ?"
		args = append(args, author)
	}

	if search != "" {
		q += " AND p.name LIKE CONCAT('%', ?, '%')"
		args = append(args, search)
	}

	if count != 0 {
		q += " LIMIT ?, ?"
		args = append(args, offset)
		args = append(args, count)
	}

	var list []common.Package

	rows, err := handle.Query(q, args...)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var pkg common.Package
		err := rows.Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.Description)
		if err != nil {
			return list, err
		}

		list = append(list, pkg)
	}

	return list, nil
}

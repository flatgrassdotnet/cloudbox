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

func FetchPackage(scriptid int, rev int) (common.Package, error) {
	var pkg common.Package
	err := handle.QueryRow("SELECT id, rev, type, name, dataname, author, description, data FROM packages WHERE id = ? AND rev = ?", scriptid, rev).Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.Description, &pkg.Data)
	if err != nil {
		return pkg, err
	}

	rows, err := handle.Query("SELECT f.id, f.rev, f.path, f.size, f.psize FROM files f JOIN content c ON f.id = c.fileid AND f.rev = c.filerev WHERE c.id = ? AND c.rev = ?", scriptid, rev)
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

	rows, err = handle.Query("SELECT p.id, p.rev, p.type FROM packages p JOIN includes i ON p.id = i.includeid AND p.rev = i.includerev WHERE i.id = ? AND i.rev = ?", scriptid, rev)
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

func FetchPackageList(category string) ([]common.Package, error) {
	var list []common.Package

	rows, err := handle.Query("SELECT p.id, p.rev, p.type, p.name, p.dataname, p.author, p.description FROM packages p WHERE p.type = ? AND p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id)", category)
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

func FetchPackageListPaged(category string, query string, offset int, count int) ([]common.Package, error) {
	var list []common.Package

	rows, err := handle.Query("SELECT p.id, p.rev, p.type, p.name, p.dataname, p.author, p.description FROM packages p WHERE p.type = ? AND p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id) AND p.name LIKE CONCAT('%', ?, '%') LIMIT ?, ?", category, query, offset, count)
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

func FetchAuthorPackageListPaged(author uint64, query string, offset int, count int) ([]common.Package, error) {
	var list []common.Package

	rows, err := handle.Query("SELECT p.id, p.rev, p.type, p.name, p.dataname, p.author, p.description FROM packages p WHERE p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id) AND p.author = ? AND p.name LIKE CONCAT('%', ?, '%') LIMIT ?, ?", author, query, offset, count)
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

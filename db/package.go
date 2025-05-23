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
	"fmt"

	"github.com/flatgrassdotnet/cloudbox/common"
)

func InsertPackage(pkg common.Package) (int, error) {
	r, err := handle.Exec("INSERT INTO packages (type, name, dataname, author, description, data) VALUES (?, ?, ?, ?, ?, ?)", pkg.Type, pkg.Name, pkg.Dataname, pkg.Author, pkg.Description, pkg.Data)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func InsertPackageInclude(id int, rev int, iid int, irev int) (int, error) {
	r, err := handle.Exec("INSERT INTO includes (id, rev, includeid, includerev) VALUES (?, ?, ?, ?)", id, rev, iid, irev)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil
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
	err := handle.QueryRow("SELECT id, rev, type, name, COALESCE(dataname, \"\"), COALESCE(author, \"\"), COALESCE(description, \"\"), data FROM packages WHERE id = ? AND rev = ?", id, rev).Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.Description, &pkg.Data)
	if err != nil {
		return pkg, err
	}

	rows, err := handle.Query("SELECT f.id, f.path, f.size, f.psize FROM files f JOIN content c ON f.id = c.fileid WHERE c.id = ?", id)
	if err != nil {
		return pkg, err
	}

	for rows.Next() {
		var content common.Content
		err := rows.Scan(&content.ID, &content.Path, &content.Size, &content.PSize)
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

func FetchPackageList(category string, dataname string, author string, search string, offset int, count int, sort string, safemode bool) ([]common.Package, error) {
	var args []any
	q := `SELECT 
	p.id, 
	p.rev, 
	p.type, 
	p.name, 
	COALESCE(p.dataname, ""), 
	COALESCE(p.author, ""), 
	COALESCE(pr.personaname, s.author, ""), 
	COALESCE(pr.avatarmedium, ""), 
	COALESCE(p.description, s.description, ""), 
	COALESCE(s.downloads, 0), 
	COALESCE(s.favorites, 0), 
	COALESCE(s.goods, 0), 
	COALESCE(s.bads, 0), 
	p.time 
	FROM packages p 
	LEFT JOIN profiles pr
	ON p.author = pr.steamid 
	LEFT JOIN scraped s
	ON p.id = s.id 
	AND s.rev = (SELECT MAX(s2.rev) FROM scraped s2 WHERE s2.id = s.id) 
	WHERE p.rev = (SELECT MAX(p2.rev) FROM packages p2 WHERE p2.id = p.id)`

	if category != "" {
		q += " AND p.type = ?"
		args = append(args, category)
	}

	if author != "" {
		q += " AND p.author = ?"
		args = append(args, author)
	}

	if search != "" {
		q += " AND p.name LIKE CONCAT('%', ?, '%')"
		args = append(args, search)
	}

	if dataname != "" {
		q += " AND dataname = ?"
		args = append(args, dataname)
	}

	if safemode {
		q += " AND incompatible = 0"
	}

	// dangerous!
	if sort != "" {
		q += fmt.Sprintf(" ORDER BY %s DESC", sort)
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
		err := rows.Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.AuthorName, &pkg.AuthorIcon, &pkg.Description, &pkg.Downloads, &pkg.Favorites, &pkg.Goods, &pkg.Bads, &pkg.Uploaded)
		if err != nil {
			return list, err
		}

		list = append(list, pkg)
	}

	return list, nil
}

func FetchPackageListAll(category string, dataname string, author string, search string, offset int, count int, sort string) ([]common.Package, error) {
	var args []any
	q := `WITH latest_packages AS (
	SELECT *
	FROM (
		SELECT *,
			ROW_NUMBER() OVER (PARTITION BY id ORDER BY rev DESC) AS rn
		FROM packages
	) AS ranked
	WHERE rn = 1
),
latest_scraped AS (
	SELECT *
	FROM (
		SELECT *,
			ROW_NUMBER() OVER (PARTITION BY id ORDER BY rev DESC) AS rn
		FROM scraped
	) AS ranked
	WHERE rn = 1
),
combined AS (
	SELECT 
		p.id,
		p.rev,
		p.type,
		p.name,
		COALESCE(p.dataname, '') AS dataname,
		COALESCE(p.author, '') AS author,
		COALESCE(pr.personaname, s.author, '') AS personaname,
		COALESCE(pr.avatarmedium, '') AS avatarmedium,
		COALESCE(p.description, s.description, '') AS description,
		COALESCE(s.downloads, 0) AS downloads,
		COALESCE(s.favorites, 0) AS favorites,
		COALESCE(s.goods, 0) AS goods,
		COALESCE(s.bads, 0) AS bads,
		p.time
	FROM latest_packages p
	LEFT JOIN profiles pr ON p.author = pr.steamid
	LEFT JOIN latest_scraped s ON p.id = s.id

	UNION ALL

	SELECT
		s.id,
		s.rev,
		s.type,
		s.name,
		'' AS dataname,
		'' AS author,
		s.author AS personaname,
		'' AS avatarmedium,
		s.description,
		s.downloads,
		s.favorites,
		s.goods,
		s.bads,
		STR_TO_DATE('2012-10-25', '%Y-%c-%d') AS time
	FROM latest_scraped s
	WHERE NOT EXISTS (SELECT 1 FROM latest_packages p WHERE p.id = s.id)
)

SELECT *
FROM combined
WHERE 1 = 1`

	if category != "" {
		q += " AND type = ?"
		args = append(args, category)
	}

	if author != "" {
		q += " AND author = ?"
		args = append(args, author)
	}

	if search != "" {
		q += " AND name LIKE CONCAT('%', ?, '%')"
		args = append(args, search)
	}

	if dataname != "" {
		q += " AND dataname = ?"
		args = append(args, dataname)
	}

	// dangerous!
	if sort != "" {
		q += fmt.Sprintf(" ORDER BY %s DESC", sort)
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
		err := rows.Scan(&pkg.ID, &pkg.Revision, &pkg.Type, &pkg.Name, &pkg.Dataname, &pkg.Author, &pkg.AuthorName, &pkg.AuthorIcon, &pkg.Description, &pkg.Downloads, &pkg.Favorites, &pkg.Goods, &pkg.Bads, &pkg.Uploaded)
		if err != nil {
			return list, err
		}

		list = append(list, pkg)
	}

	return list, nil
}

func FetchFileInfoFromPath(path string) (int, error) {
	var id int
	err := handle.QueryRow("SELECT id FROM files WHERE path = ?", path).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

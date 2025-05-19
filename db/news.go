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

import "github.com/flatgrassdotnet/cloudbox/common"

func FetchNewsEntries() ([]common.NewsEntry, error) {
	var entries []common.NewsEntry
	rows, err := handle.Query("SELECT id, title, body, author, time FROM news")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var entry common.NewsEntry
		err := rows.Scan(&entry.ID, &entry.Title, &entry.Body, &entry.Author, &entry.Time)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

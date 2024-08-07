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

func InsertLogin(version int, steamid int, vac string, ticket []byte) error {
	_, err := handle.Exec("INSERT INTO logins (version, steamid, vac, ticket) VALUES (?, ?, ?, ?)", version, steamid, vac, ticket)
	if err != nil {
		return err
	}

	return nil
}

func FetchSteamIDFromTicket(ticket []byte) (uint64, error) {
	var steamid uint64
	err := handle.QueryRow("SELECT steamid FROM logins WHERE ticket = ?", ticket).Scan(&steamid)
	if err != nil {
		return 0, err
	}

	return steamid, nil
}

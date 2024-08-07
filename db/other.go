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
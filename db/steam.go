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
	"reboxed/common"
)

func InsertPlayerSummary(s common.PlayerSummaryInfo) error {
	_, err := handle.Exec("REPLACE INTO profiles (steamid, personaname, avatar, avatarmedium, avatarfull) VALUES (?, ?, ?, ?, ?)", s.SteamID, s.PersonaName, s.Avatar, s.AvatarMedium, s.AvatarFull)
	if err != nil {
		return err
	}

	return nil
}

func FetchPlayerSummary(steamid string) (common.PlayerSummaryInfo, error) {
	var s common.PlayerSummaryInfo
	err := handle.QueryRow("SELECT personaname, avatar, avatarmedium, avatarfull FROM profiles WHERE time > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 WEEK) AND steamid = ?", steamid).Scan(&s.PersonaName, &s.Avatar, &s.AvatarMedium, &s.AvatarFull)
	if err != nil {
		return s, err
	}

	return s, nil
}

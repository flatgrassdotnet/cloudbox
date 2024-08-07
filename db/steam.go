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
	"time"
)

func InsertPlayerSummary(s common.PlayerSummaryInfo) error {
	_, err := handle.Exec("REPLACE INTO profiles (steamid, communityvisibilitystate, profilestate, personaname, lastlogoff, profileurl, avatar, avatarmedium, avatarfull) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", s.SteamID, s.CommunityVisibilityState, s.ProfileState, s.PersonaName, time.Unix(int64(s.LastLogoff), 0), s.ProfileURL, s.Avatar, s.AvatarMedium, s.AvatarFull)
	if err != nil {
		return err
	}

	return nil
}

func FetchPlayerSummary(steamid uint64) (common.PlayerSummaryInfo, error) {
	// workaround
	var lastlogoff time.Time

	var s common.PlayerSummaryInfo
	err := handle.QueryRow("SELECT communityvisibilitystate, profilestate, personaname, lastlogoff, profileurl, avatar, avatarmedium, avatarfull FROM profiles WHERE time > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 WEEK) AND steamid = ?", steamid).Scan(&s.CommunityVisibilityState, &s.ProfileState, &s.PersonaName, &lastlogoff, &s.ProfileURL, &s.Avatar, &s.AvatarMedium, &s.AvatarFull)
	if err != nil {
		return s, err
	}

	// workaround
	s.LastLogoff = int(lastlogoff.Unix())

	return s, nil
}

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

package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
)

var SteamAPIKey string

type AuthenticateUserTicketResponse struct {
	Response struct {
		Params common.UserTicketInfo `json:"params"`
		Error  struct {
			ErrorCode int    `json:"errorcode"`
			ErrorDesc string `json:"errordesc"`
		} `json:"error"`
	} `json:"response"`
}

func AuthenticateUserTicket(ticket string) (common.UserTicketInfo, error) {
	v := make(url.Values)

	v.Set("key", SteamAPIKey)
	v.Set("appid", "4000") // garry's mod
	v.Set("ticket", ticket)

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v0001/?%s", v.Encode()))
	if err != nil {
		return common.UserTicketInfo{}, err
	}

	defer r.Body.Close()

	var rd AuthenticateUserTicketResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return common.UserTicketInfo{}, err
	}

	// no steamid, something is wrong
	if rd.Response.Params.SteamID == "" {
		return common.UserTicketInfo{}, fmt.Errorf(rd.Response.Error.ErrorDesc)
	}

	return rd.Response.Params, nil
}

type GetPlayerSummariesResponse struct {
	Response struct {
		Players []common.PlayerSummaryInfo `json:"players"`
	} `json:"response"`
}

func GetPlayerSummary(steamid string) (common.PlayerSummaryInfo, error) {
	// fetch from cache
	s, err := db.FetchPlayerSummary(steamid)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return s, fmt.Errorf("failed to fetch player summary: %s", err)
		}
	} else {
		return s, nil
	}

	// otherwise get new data
	v := make(url.Values)

	v.Set("key", SteamAPIKey)
	v.Set("steamids", steamid)

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?%s", v.Encode()))
	if err != nil {
		return s, err
	}

	defer r.Body.Close()

	var rd GetPlayerSummariesResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return s, err
	}

	if len(rd.Response.Players) == 0 {
		return s, fmt.Errorf("no players returned")
	}

	// insert into cache
	err = db.InsertPlayerSummary(rd.Response.Players[0])
	if err != nil {
		return s, fmt.Errorf("failed to insert player summary: %s", err)
	}

	return rd.Response.Players[0], nil
}

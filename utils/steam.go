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

package utils

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

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
		return common.UserTicketInfo{}, errors.New(rd.Response.Error.ErrorDesc)
	}

	return rd.Response.Params, nil
}

type GetPlayerSummariesResponse struct {
	Response struct {
		Players []common.PlayerSummaryInfo `json:"players"`
	} `json:"response"`
}

func GetPlayerSummaries(steamids ...string) ([]common.PlayerSummaryInfo, error) {
	var summaries []common.PlayerSummaryInfo

	// fetch from cache
	for i, steamid := range steamids {
		summary, err := db.FetchPlayerSummary(steamid)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("failed to fetch player summary: %s", err)
			}

			continue
		}

		// remove from todo
		steamids = slices.Delete(steamids, i, i)

		summaries = append(summaries, summary)
	}

	// return now if there's none left
	if len(steamids) == 0 {
		return summaries, nil
	}

	// comma separated steamids
	buf := new(bytes.Buffer)

	cw := csv.NewWriter(buf)
	cw.Write(steamids)
	cw.Flush()

	v := make(url.Values)

	v.Set("key", SteamAPIKey)
	v.Set("steamids", buf.String())

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?%s", v.Encode()))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	var rd GetPlayerSummariesResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return nil, err
	}

	if len(rd.Response.Players) == 0 {
		return nil, fmt.Errorf("no players returned")
	}

	for _, summary := range rd.Response.Players {
		// insert into cache
		err = db.InsertPlayerSummary(summary)
		if err != nil {
			return nil, fmt.Errorf("failed to insert player summary: %s", err)
		}

		// add to summaries
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

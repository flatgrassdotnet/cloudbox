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

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type AuthenticateUserTicketResponse struct {
	Response struct {
		Params SteamUserInfo `json:"params"`
		Error  struct {
			ErrorCode int    `json:"errorcode"`
			ErrorDesc string `json:"errordesc"`
		} `json:"error"`
	} `json:"response"`
}

type SteamUserInfo struct {
	Result          string `json:"result"`
	SteamID         string `json:"steamid"`
	OwnerSteamID    string `json:"ownersteamid"`
	VACBanned       bool   `json:"vacbanned"`
	PublisherBanned bool   `json:"publisherbanned"`
}

var WebAPIKey string

func GetSteamUserInfo(ticket string) (SteamUserInfo, error) {
	v := make(url.Values)

	v.Set("key", WebAPIKey)
	v.Set("appid", "4000")
	v.Set("ticket", ticket)

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v0001/?%s", v.Encode()))
	if err != nil {
		return SteamUserInfo{}, err
	}

	defer r.Body.Close()

	var rd AuthenticateUserTicketResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return SteamUserInfo{}, err
	}

	// no steamid, something is wrong
	if rd.Response.Params.SteamID == "" {
		return SteamUserInfo{}, fmt.Errorf(rd.Response.Error.ErrorDesc)
	}

	return rd.Response.Params, nil
}

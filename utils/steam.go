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
	"strconv"
)

var SteamAPIKey string

type AuthenticateUserTicketResponse struct {
	Response struct {
		Params UserTicketInfo `json:"params"`
		Error  struct {
			ErrorCode int    `json:"errorcode"`
			ErrorDesc string `json:"errordesc"`
		} `json:"error"`
	} `json:"response"`
}

type UserTicketInfo struct {
	Result          string `json:"result"`
	SteamID         string `json:"steamid"`
	OwnerSteamID    string `json:"ownersteamid"`
	VACBanned       bool   `json:"vacbanned"`
	PublisherBanned bool   `json:"publisherbanned"`
}

func AuthenticateUserTicket(ticket string) (UserTicketInfo, error) {
	v := make(url.Values)

	v.Set("key", SteamAPIKey)
	v.Set("appid", "4000") // garry's mod
	v.Set("ticket", ticket)

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v0001/?%s", v.Encode()))
	if err != nil {
		return UserTicketInfo{}, err
	}

	defer r.Body.Close()

	var rd AuthenticateUserTicketResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return UserTicketInfo{}, err
	}

	// no steamid, something is wrong
	if rd.Response.Params.SteamID == "" {
		return UserTicketInfo{}, fmt.Errorf(rd.Response.Error.ErrorDesc)
	}

	return rd.Response.Params, nil
}

type GetPlayerSummariesResponse struct {
	Response struct {
		Players []PlayerSummaryInfo `json:"players"`
	} `json:"response"`
}

type PlayerSummaryInfo struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	LastLogoff               int    `json:"lastlogoff"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
}

func GetPlayerSummary(steamid int64) (PlayerSummaryInfo, error) {
	v := make(url.Values)

	v.Set("key", SteamAPIKey)
	v.Set("steamids", strconv.Itoa(int(steamid)))

	r, err := http.Get(fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?%s", v.Encode()))
	if err != nil {
		return PlayerSummaryInfo{}, err
	}

	defer r.Body.Close()

	var rd GetPlayerSummariesResponse
	err = json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		return PlayerSummaryInfo{}, err
	}

	if len(rd.Response.Players) == 0 {
		return PlayerSummaryInfo{}, fmt.Errorf("no players returned")
	}

	return rd.Response.Players[0], nil
}

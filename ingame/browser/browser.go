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

package browser

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"reboxed/common"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
	"strings"
)

type Browser struct {
	InGame   bool
	SteamID  string
	MapName  string
	Search   string
	Category string
	Packages []common.Package
	PrevLink string
	NextLink string
}

const itemsPerPage = 50

//go:embed browser.tmpl
var tmpl string

var (
	categories = map[string]string{
		"mine":     "mine",
		"entities": "entity",
		"weapons":  "weapon",
		"props":    "prop",
		"saves":    "savemap",
		"maps":     "map",
	}
	t = template.Must(template.New("Browser").Funcs(template.FuncMap{"StripHTTPS": func(url string) string { s, _ := strings.CutPrefix(url, "https:"); return s }}).Parse(tmpl))
)

func Handle(w http.ResponseWriter, r *http.Request) {
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	var steamid string
	if r.Header.Get("TICKET") != "" {
		ticket, err := base64.StdEncoding.DecodeString(r.Header.Get("TICKET"))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode ticket value: %s", err))
			return
		}

		steamid, err = db.FetchSteamIDFromTicket(ticket)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to fetch steamid from ticket: %s", err))
			return
		}
	}

	c := category
	var author string
	if category == "mine" {
		c = "" // all categories
		author = steamid
	}

	list, err := db.FetchPackageList(c, author, r.URL.Query().Get("search"), (page-1)*itemsPerPage, itemsPerPage)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package list: %s", err))
		return
	}

	prev := fmt.Sprintf("?page=%d", page-1)
	if page <= 1 {
		prev = "#"
	}

	next := fmt.Sprintf("?page=%d", page+1)
	if len(list) < itemsPerPage {
		next = "#"
	}

	err = t.Execute(w, Browser{
		InGame:   strings.Contains(r.UserAgent(), "Valve"),
		SteamID:  steamid,
		MapName:  r.Header.Get("MAP"),
		Search:   r.URL.Query().Get("search"),
		Category: category,
		Packages: list,
		PrevLink: prev,
		NextLink: next,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}

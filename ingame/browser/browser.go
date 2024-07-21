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
	"fmt"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

var categories = map[string]string{
	"entities": "entity",
	"weapons":  "weapon",
	"props":    "prop",
	"saves":    "savemap",
	"maps":     "map",
}

const header = `<html>
<title>reboxed</title>
<style>
	body {margin: 0px; font-family: Helvetica; background-color: #36393D; color: #EEE;}
	a {color: #FFF; text-decoration: none;}
	a:hover {color: #0AF;}
	.nav {padding: 8px; background-color: #4096EE; height: 20px; border-bottom: 1px solid #90C6FE; box-shadow: 0px 16px 16px rgba(0, 0, 0, 0.1);}
	.nav a {margin: 20px; font-size: 20px; font-weight: bolder;}
	.logo h1 {margin: 0px; font-size: 20px; font-style: italic; float: right; color: #FFF;}
	.pagenav {float: right;}
	.pagenav a {margin: 8px; font-weight: bolder;}
	.content {padding: 16px 8px;}
	.item {margin-left: 2px; margin-right: 2px; display: inline-block; font-size: 11px; font-weight: bolder; width: 128px; height: 125px; text-align: center; text-shadow: 1px 1px 1px #000; text-overflow: ellipsis; overflow: hidden; white-space: nowrap; letter-spacing: -0.1px;}
	.item img {width: 128px; height: 100px;}
	.thumb {background-position: center;}
</style>
`

const itemsPerPage = 25

func Handle(w http.ResponseWriter, r *http.Request) {
	var category string
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	body := header

	body += `<div class="nav"><div class="logo"><h1>reboxed</h1></div>`

	if r.Header.Get("GMOD_VERSION") != "" { // ingame
		if category != "map" { // don't show nav buttons if browsing maps
			body += `<a href="/browse/entities">Entities</a><a href="/browse/weapons">Weapons</a><a href="/browse/props">Props</a><a href="/browse/saves">Saves</a>`
		}
	} else { // not ingame
		body += `<a href="/browse/entities">Entities</a><a href="/browse/weapons">Weapons</a><a href="/browse/props">Props</a><a href="/browse/saves">Saves</a><a href="/browse/maps">Maps</a>`
	}

	body += `</div><div class="content">`

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	list, err := db.FetchPackageListByTypePaged(category, (page-1)*itemsPerPage, itemsPerPage)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package list: %s", err))
		return
	}

	for _, pkg := range list {
		action := "spawn"
		if category == "map" {
			action = "install"
		}

		if category == "savemap" {
			if r.Header.Get("MAP") != "" && r.Header.Get("MAP") != pkg.Dataname {
				continue
			}
		}

		link := fmt.Sprintf(`href="garrysmod://%s/%s/%d/%d"`, action, category, pkg.ID, pkg.Revision)

		body += fmt.Sprintf(`<div class="item"><a %s><div class="thumb" style="background-image: url(//image.reboxed.fun/%d_thumb_128.png), url(//image.reboxed.fun/no_thumb_128.png);"><img src="//image.reboxed.fun/overlay_128.png"></div>%s</a></div>`, link, pkg.ID, pkg.Name)
	}

	body += "</div>"

	body += fmt.Sprintf(`<div class="pagenav"><a href="?page=%d">Previous</a>%d<a href="?page=%d">Next</a></div>`, page-1, page, page+1)

	body += "</html>"

	w.Write([]byte(body))
}

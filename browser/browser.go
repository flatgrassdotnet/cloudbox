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
)

var categories = map[string]string{
	"entities": "entity",
	"weapons":  "weapon",
	"props":    "prop",
	"saves":    "savemap",
	"maps":     "map",
}

const header = `<html>
<style>
	body {margin: 0px; font-family: Helvetica; background-color: #36393D; color: #EEE;}
	a {color: #FFF; text-decoration: none;}
	a:hover {color: #0AF;}
	.nav {padding: 8px; background-color: #4096EE; height: 20px; text-align: center; border-bottom: 1px solid #90C6FE; box-shadow: 0px 16px 16px rgba(0, 0, 0, 0.1);}
	.nav a {margin: 20px; font-size: 20px; font-weight: bolder;}
	.logo h1 {margin: 0px; font-size: 20px; font-style: italic; float: left; color: #FFF;}
	.logo img {margin: 0px; width: 20px; height: 20px; float: left;}
	.content {padding: 16px 8px;}
	.item {margin-left: 2px; margin-right: 2px; display: inline-block; font-size: 11px; font-weight: bolder; width: 128px; height: 125px; text-align: center; text-shadow: 1px 1px 1px #000; text-overflow: ellipsis; overflow: hidden; white-space: nowrap; letter-spacing: -0.1px;}
	.item img {width: 128px; height: 100px;}
	.thumb {background-position: center;}
</style>
`

func Handle(w http.ResponseWriter, r *http.Request) {
	var category string
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	body := header

	body += `<div class="nav"><div class="logo"><img src="//img.reboxed.fun/logo.png"><h1>reboxed</h1></div>`

	if category != "map" {
		body += `<a href="/browse/entities">Entities</a><a href="/browse/weapons">Weapons</a><a href="/browse/saves">Saves</a>`
	}

	body += `</div><div class="content">`

	list, err := db.FetchPackageListByType(category)
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

		body += fmt.Sprintf(`<div class="item"><a %s><div class="thumb" style="background-image: url(//img.reboxed.fun/%d_thumb_128.png), url(//img.reboxed.fun/no_thumb_128.png);"><img src="//img.reboxed.fun/overlay_128.png"></div>%s</a></div>`, link, pkg.ID, pkg.Name)
	}

	body += "</div></html>"

	w.Write([]byte(body))
}

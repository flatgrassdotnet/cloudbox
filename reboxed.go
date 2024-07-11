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

package main

import (
	"flag"
	"log"
	"net/http"
	"reboxed/api/packages"
	"reboxed/browser"
	"reboxed/db"
	"reboxed/stats"
	"reboxed/toyboxapi"
	"reboxed/utils"
)

func main() {
	dbuser := flag.String("dbuser", "reboxed", "database user's name")
	dbpass := flag.String("dbpass", "", "database user's password")
	dbaddr := flag.String("dbaddr", "localhost", "database server address")
	dbname := flag.String("dbname", "reboxed", "database name")
	apikey := flag.String("apikey", "", "steam web api key")
	flag.Parse()

	err := db.Init(*dbuser, *dbpass, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("failed to init database: %s", err)
	}

	utils.WebAPIKey = *apikey

	// browser
	http.HandleFunc("GET /browse/{category}/", browser.Handle)

	// static content - using nginx now
	//http.Handle("GET cdn.reboxed.fun/", http.FileServer(http.Dir("data/content")))
	//http.Handle("GET img.reboxed.fun/", http.FileServer(http.Dir("data/img")))

	// api
	http.HandleFunc("GET api.reboxed.fun/packages/list", packages.List)
	http.HandleFunc("GET api.reboxed.fun/packages/get", packages.Get)

	// stats.garrysmod.com
	http.HandleFunc("GET stats.garrysmod.com/API/mapload_001/", stats.MapLoad)

	// toyboxapi.garrysmod.com
	http.HandleFunc("POST toyboxapi.garrysmod.com/auth_003/", toyboxapi.Auth)
	http.HandleFunc("POST toyboxapi.garrysmod.com/error_003/", toyboxapi.Error)
	http.HandleFunc("GET toyboxapi.garrysmod.com/getinstall_003/", toyboxapi.GetPackage)
	http.HandleFunc("GET toyboxapi.garrysmod.com/getscript_003/", toyboxapi.GetPackage)
	//http.HandleFunc("POST toyboxapi.garrysmod.com/upload_003/", toyboxapi.Upload)

	// toybox.garrysmod.com
	//http.HandleFunc("GET toybox.garrysmod.com/API/publishsave_002/", toybox.PublishSave)

	// redirects
	http.HandleFunc("GET toybox.garrysmod.com/ingame/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "//toybox.garrysmod.com/browse/entities", http.StatusSeeOther)
	})
	http.HandleFunc("GET toybox.garrysmod.com/IG/maps/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "//toybox.garrysmod.com/browse/maps", http.StatusSeeOther)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error while serving: %s", err)
	}
}

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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/api/auth"
	"github.com/flatgrassdotnet/cloudbox/api/content"
	"github.com/flatgrassdotnet/cloudbox/api/packages"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/ingame/stats"
	"github.com/flatgrassdotnet/cloudbox/ingame/toyboxapi"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func main() {
	dbuser := flag.String("dbuser", "cloudbox", "database user's name")
	dbpass := flag.String("dbpass", "", "database user's password")
	dbaddr := flag.String("dbaddr", "localhost", "database server address")
	dbname := flag.String("dbname", "cloudbox", "database name")
	apikey := flag.String("apikey", "", "steam web api key")
	statswebhook := flag.String("statswebhook", "", "discord stats webhook url")
	savewebhook := flag.String("savewebhook", "", "discord save webhook url")
	port := flag.Int("port", 80, "web server listen port")
	flag.Parse()

	err := db.Init(*dbuser, *dbpass, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("failed to init database: %s", err)
	}

	utils.SteamAPIKey = *apikey
	utils.DiscordStatsWebhookURL = *statswebhook
	utils.DiscordSaveWebhookURL = *savewebhook

	// cloudbox api
	http.HandleFunc("GET api.cl0udb0x.com/auth/getid", auth.GetID)
	http.HandleFunc("GET api.cl0udb0x.com/packages/list", packages.List)
	http.HandleFunc("GET api.cl0udb0x.com/packages/get", packages.Get)
	http.HandleFunc("GET api.cl0udb0x.com/packages/getgma", packages.GetGMA)
	http.HandleFunc("GET api.cl0udb0x.com/packages/publishsave", packages.PublishSave)
	http.HandleFunc("GET api.cl0udb0x.com/content/get", content.Get)

	// stats.garrysmod.com
	http.HandleFunc("GET /API/mapload_001/", stats.MapLoad)

	// toyboxapi.garrysmod.com
	http.HandleFunc("POST /auth_003/", toyboxapi.Auth)
	http.HandleFunc("GET /error_003/", toyboxapi.Error)
	http.HandleFunc("GET /getinstall_003/", toyboxapi.GetPackage)
	http.HandleFunc("GET /getscript_003/", toyboxapi.GetPackage)
	http.HandleFunc("POST /upload_003/", toyboxapi.Upload)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error while serving: %s", err)
	}
}

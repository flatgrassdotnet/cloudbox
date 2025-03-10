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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/api/auth"
	"github.com/flatgrassdotnet/cloudbox/api/content"
	"github.com/flatgrassdotnet/cloudbox/api/news"
	"github.com/flatgrassdotnet/cloudbox/api/packages"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/ingame/publishsave"
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
	http.HandleFunc("GET /auth/getid", auth.GetID)
	http.HandleFunc("GET /news/list", news.List)
	http.HandleFunc("GET /packages/list", packages.List)
	http.HandleFunc("GET /packages/get", packages.Get)
	http.HandleFunc("GET /packages/getscript", packages.GetScript)
	http.HandleFunc("GET /packages/getgma", packages.GetGMA)
	http.HandleFunc("GET /content/get", content.Get)
	http.HandleFunc("GET /content/getzip", content.GetZIP)
	http.HandleFunc("GET /content/fastdl", content.FastDL)

	// stats.garrysmod.com
	http.HandleFunc("GET stats.garrysmod.com/API/mapload_001/", stats.MapLoad) // v102 - v142

	// toyboxapi.garrysmod.com
	// auth
	http.HandleFunc("GET toyboxapi.garrysmod.com/auth_001/", toyboxapi.Auth)  // v104 - v106
	http.HandleFunc("POST toyboxapi.garrysmod.com/auth_002/", toyboxapi.Auth) // v107 - v133
	http.HandleFunc("POST toyboxapi.garrysmod.com/auth_003/", toyboxapi.Auth) // v134 - v142

	// getinstall
	http.HandleFunc("GET toyboxapi.garrysmod.com/getinstall_003/", toyboxapi.GetPackage) // v134 - v142

	// getscript
	http.HandleFunc("GET toyboxapi.garrysmod.com/getscript_001/", toyboxapi.GetPackage) // v100 - v133
	http.HandleFunc("GET toyboxapi.garrysmod.com/getscript_003/", toyboxapi.GetPackage) // v134 - v142

	// upload
	http.HandleFunc("POST toyboxapi.garrysmod.com/upload_001/", toyboxapi.Upload) // v109 - v133
	http.HandleFunc("POST toyboxapi.garrysmod.com/upload_003/", toyboxapi.Upload) // v134 - v142

	// error
	http.HandleFunc("GET toyboxapi.garrysmod.com/error_001/", toyboxapi.Error) // v98 - v133
	http.HandleFunc("GET toyboxapi.garrysmod.com/error_003/", toyboxapi.Error) // v134 - v142

	// publishsave
	http.HandleFunc("GET toyboxapi.garrysmod.com/publishsave_001/", publishsave.Save) // v106 - v108
	http.HandleFunc("GET toyboxapi.garrysmod.com/publishsave_002/", publishsave.Save) // v109 - v142

	http.HandleFunc("POST toyboxapi.garrysmod.com/publishsave_001/", publishsave.Publish) // v106 - v108
	http.HandleFunc("POST toyboxapi.garrysmod.com/publishsave_002/", publishsave.Publish) // v109 - v142

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error while serving: %s", err)
	}
}

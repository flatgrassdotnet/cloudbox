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

package publishsave

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"net/http"
	"strconv"

	"github.com/blezek/tga"
	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

var tp = template.Must(template.New("publish.html").ParseGlob("data/templates/publishsave/*.html"))

func Publish(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse form data: %s", err))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	sid, err := strconv.Atoi(r.URL.Query().Get("sid"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse sid value: %s", err))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "No Name"
	}

	desc := r.URL.Query().Get("desc")

	ticket, err := base64.StdEncoding.DecodeString(r.Header.Get("TICKET"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode ticket value: %s", err))
		return
	}

	steamid, err := db.FetchSteamIDFromTicket(ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch steamid from ticket: %s", err))
		return
	}

	save, err := db.FetchUpload(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch upload: %s", err))
		return
	}

	pkgID, err := db.InsertPackage(common.Package{Type: "savemap", Name: name, Dataname: save.Metadata, Author: steamid, Description: desc, Data: save.Data})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert package: %s", err))
		return
	}

	for _, include := range save.Includes {
		rev, err := db.FetchPackageLatestRevision(include)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to fetch package latest revision: %s", err))
			return
		}

		// save revision should always be 1
		_, err = db.InsertPackageInclude(pkgID, 1, include, rev)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to insert package include: %s", err))
			return
		}
	}

	err = db.DeleteUpload(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to delete upload: %s", err))
		return
	}

	thumb, err := db.FetchUpload(sid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch upload: %s", err))
		return
	}

	img, err := tga.Decode(bytes.NewReader(thumb.Data))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode thumbnail tga: %s", err))
		return
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to encode thumbnail png: %s", err))
		return
	}

	err = db.DeleteUpload(sid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to delete upload: %s", err))
		return
	}

	err = db.PutThumbnail(id, buf)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to upload thumbnail: %s", err))
		return
	}

	err = tp.Execute(w, nil)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}

	// webhook related
	s, err := utils.GetPlayerSummaries(steamid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	err = utils.SendDiscordMessage(utils.DiscordSaveWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title:       name,
			Description: desc,
			Color:       0xB8E3FF,
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s[0].PersonaName,
				IconURL: s[0].Avatar,
			},
			Image: utils.DiscordWebhookEmbedImage{
				URL: fmt.Sprintf("https://img.cl0udb0x.com/%d_thumb_128.png", pkgID),
			},
		},
		},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to send discord webhook message: %s", err))
		return
	}
}

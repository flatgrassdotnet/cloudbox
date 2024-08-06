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

package publishsave

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"net/http"
	"os"
	"reboxed/db"
	"reboxed/utils"
	"strconv"

	"github.com/blezek/tga"
)

type PublishSavePost struct{}

//go:embed post.tmpl
var tmplPost string

var tp = template.Must(template.New("PublishSavePost").Parse(tmplPost))

func Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse form data: %s", err))
		return
	}

	name := r.PostForm.Get("name")
	if name == "" {
		name = "No Name"
	}

	desc := r.PostForm.Get("desc")

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

	pkgID, err := db.InsertPackage("savemap", name, save.Metadata, steamid, desc, save.Data)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert package: %s", err))
		return
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

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to encode thumbnail png: %s", err))
		return
	}

	err = os.WriteFile(fmt.Sprintf("data/img/%d_thumb_128.png", pkgID), buf.Bytes(), 0644)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to write thumbnail: %s", err))
		return
	}

	err = db.DeleteUpload(sid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to delete upload: %s", err))
		return
	}

	err = tp.Execute(w, PublishSavePost{})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}

	// webhook related
	s, err := utils.GetPlayerSummary(steamid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	err = utils.SendDiscordMessage(utils.DiscordSaveWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title:       name,
			Description: desc,
			Color:       10607359, // #A1DAFF
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s.PersonaName,
				IconURL: s.Avatar,
			},
			Image: utils.DiscordWebhookEmbedImage{
				URL: fmt.Sprintf("https://img.reboxed.fun/%d_thumb_128.png", pkgID),
			},
		},
		},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to send discord webhook message: %s", err))
		return
	}
}

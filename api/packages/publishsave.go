package packages

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/blezek/tga"
	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func PublishSave(w http.ResponseWriter, r *http.Request) {
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

	ticket, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("ticket"))
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

	if save.Include != "" {
		for _, inc := range strings.Split(save.Include, ",") {
			// usually means it hit the end but maybe not
			if inc == "" {
				continue
			}

			i, err := strconv.Atoi(inc)
			if err != nil {
				utils.WriteError(w, r, fmt.Sprintf("failed to parse inc value: %s", err))
				return
			}

			rev, err := db.FetchPackageLatestRevision(i)
			if err != nil {
				utils.WriteError(w, r, fmt.Sprintf("failed to fetch package latest revision: %s", err))
				return
			}

			// save revision should always be 1 unless something has gone horribly wrong
			_, err = db.InsertPackageInclude(pkgID, 1, i, rev)
			if err != nil {
				utils.WriteError(w, r, fmt.Sprintf("failed to insert package include: %s", err))
				return
			}
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

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to encode thumbnail png: %s", err))
		return
	}

	err = os.WriteFile(filepath.Join("data", "img", strconv.Itoa(pkgID)+"_thumb_128.png"), buf.Bytes(), 0644)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to write thumbnail: %s", err))
		return
	}

	err = db.DeleteUpload(sid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to delete upload: %s", err))
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
			Color:       0xB8E3FF,
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s.PersonaName,
				IconURL: s.Avatar,
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

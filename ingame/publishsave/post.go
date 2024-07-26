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

	pkgID, err := db.InsertPackage("savemap", name, save.Metadata, steamid, r.PostForm.Get("desc"), save.Data)
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
}

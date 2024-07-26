package toyboxapi

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	// steamid64
	steamid, err := strconv.Atoi(r.URL.Query().Get("steamid"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse steamid value: %s", err))
		return
	}

	// "save" or "save_image"
	uploadType := r.URL.Query().Get("type")
	if uploadType != "save" && uploadType != "save_image" {
		utils.WriteError(w, r, "invalid upload type")
		return
	}

	// metadata
	meta := r.URL.Query().Get("meta")

	// unknown
	inc := r.URL.Query().Get("inc")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to read post body: %s", err))
		return
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	_, err = base64.StdEncoding.Decode(decoded, body)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse post body: %s", err))
		return
	}

	id, err := db.InsertUpload(steamid, utils.Upload{Type: uploadType, Metadata: meta, Include: inc, Data: decoded})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert upload: %s", err))
		return
	}

	w.Write([]byte(strconv.Itoa(id)))
}

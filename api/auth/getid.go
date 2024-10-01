package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func GetID(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(steamid))
}

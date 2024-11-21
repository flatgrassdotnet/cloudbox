package news

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func List(w http.ResponseWriter, r *http.Request) {
	entries, err := db.FetchNewsEntries()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch news entries: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to encode response: %s", err))
		return
	}
}

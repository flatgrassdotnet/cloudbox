package content

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func FastDL(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		utils.WriteError(w, r, "missing file value")
		return
	}

	id, rev, err := db.FetchFileInfoFromPath(strings.TrimPrefix(file, "/"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to fetch file info: %s", err))
		return
	}

	data, err := utils.GetContentFile(id, rev)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get content file data: %s", err))
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

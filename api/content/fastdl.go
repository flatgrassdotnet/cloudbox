package content

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
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

	id, err := db.FetchFileInfoFromPath(strings.TrimPrefix(file, "/"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to fetch file info: %s", err))
		return
	}

	f, err := db.GetContentFile(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open content file: %s", err))
		return
	}

	defer f.Body.Close()

	w.Header().Set("Content-Length", strconv.Itoa(int(*f.ContentLength)))
	io.Copy(w, f.Body)
}

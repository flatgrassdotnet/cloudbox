package content

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/utils"
)

func GetZIP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	f, err := utils.GetContentFile(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open content file for reading: %s", err))
		return
	}

	defer f.Close()

	// GM12 won't show download progress without Content-Length
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)

	file, err := zw.Create("file")
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open zip file for writing: %s", err))
		return
	}

	io.Copy(file, f)

	zw.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	io.Copy(w, buf)
}

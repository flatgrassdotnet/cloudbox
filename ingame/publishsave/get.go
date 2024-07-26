package publishsave

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"reboxed/utils"
	"strconv"
)

type PublishSaveGet struct {
	ID  int
	SID int
}

//go:embed get.tmpl
var tmplGet string

var tg = template.Must(template.New("PublishSaveGet").Parse(tmplGet))

func Get(w http.ResponseWriter, r *http.Request) {
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

	err = tg.Execute(w, PublishSaveGet{
		ID:  id,
		SID: sid,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}

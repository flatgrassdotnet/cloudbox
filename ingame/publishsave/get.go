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

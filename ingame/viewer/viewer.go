/*
	cloudbox - the toybox server emulator
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

package viewer

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

var t = template.Must(template.New("viewer.html").ParseFiles("data/templates/viewer/viewer.html"))

func Handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse id: %s", err), http.StatusBadRequest)
		return
	}

	rev, err := db.FetchPackageLatestRevision(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch rev: %s", err), http.StatusInternalServerError)
		return
	}

	pkg, err := db.FetchPackage(id, rev)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch package data: %s", err), http.StatusInternalServerError)
	}

	err = t.Execute(w, pkg)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}

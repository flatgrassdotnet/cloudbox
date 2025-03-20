/*
	cloudbox - the toybox server emulator
	Copyright (C) 2024-2025  patapancakes <patapancakes@pagefault.games>

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

package common

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Package struct {
	ID          int       `json:"id"`
	Revision    int       `json:"rev"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Dataname    string    `json:"dataname,omitempty"`
	Author      string    `json:"author,omitempty"`
	AuthorName  string    `json:"authorname,omitempty"`
	AuthorIcon  string    `json:"authoricon,omitempty"`
	Description string    `json:"description,omitempty"`
	Data        []byte    `json:"data,omitempty"`
	Content     []Content `json:"content,omitempty"`
	Includes    []Include `json:"includes,omitempty"`
	Uploaded    time.Time `json:"uploaded,omitempty"`

	Downloads int `json:"downloads,omitempty"`
	Favorites int `json:"favorites,omitempty"`
	Goods     int `json:"goods,omitempty"`
	Bads      int `json:"bads,omitempty"`
}

type Content struct {
	ID    int    `json:"id"`
	Path  string `json:"path"`
	Size  int    `json:"size"`
	PSize int    `json:"psize"`
}

type Include struct {
	ID       int    `json:"id"`
	Revision int    `json:"rev"`
	Type     string `json:"type"`
}

func (pkg Package) Marshal() []byte {
	script := make(VDF)

	script["scriptid"] = pkg.ID
	script["revision"] = pkg.Revision
	script["type"] = pkg.Type
	script["dataname"] = pkg.Dataname
	script["name"] = pkg.Name

	// maps have extra stuff
	if pkg.Type == "map" {
		script["uid"] = fmt.Sprintf("map_%d", pkg.ID)
		script["luamenu_installed"] = "OnMapDownloaded();"
		script["luamenu_action"] = fmt.Sprintf("OnMapSelected('%s');", pkg.RealMapName())
	}

	if len(pkg.Content) != 0 {
		content := make(VDF)

		for _, c := range pkg.Content {
			item := make(VDF)

			item["id"] = c.ID
			item["rev"] = 1
			item["name"] = c.Path
			item["url"] = fmt.Sprintf("http://api.cl0udb0x.com/content/getzip?id=%d", c.ID)
			item["size"] = c.PSize

			content[fmt.Sprintf("content_%d", c.ID)] = item
		}

		script["content"] = content
	}

	if len(pkg.Includes) != 0 {
		includes := make(VDF)

		for _, i := range pkg.Includes {
			item := make(VDF)

			item["id"] = i.ID
			item["rev"] = i.Revision
			item["type"] = i.Type

			includes[fmt.Sprintf("include_%d", i.ID)] = item
		}

		script["includes"] = includes
	}

	root := make(VDF)
	root["script"] = script

	if len(pkg.Data) != 0 {
		return append([]byte(root.Marshal()), pkg.Data...)
	}

	return []byte(root.Marshal())
}

func (pkg Package) RealMapName() string {
	mapname := pkg.Name
	for _, c := range pkg.Content {
		if filepath.Ext(c.Path) != ".bsp" {
			continue
		}

		mapname = strings.TrimSuffix(filepath.Base(c.Path), filepath.Ext(c.Path))
	}

	return mapname
}

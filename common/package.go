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
	ID       int       `json:"id"`
	Revision int       `json:"rev"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Dataname string    `json:"dataname,omitempty"`
	Content  []Content `json:"content,omitempty"`
	Includes []Include `json:"includes,omitempty"`
	Data     []byte    `json:"data,omitempty"`

	// only used by Install packages

	LuaMenuInstalled   string `json:"-"`
	LuaMenuAction      string `json:"-"`
	LuaClientInstalled string `json:"-"`
	LuaClientAction    string `json:"-"`
	LuaServerInstalled string `json:"-"`
	LuaServerAction    string `json:"-"`

	// metadata

	Author      string    `json:"author,omitempty"`
	AuthorName  string    `json:"authorname,omitempty"`
	AuthorIcon  string    `json:"authoricon,omitempty"`
	Description string    `json:"description,omitempty"`
	Uploaded    time.Time `json:"uploaded,omitempty"`

	Downloads int `json:"downloads,omitempty"`
	Favorites int `json:"favorites,omitempty"`
	Goods     int `json:"goods,omitempty"`
	Bads      int `json:"bads,omitempty"`
}

type Content struct {
	ID       int    `json:"id"`
	Revision int    `json:"rev"` // not stored by cloudbox, always 1
	Path     string `json:"path"`
	Size     int    `json:"size"`  // raw size
	PSize    int    `json:"psize"` // compressed size
}

type Include struct {
	ID       int    `json:"id"`
	Revision int    `json:"rev"`
	Type     string `json:"type"`
}

func (pkg Package) Marshal(install bool) []byte {
	script := make(VDF)

	script["scriptid"] = pkg.ID
	script["revision"] = pkg.Revision
	script["type"] = pkg.Type
	script["dataname"] = pkg.Dataname
	script["name"] = pkg.Name

	if install {
		script["uid"] = pkg.UID()
		if pkg.LuaMenuInstalled != "" {
			script["luamenu_installed"] = pkg.LuaMenuInstalled
		}
		if pkg.LuaMenuAction != "" {
			script["luamenu_action"] = pkg.LuaMenuAction
		}
		if pkg.LuaClientInstalled != "" {
			script["luaclient_installed"] = pkg.LuaClientInstalled
		}
		if pkg.LuaClientAction != "" {
			script["luaclient_action"] = pkg.LuaClientAction
		}
		if pkg.LuaServerInstalled != "" {
			script["luaserver_installed"] = pkg.LuaServerInstalled
		}
		if pkg.LuaServerAction != "" {
			script["luaserver_action"] = pkg.LuaServerAction
		}
	}

	if len(pkg.Content) != 0 {
		content := make(VDF)

		for _, c := range pkg.Content {
			item := make(VDF)

			item["id"] = c.ID
			item["rev"] = c.Revision
			item["name"] = c.Path
			item["url"] = fmt.Sprintf("http://api.cl0udb0x.com/content/getzip?id=%d", c.ID)
			item["size"] = c.PSize

			// name doesn't matter
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

			// name doesn't matter
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

func (pkg Package) BSPName() string {
	for _, c := range pkg.Content {
		if filepath.Ext(c.Path) == ".bsp" {
			return strings.TrimSuffix(filepath.Base(c.Path), filepath.Ext(c.Path))
		}
	}

	return pkg.Name
}

// used by Install packages only
func (pkg Package) UID() string {
	return fmt.Sprintf("%s_%d", pkg.Type, pkg.ID)
}

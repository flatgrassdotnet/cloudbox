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

package common

import "fmt"

type VDF map[string]any

func (vdf VDF) Marshal() string {
	var output string

	vdf.encode(&output)

	return output + "\n"
}

// very crude VDF string encoder specifically targeting toybox's implementation
// all values become strings, even though VDF supports non-string values
func (vdf VDF) encode(output *string) {
	for k, v := range vdf {
		switch data := v.(type) {
		case int:
			*output += fmt.Sprintf("\"%s\"\t\"%d\"\n", k, data)
		case string:
			*output += fmt.Sprintf("\"%s\"\t\"%s\"\n", k, data)
		case VDF:
			*output += fmt.Sprintf("\"%s\"\n{\n", k)
			data.encode(output)
			*output += "}\n"
		}
	}
}

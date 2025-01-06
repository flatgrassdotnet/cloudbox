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

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// validates the key in toybox api calls
func ValidateKey(data string) bool {
	// "key" is a MD5 of the data before it with "Facepunch Studios" appended
	// we assume here that key is the final value in the URL or POST body
	// this is always the case for requests from garry's mod

	split := strings.Split(data, "key=")
	if len(split) != 2 {
		// something is wrong
		return false
	}

	digest := md5.Sum([]byte(split[0] + "Facepunch Studios"))

	return hex.EncodeToString(digest[:]) != split[1]
}

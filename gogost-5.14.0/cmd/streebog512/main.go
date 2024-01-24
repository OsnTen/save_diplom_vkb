// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2024 Sergey Matveev <stargrave@stargrave.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy_db of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Command-line 34.11-2012 512-bit hash function.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"go.cypherpunks.ru/gogost/v5"
	"go.cypherpunks.ru/gogost/v5/gost34112012512"
)

var (
	version = flag.Bool("version", false, "Print version information")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Println(gogost.Version)
		return
	}
	h := gost34112012512.New()
	if _, err := io.Copy(h, os.Stdin); err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}

// Quicknote stores and searches tens of thousands of short notes.
//
// Copyright (C) 2017  Andrew Miller <amiller@amilx.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package db

import (
	"errors"

	"github.com/anmil/quicknote"
	"github.com/anmil/quicknote/db/postgres"
	"github.com/anmil/quicknote/db/sqlite"
)

// ErrProviderNotSupported database provider given is not supported
var ErrProviderNotSupported = errors.New("Unsupported database provider")

// NewDatabase returns a new database for the given provider
func NewDatabase(provider string, options ...string) (quicknote.DB, error) {
	switch provider {
	case "sqlite":
		return sqlite.NewDatabase(options...)
	case "postgres":
		return postgres.NewDatabase(options...)
	default:
		return nil, ErrProviderNotSupported
	}
}

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

package quicknote

import (
	"fmt"
	"time"
)

// Book is a collection of notes
type Book struct {
	ID       int64
	Created  time.Time
	Modified time.Time

	Name string
}

// NewBook returns a new Book
func NewBook() *Book {
	return &Book{}
}

func (b *Book) String() string {
	return fmt.Sprintf("<Book ID: %d Name: %s>", b.ID, b.Name)
}

type Books []*Book

func (b Books) Len() int {
	return len(b)
}

func (b Books) Less(i, j int) bool {
	return b[i].ID < b[j].ID
}

func (b Books) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
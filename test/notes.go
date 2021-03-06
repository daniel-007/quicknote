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

package test

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/anmil/quicknote"
)

var notesJSON = `[
  {
    "id": 603,
    "created": "2017-03-25T21:35:27.287881752-04:00",
    "modified": "2017-03-25T21:35:27.287881829-04:00",
    "type": "basic",
    "title": "This is test 1 of the basic parser",
    "body": "#basic #test #parser\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor quis nibh.",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser"
    ]
  },
  {
    "id": 604,
    "created": "2017-03-25T21:35:27.293783239-04:00",
    "modified": "2017-03-25T21:35:27.29378334-04:00",
    "type": "basic",
    "title": "This is #test 2 of the #basic #parser",
    "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor quis nibh.",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser"
    ]
  },
  {
    "id": 605,
    "created": "2017-03-25T21:35:27.305349683-04:00",
    "modified": "2017-03-25T21:35:27.305349813-04:00",
    "type": "basic",
    "title": "This is #test 2 of the #basic #parser",
    "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nNulla tincidunt diam eu purus laoreet condimentum. Duis\ntempus, turpis vitae varius ullamcorper, sapien erat\ncursus lacus, et lacinia ligula dolor #quis nibh.#",
    "book": "test",
    "tags": [
      "basic",
      "test",
      "parser",
      "quis"
    ]
  }
]`

type JsonNote struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Book     string    `json:"book"`
	Tags     []string  `json:"tags"`
}

func init() {
	noteBooks = make(map[string]*quicknote.Book)
	noteTags = make(map[string]*quicknote.Tag)
}

func GetTestNotes() quicknote.Notes {
	return GetTestNotesCust(notesJSON)
}

func GetTestNotesCust(s string) quicknote.Notes {
	reader := strings.NewReader(s)
	dec := json.NewDecoder(reader)

	jsonNotes := make([]*JsonNote, 0)
	err := dec.Decode(&jsonNotes)
	if err != nil {
		panic(err)
	}

	testNotes := make(quicknote.Notes, len(jsonNotes))
	for idx, jn := range jsonNotes {
		tags := make(quicknote.Tags, len(jn.Tags))
		for i, t := range jn.Tags {
			tags[i] = getTag(t)
		}

		n := quicknote.NewNote()
		n.ID = jn.ID
		n.Created = jn.Created
		n.Modified = jn.Modified
		n.Type = jn.Type
		n.Title = jn.Title
		n.Body = jn.Body
		n.Book = getBook(jn.Book)
		n.Tags = tags

		testNotes[idx] = n
	}

	return testNotes
}

func CheckNotes(t *testing.T, notes1, notes2 quicknote.Notes) {
	nnNotes := quicknote.Notes{}
	for _, t := range notes1 {
		nnNotes = append(nnNotes, t)
	}
	sort.Sort(nnNotes)

	nNotes := quicknote.Notes{}
	for _, t := range notes2 {
		nNotes = append(nNotes, t)
	}
	sort.Sort(nNotes)

	if !NoteSliceEq(nnNotes, nNotes) {
		t.Fatal("Did not received the corrected Notes")
	}
}

func NoteSliceEq(a, b quicknote.Notes) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Title != b[i].Title {
			return false
		}
		if a[i].Body != b[i].Body {
			return false
		}
	}
	return true
}

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

package sqlite

import (
	"testing"

	"github.com/anmil/quicknote"
	"github.com/anmil/quicknote/test"
)

func TestCreateNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)
}

func TestGetNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	notes := test.GetTestNotes()
	saveNotes(t, db, notes)

	getNoteByID(t, db, notes[0])
	getNoteByNote(t, db, notes[0])
	getNotesByID(t, db, notes)
	getNotesByBook(t, db, notes)
	getNotesAll(t, db, notes)
}

func TestEditNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	bk := quicknote.NewBook()
	bk.Name = "NewBook"

	err := db.CreateBook(bk)
	if err != nil {
		t.Fatal(err)
	}

	var ids []int64
	notes := test.GetTestNotes()

	for _, n := range notes {
		saveNote(t, db, n)
		ids = append(ids, n.ID)
		n.Book = bk
	}

	if err := db.EditNoteByIDBook(ids, bk); err != nil {
		t.Fatal(err)
	}

	getNotesByBook(t, db, notes)
}

func TestEditNoteBookSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)

	n.Title = "New title"
	if err := db.EditNote(n); err != nil {
		t.Fatal(err)
	}

	getNoteByID(t, db, n)
}

func TestDeleteNoteSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	getNoteByID(t, db, n)

	if err := db.DeleteNote(n); err != nil {
		t.Fatal(err)
	}

	if nn, err := db.GetNoteByID(n.ID); err != nil {
		t.Fatal(err)
	} else if nn != nil {
		t.Fatal("Expected nil, got a note")
	}
}

func saveNotes(t *testing.T, db *Database, notes quicknote.Notes) {
	for _, n := range notes {
		saveNote(t, db, n)
	}
}

func saveNote(t *testing.T, db *Database, n *quicknote.Note) {
	if bk, err := db.GetBookByName(n.Book.Name); err != nil {
		t.Fatal(err)
	} else if bk == nil {
		if err := db.CreateBook(n.Book); err != nil {
			t.Fatal(err)
		}
	}

	for _, tag := range n.Tags {
		if bk, err := db.GetTagByName(tag.Name); err != nil {
			t.Fatal(err)
		} else if bk == nil {
			if err := db.CreateTag(tag); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := db.CreateNote(n); err != nil {
		t.Fatal(err)
	}
}

func getNoteByID(t *testing.T, db *Database, n *quicknote.Note) {
	if nn, err := db.GetNoteByID(n.ID); err != nil {
		t.Fatal(err)
	} else if nn == nil {
		t.Fatal("Expected 1 note, got nil")
	} else if nn.ID != n.ID {
		t.Fatalf("Expected note with ID %d, got %d", n.ID, nn.ID)
	} else {
		test.CheckTags(t, nn.Tags, n.Tags)
	}
}

func getNoteByNote(t *testing.T, db *Database, n *quicknote.Note) {
	nn := quicknote.NewNote()
	nn.Book = n.Book
	nn.Type = n.Type
	nn.Title = n.Title
	nn.Body = n.Body
	if err := db.GetNoteByNote(nn); err != nil {
		t.Fatal(err)
	} else if nn == nil {
		t.Fatal("Expected 1 note, got nil")
	} else if nn.ID != n.ID {
		t.Fatalf("Expected note with ID %d, got %d", n.ID, nn.ID)
	} else if !nn.Created.Equal(n.Created) {
		t.Fatalf("Expected note with Created %s, got %s", n.Created, nn.Created)
	} else if !nn.Modified.Equal(n.Modified) {
		t.Fatalf("Expected note with Modified %s, got %s", n.Modified, nn.Modified)
	}
}

func getNotesByID(t *testing.T, db *Database, notes quicknote.Notes) {
	var ids []int64
	for _, n := range notes {
		ids = append(ids, n.ID)
	}

	if nn, err := db.GetNotesByIDs(ids); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}

func getNotesByBook(t *testing.T, db *Database, notes quicknote.Notes) {
	if nn, err := db.GetAllBookNotes(notes[0].Book, "modified", "asc"); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}

func getNotesAll(t *testing.T, db *Database, notes quicknote.Notes) {
	if nn, err := db.GetAllNotes("modified", "asc"); err != nil {
		t.Fatal(err)
	} else if len(nn) != len(notes) {
		t.Fatalf("Expected %d notes, got %d", len(notes), len(nn))
	} else {
		test.CheckNotes(t, nn, notes)
		for i := 0; i < len(nn); i++ {
			test.CheckTags(t, nn[i].Tags, notes[i].Tags)
		}
	}
}

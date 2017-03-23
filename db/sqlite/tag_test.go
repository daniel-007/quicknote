package sqlite

import (
	"testing"

	"github.com/anmil/quicknote/note"
	"github.com/anmil/quicknote/test"
)

func TestGetTagByNameSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	if tag, err := db.GetTagByName(n.Tags[0].Name); err != nil {
		t.Fatal(err)
	} else if tag == nil {
		t.Fatal("Expected 1 tag, got nil")
	} else if tag.Name != n.Tags[0].Name {
		t.Fatalf("Expected tag %s, got %s", n.Tags[0].Name, tag.Name)
	}
}

func TestGetOrCreateTagSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	tag1, err := db.GetOrCreateTagByName("NewTag")
	if err != nil {
		t.Fatal(err)
	} else if tag1 == nil {
		t.Fatal("Expected 1 tag, got nil")
	}

	if tag2, err := db.GetTagByName(tag1.Name); err != nil {
		t.Fatal(err)
	} else if tag2 == nil {
		t.Fatal("Expected 1 tag, got nil")
	} else if tag2.Name != tag1.Name {
		t.Fatalf("Expected tag %s, got %s", tag1.Name, tag2.Name)
	}
}

func TestGetTagsSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	notes := test.GetTestNotes()
	saveNotes(t, db, notes)

	tags, err := db.GetAllBookTags(notes[0].Book)
	if err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, test.AllTags, tags)
	}

	if tags, err = db.GetAllTags(); err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, test.AllTags, tags)
	}
}

func TestLoadNoteTagsSQLiteUnit(t *testing.T) {
	db := openDatabase(t)
	defer closeDatabase(db, t)

	n := test.GetTestNotes()[0]
	saveNote(t, db, n)

	tags := n.Tags
	n.Tags = make([]*note.Tag, 0)

	if err := db.LoadNoteTags(n); err != nil {
		t.Fatal(err)
	} else {
		test.CheckTags(t, tags, n.Tags)
	}
}
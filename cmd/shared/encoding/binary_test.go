package encoding

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/anmil/quicknote"
	"github.com/anmil/quicknote/test"
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

var exportTestBytes = []byte{
	0x51, 0x4e, 0x4f, 0x54, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x58, 0xd7, 0x20, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x05, 0x00, 0x00, 0x00, 0x00, 0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00,
	0x00, 0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe3, 0x5a, 0x00, 0x00,
	0x00, 0x00, 0x58, 0xcf, 0xe3, 0x5a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x05, 0x62, 0x61, 0x73, 0x69, 0x63, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe3, 0x5a,
	0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe3, 0x5a, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe3,
	0x5a, 0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe3, 0x5a, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x06, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x02,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x5b, 0x00, 0x00, 0x00, 0x00,
	0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00, 0x58, 0xd7, 0x1a, 0xdf,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x62, 0x61, 0x73, 0x69,
	0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x54, 0x68, 0x69,
	0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x31, 0x20,
	0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x62, 0x61, 0x73, 0x69, 0x63,
	0x20, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xec, 0x23, 0x62, 0x61, 0x73, 0x69, 0x63, 0x20, 0x23, 0x74,
	0x65, 0x73, 0x74, 0x20, 0x23, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x0a,
	0x0a, 0x4c, 0x6f, 0x72, 0x65, 0x6d, 0x20, 0x69, 0x70, 0x73, 0x75, 0x6d,
	0x20, 0x64, 0x6f, 0x6c, 0x6f, 0x72, 0x20, 0x73, 0x69, 0x74, 0x20, 0x61,
	0x6d, 0x65, 0x74, 0x2c, 0x20, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x63, 0x74,
	0x65, 0x74, 0x75, 0x72, 0x20, 0x61, 0x64, 0x69, 0x70, 0x69, 0x73, 0x63,
	0x69, 0x6e, 0x67, 0x20, 0x65, 0x6c, 0x69, 0x74, 0x2e, 0x0a, 0x4e, 0x75,
	0x6c, 0x6c, 0x61, 0x20, 0x74, 0x69, 0x6e, 0x63, 0x69, 0x64, 0x75, 0x6e,
	0x74, 0x20, 0x64, 0x69, 0x61, 0x6d, 0x20, 0x65, 0x75, 0x20, 0x70, 0x75,
	0x72, 0x75, 0x73, 0x20, 0x6c, 0x61, 0x6f, 0x72, 0x65, 0x65, 0x74, 0x20,
	0x63, 0x6f, 0x6e, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x75, 0x6d, 0x2e,
	0x20, 0x44, 0x75, 0x69, 0x73, 0x0a, 0x74, 0x65, 0x6d, 0x70, 0x75, 0x73,
	0x2c, 0x20, 0x74, 0x75, 0x72, 0x70, 0x69, 0x73, 0x20, 0x76, 0x69, 0x74,
	0x61, 0x65, 0x20, 0x76, 0x61, 0x72, 0x69, 0x75, 0x73, 0x20, 0x75, 0x6c,
	0x6c, 0x61, 0x6d, 0x63, 0x6f, 0x72, 0x70, 0x65, 0x72, 0x2c, 0x20, 0x73,
	0x61, 0x70, 0x69, 0x65, 0x6e, 0x20, 0x65, 0x72, 0x61, 0x74, 0x0a, 0x63,
	0x75, 0x72, 0x73, 0x75, 0x73, 0x20, 0x6c, 0x61, 0x63, 0x75, 0x73, 0x2c,
	0x20, 0x65, 0x74, 0x20, 0x6c, 0x61, 0x63, 0x69, 0x6e, 0x69, 0x61, 0x20,
	0x6c, 0x69, 0x67, 0x75, 0x6c, 0x61, 0x20, 0x64, 0x6f, 0x6c, 0x6f, 0x72,
	0x20, 0x71, 0x75, 0x69, 0x73, 0x20, 0x6e, 0x69, 0x62, 0x68, 0x2e, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x03, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x5c,
	0x00, 0x00, 0x00, 0x00, 0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00,
	0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
	0x62, 0x61, 0x73, 0x69, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x25, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x23, 0x74, 0x65,
	0x73, 0x74, 0x20, 0x32, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x23, 0x62, 0x61, 0x73, 0x69, 0x63, 0x20, 0x23, 0x70, 0x61, 0x72, 0x73,
	0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xd6, 0x4c, 0x6f,
	0x72, 0x65, 0x6d, 0x20, 0x69, 0x70, 0x73, 0x75, 0x6d, 0x20, 0x64, 0x6f,
	0x6c, 0x6f, 0x72, 0x20, 0x73, 0x69, 0x74, 0x20, 0x61, 0x6d, 0x65, 0x74,
	0x2c, 0x20, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x63, 0x74, 0x65, 0x74, 0x75,
	0x72, 0x20, 0x61, 0x64, 0x69, 0x70, 0x69, 0x73, 0x63, 0x69, 0x6e, 0x67,
	0x20, 0x65, 0x6c, 0x69, 0x74, 0x2e, 0x0a, 0x4e, 0x75, 0x6c, 0x6c, 0x61,
	0x20, 0x74, 0x69, 0x6e, 0x63, 0x69, 0x64, 0x75, 0x6e, 0x74, 0x20, 0x64,
	0x69, 0x61, 0x6d, 0x20, 0x65, 0x75, 0x20, 0x70, 0x75, 0x72, 0x75, 0x73,
	0x20, 0x6c, 0x61, 0x6f, 0x72, 0x65, 0x65, 0x74, 0x20, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x75, 0x6d, 0x2e, 0x20, 0x44, 0x75,
	0x69, 0x73, 0x0a, 0x74, 0x65, 0x6d, 0x70, 0x75, 0x73, 0x2c, 0x20, 0x74,
	0x75, 0x72, 0x70, 0x69, 0x73, 0x20, 0x76, 0x69, 0x74, 0x61, 0x65, 0x20,
	0x76, 0x61, 0x72, 0x69, 0x75, 0x73, 0x20, 0x75, 0x6c, 0x6c, 0x61, 0x6d,
	0x63, 0x6f, 0x72, 0x70, 0x65, 0x72, 0x2c, 0x20, 0x73, 0x61, 0x70, 0x69,
	0x65, 0x6e, 0x20, 0x65, 0x72, 0x61, 0x74, 0x0a, 0x63, 0x75, 0x72, 0x73,
	0x75, 0x73, 0x20, 0x6c, 0x61, 0x63, 0x75, 0x73, 0x2c, 0x20, 0x65, 0x74,
	0x20, 0x6c, 0x61, 0x63, 0x69, 0x6e, 0x69, 0x61, 0x20, 0x6c, 0x69, 0x67,
	0x75, 0x6c, 0x61, 0x20, 0x64, 0x6f, 0x6c, 0x6f, 0x72, 0x20, 0x71, 0x75,
	0x69, 0x73, 0x20, 0x6e, 0x69, 0x62, 0x68, 0x2e, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
	0x00, 0x58, 0xcf, 0xe6, 0x5a, 0x00, 0x00, 0x00, 0x00, 0x58, 0xcf, 0xe6,
	0x5a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x71, 0x75, 0x69,
	0x73, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x5d, 0x00, 0x00,
	0x00, 0x00, 0x58, 0xd7, 0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00, 0x58, 0xd7,
	0x1a, 0xdf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x62, 0x61,
	0x73, 0x69, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x25, 0x54,
	0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x23, 0x74, 0x65, 0x73, 0x74,
	0x20, 0x32, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x23, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x20, 0x23, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xd8, 0x4c, 0x6f, 0x72, 0x65,
	0x6d, 0x20, 0x69, 0x70, 0x73, 0x75, 0x6d, 0x20, 0x64, 0x6f, 0x6c, 0x6f,
	0x72, 0x20, 0x73, 0x69, 0x74, 0x20, 0x61, 0x6d, 0x65, 0x74, 0x2c, 0x20,
	0x63, 0x6f, 0x6e, 0x73, 0x65, 0x63, 0x74, 0x65, 0x74, 0x75, 0x72, 0x20,
	0x61, 0x64, 0x69, 0x70, 0x69, 0x73, 0x63, 0x69, 0x6e, 0x67, 0x20, 0x65,
	0x6c, 0x69, 0x74, 0x2e, 0x0a, 0x4e, 0x75, 0x6c, 0x6c, 0x61, 0x20, 0x74,
	0x69, 0x6e, 0x63, 0x69, 0x64, 0x75, 0x6e, 0x74, 0x20, 0x64, 0x69, 0x61,
	0x6d, 0x20, 0x65, 0x75, 0x20, 0x70, 0x75, 0x72, 0x75, 0x73, 0x20, 0x6c,
	0x61, 0x6f, 0x72, 0x65, 0x65, 0x74, 0x20, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x6d, 0x65, 0x6e, 0x74, 0x75, 0x6d, 0x2e, 0x20, 0x44, 0x75, 0x69, 0x73,
	0x0a, 0x74, 0x65, 0x6d, 0x70, 0x75, 0x73, 0x2c, 0x20, 0x74, 0x75, 0x72,
	0x70, 0x69, 0x73, 0x20, 0x76, 0x69, 0x74, 0x61, 0x65, 0x20, 0x76, 0x61,
	0x72, 0x69, 0x75, 0x73, 0x20, 0x75, 0x6c, 0x6c, 0x61, 0x6d, 0x63, 0x6f,
	0x72, 0x70, 0x65, 0x72, 0x2c, 0x20, 0x73, 0x61, 0x70, 0x69, 0x65, 0x6e,
	0x20, 0x65, 0x72, 0x61, 0x74, 0x0a, 0x63, 0x75, 0x72, 0x73, 0x75, 0x73,
	0x20, 0x6c, 0x61, 0x63, 0x75, 0x73, 0x2c, 0x20, 0x65, 0x74, 0x20, 0x6c,
	0x61, 0x63, 0x69, 0x6e, 0x69, 0x61, 0x20, 0x6c, 0x69, 0x67, 0x75, 0x6c,
	0x61, 0x20, 0x64, 0x6f, 0x6c, 0x6f, 0x72, 0x20, 0x23, 0x71, 0x75, 0x69,
	0x73, 0x20, 0x6e, 0x69, 0x62, 0x68, 0x2e, 0x23, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
}

func TestBinaryEncoder(t *testing.T) {
	buff := &bytes.Buffer{}

	enc := NewBinaryEncoder(buff)
	_, err := enc.WriteHeader()
	if err != nil {
		t.Error(err)
		return
	}

	notes := test.GetTestNotesCust(notesJSON)
	for _, n := range notes {
		_, err = enc.WriteNote(n)
		if err != nil {
			t.Error(err)
			return
		}
	}

	encodedBytes := buff.Bytes()

	if len(encodedBytes) != len(exportTestBytes) {
		t.Errorf("Excepted len %d, got %d", len(encodedBytes), len(exportTestBytes))
	}
}

func TestBinaryDecoder(t *testing.T) {
	buff := bytes.NewReader(exportTestBytes)

	r := bufio.NewReader(buff)
	dec := NewBinaryDecoder(r)
	err := dec.ParseHeader()

	if err != nil {
		t.Error(err)
		return
	}

	if dec.Header.Version != 1 {
		t.Error("Expected Version 1, got", dec.Header.Version)
	}
	if dec.Header.Created.Unix() != 1490493535 {
		t.Error("Expected Timestamp 1490493535, got", dec.Header.Created.Unix())
	}

	notesChan, err := dec.ParseNotes()
	if err != nil {
		t.Error(err)
		return
	}

	notes := make(quicknote.Notes, 0)
	for n := range notesChan {
		notes = append(notes, n)
	}
	if err != nil {
		t.Error(err)
		return
	}

	tnotes := test.GetTestNotesCust(notesJSON)
	test.CheckNotes(t, notes, tnotes)
	for idx, n := range notes {
		test.CheckTags(t, n.Tags, tnotes[idx].Tags)
	}

	if notes[0].Book.Name != tnotes[0].Book.Name {
		t.Errorf("Expected book %s, got %s", tnotes[0].Book.Name, notes[0].Book.Name)
	}
}

func BenchmarkBinaryEncoder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buff := &bytes.Buffer{}

		enc := NewBinaryEncoder(buff)
		_, err := enc.WriteHeader()
		if err != nil {
			b.Error(err)
			return
		}

		notes := test.GetTestNotesCust(notesJSON)
		for _, n := range notes {
			_, err = enc.WriteNote(n)
			if err != nil {
				b.Error(err)
				return
			}
		}
	}
}

func BenchmarkBinaryDecoder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buff := bytes.NewReader(exportTestBytes)

		r := bufio.NewReader(buff)
		dec := NewBinaryDecoder(r)
		err := dec.ParseHeader()
		if err != nil {
			b.Error(err)
			return
		}

		notesChan, err := dec.ParseNotes()
		if err != nil {
			b.Error(err)
			return
		}

		notes := make(quicknote.Notes, 0)
		for n := range notesChan {
			notes = append(notes, n)
		}
		if dec.Err != nil {
			b.Error(err)
			return
		}
	}
}

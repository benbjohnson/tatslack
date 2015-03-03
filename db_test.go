package tatslack_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/benbjohnson/tatslack"
)

// Ensure we can save messages to the database for a channel.
func TestDB_SaveMessages(t *testing.T) {
	db := OpenDB()
	defer db.Close()

	// Write messages to database.
	a := []*tatslack.Message{
		{Type: "message", TS: "1358546515.000008", User: "U2147483896", Text: "hello"},
		{Type: "message", TS: "1358546515.000007", User: "U2147483896", Text: "goodbyte"},
	}
	if err := db.SaveMessages("C1234567890", a); err != nil {
		t.Fatal(err)
	}

	// Read messages back.
	messages, err := db.Messages("C1234567890")
	if err != nil {
		t.Fatal(err)
	} else if len(messages) != 2 {
		t.Fatalf("unexpected len: %d", len(messages))
	} else if !reflect.DeepEqual(messages[0], &tatslack.Message{Type: "message", TS: "1358546515.000007", User: "U2147483896", Text: "goodbyte"}) {
		t.Fatalf("unexpected message(0): %#v", messages[0])
	} else if !reflect.DeepEqual(messages[1], &tatslack.Message{Type: "message", TS: "1358546515.000008", User: "U2147483896", Text: "hello"}) {
		t.Fatalf("unexpected message(1): %#v", messages[1])
	}
}

func OpenDB() *tatslack.DB {
	db, err := tatslack.Open(tempfile())
	if err != nil {
		panic(err)
	}
	return db
}

func tempfile() string {
	f, _ := ioutil.TempFile("", "tatslack")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

package store_test

import (
	"context"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/tchaudhry91/algoprom/store"
)

var db *store.BoltStore

// This is not great, but this is how I'm debugging for now
func init() {
	store, err := store.NewBoltStore("/home/tchaudhry/Temp/algoprom.db", log.Default())
	if err != nil {
		panic("Could not open DB")
	}
	db = store
}

func TestGetCheck(t *testing.T) {
	output, err := db.GetCheck(context.Background(), "Sample HTTP Check - 2", "1739348709")
	if err != nil {
		t.Fatalf("Could not fetch Check:%v", err)
	}
	t.Log(output)
}

func TestGetBuckets(t *testing.T) {
	output, err := db.GetAllCheckNames(context.Background())
	if err != nil {
		t.Fatalf("Could not fetch Check Names:%v", err)
	}
	t.Log(output)
}

func TestGetNamedCheck(t *testing.T) {
	outputs, err := db.GetNamedCheck(context.Background(), "Sample HTTP Check", 100)
	if err != nil {
		t.Fatalf("Could not fetch named check")
	}
	t.Log(outputs)
}

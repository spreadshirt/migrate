package iofs

import (
	"testing"
	"testing/fstest"

	st "github.com/golang-migrate/migrate/v4/source/testing"
)

func Test(t *testing.T) {
	source := fstest.MapFS{
		"migrations/01_test.down.sql": &fstest.MapFile{},
		"migrations/01_test.up.sql":   &fstest.MapFile{},
		"migrations/03_test.up.sql":   &fstest.MapFile{},
		"migrations/04_test.down.sql": &fstest.MapFile{},
		"migrations/04_test.up.sql":   &fstest.MapFile{},
		"migrations/05_test.down.sql": &fstest.MapFile{},
		"migrations/07_test.down.sql": &fstest.MapFile{},
		"migrations/07_test.up.sql":   &fstest.MapFile{},
	}

	d, err := New(source, "migrations")
	if err != nil {
		t.Fatalf("New: %s", err.Error())
	}

	// For test details refer to:
	// https://github.com/golang-migrate/migrate/blob/8709f196bf6b98105c29ee7104a4ea85a49a8ed4/source/testing/testing.go
	st.Test(t, d)
}

package iofs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/source"
)

func init() {
	source.Register("iofs", &Driver{})
}

// Driver implements source.Driver for fs.ReadDirFS.
type Driver struct {
	source     fs.ReadDirFS
	migrations *source.Migrations
}

// New instantiates a migration driver from the given fs.ReadDirFS and looks for migrations below migrationRoot.
func New(migrationSource fs.ReadDirFS, migrationRoot string) (*Driver, error) {
	files, err := migrationSource.ReadDir(migrationRoot)
	if err != nil {
		return nil, err
	}

	migrations := source.NewMigrations()
	for _, file := range files {
		m, err := source.DefaultParse(filepath.Base(file.Name()))
		if err != nil {
			// ignore unparseable files
			continue
		}
		m.Raw = file.Name()

		ok := migrations.Append(m)
		if !ok {
			return nil, source.ErrDuplicateMigration{
				Migration: *m,
			}
		}
	}

	return &Driver{
		source:     migrationSource,
		migrations: migrations,
	}, nil
}

// Open implements source.Driver.
func (iofs *Driver) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("Open is not implemented, use New instead")
}

// Close implements source.Driver.
func (iofs *Driver) Close() error {
	// There is no Close defined for the underlying fs.FS.
	return nil
}

// First implements source.Driver.
func (iofs *Driver) First() (version uint, err error) {
	if version, ok := iofs.migrations.First(); ok {
		return version, nil
	}
	return 0, os.ErrNotExist
}

// Prev implements source.Driver.
func (iofs *Driver) Prev(version uint) (prevVersion uint, err error) {
	if version, ok := iofs.migrations.Prev(version); ok {
		return version, nil
	}
	return 0, os.ErrNotExist
}

// Next implements source.Driver.
func (iofs *Driver) Next(version uint) (nextVersion uint, err error) {
	if version, ok := iofs.migrations.Next(version); ok {
		return version, nil
	}
	return 0, os.ErrNotExist
}

// ReadUp implements source.Driver.
func (iofs *Driver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := iofs.migrations.Up(version); ok {
		f, err := iofs.source.Open(m.Raw)
		return f, m.Identifier, err
	}
	return nil, "", os.ErrNotExist
}

// ReadDown implements source.Driver.
func (iofs *Driver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := iofs.migrations.Down(version); ok {
		f, err := iofs.source.Open(m.Raw)
		return f, m.Identifier, err
	}
	return nil, "", os.ErrNotExist
}

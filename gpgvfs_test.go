package gpgvfs

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/psanford/sqlite3vfs"
)

func setupEncryptedFile(content []byte) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "gpgvfs")
	if err != nil {
		panic(err)
	}

	dbPath := fmt.Sprintf("%s/%s", tmpDir, "test.db")
	os.WriteFile(dbPath, content, fs.ModePerm)

	return dbPath, func() {
		os.RemoveAll(tmpDir)
	}
}

func TestReadEncryptedFileSuccess(t *testing.T) {
	dbPath, cleanup := setupEncryptedFile(TEST_FILE_ENCRYPTED)
	defer cleanup()

	contents, err := readEncryptedFile(dbPath, PASSWORD)
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != TEST_FILE_DECRIPTED {
		t.Fatalf("Expected contents to be %s, got \"%s\"", TEST_FILE_DECRIPTED, string(contents))
	}
}

func TestReadEncryptedFileError(t *testing.T) {
	dbPath, cleanup := setupEncryptedFile(TEST_FILE_ENCRYPTED)
	defer cleanup()

	defer func() {
		if r := recover(); r == nil || r != "openpgp.ReadMessage: wrong password" {
			t.Fatalf("Expected panic, got %s", r)
		}
	}()

	readEncryptedFile(dbPath, []byte("wrong-password"))
}

func TestNewGPGVFS(t *testing.T) {
	dbPath, cleanup := setupEncryptedFile(TEST_DB_ENCRYPTED)
	defer cleanup()

	// create vfs
	vfs := NewGPGVFS(dbPath, PASSWORD)
	err := sqlite3vfs.RegisterVFS("gpgvfs", vfs)
	if err != nil {
		t.Fatal(err)
	}
	defer vfs.Close(dbPath, PASSWORD)

	// create connection
	db, err := sql.Open("sqlite3", "/data/test.db?vfs=gpgvfs")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// query
	rows, err := db.Query("SELECT id, name FROM test ORDER BY id ASC")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	type rowType struct {
		Id   int
		Name string
	}

	result := []rowType{}
	for rows.Next() {
		var row rowType
		err = rows.Scan(&row.Id, &row.Name)
		result = append(result, row)
	}

	if len(result) == 0 {
		t.Fatalf("Expected two rows, got %d", len(result))
	}

	if result[0].Id != 1 || result[0].Name != "test1" {
		t.Fatalf("Expected result[0]={1 test1}, got %v", result[0])
	}

	if result[1].Id != 2 || result[1].Name != "test2" {
		t.Fatalf("Expected result[1]={2 test2}, got %v", result[1])
	}
}

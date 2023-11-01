package gpgvfs

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/psanford/sqlite3vfs"
	"golang.org/x/crypto/openpgp"
)

type GPGVFS struct {
	EncryptedPath string
	memfs         *memfs.MemFS
}

func NewGPGVFS(EncryptedPath string, EncryptionPassphrase []byte) *GPGVFS {
	memfs := memfs.Create()
	memfs.Mkdir("/data", 0777)

	fileContents, err := readEncryptedFile(EncryptedPath, EncryptionPassphrase)

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(fileContents)))
	_, err = base64.StdEncoding.Decode(decoded, fileContents)
	if err != nil {
		panic(err)
	}

	// copy contents
	f, err := memfs.OpenFile("/data/test.db", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	if n, err := f.Write([]byte(decoded)); err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
	} else if n != len(fileContents) {
		fmt.Printf("Invalid write count: %d\n", n)
	}
	f.Close()

	return &GPGVFS{
		EncryptedPath: "/data/test.db",
		memfs:         memfs,
	}
}

func readEncryptedFile(EncryptedPath string, EncryptionPassphrase []byte) ([]byte, error) {
	dbFile, err := os.Open(EncryptedPath)

	if err != nil {
		panic(err)
	}

	attempts := 0

	md, err := openpgp.ReadMessage(dbFile, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if attempts > 1 {
			panic("openpgp.ReadMessage: wrong password")
		}

		attempts = attempts + 1
		return EncryptionPassphrase, nil
	}, nil)

	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(md.UnverifiedBody)

	return bytes, nil
}

func (gpgvfs *GPGVFS) Open(name string, flags sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	var (
		f   vfs.File
		err error
	)

	f, err = gpgvfs.memfs.OpenFile(name, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	tf := &gpgFile{Name: name, f: f}
	return tf, flags, nil
}

func (gpgvfs *GPGVFS) Close(name string, EncryptionPassphrase []byte) error {
	fileinfo, err := gpgvfs.memfs.Stat("/data/test.db")

	size := fileinfo.Size()
	p := make([]byte, size)

	f, err := gpgvfs.memfs.OpenFile("/data/test.db", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	f.ReadAt(p, 0)

	data := make([]byte, base64.StdEncoding.EncodedLen(len(p)))
	base64.StdEncoding.Encode(data, p)

	dst, _ := os.Create(name)
	encryptor, _ := openpgp.SymmetricallyEncrypt(dst, EncryptionPassphrase, nil, nil)
	encryptor.Write(data)
	encryptor.Close()
	dst.Close()

	return nil
}

func (gpgvfs *GPGVFS) Delete(name string, dirSync bool) error {
	gpgvfs.memfs.Remove(name)
	return nil
}

func (gpgvfs *GPGVFS) Access(name string, flag sqlite3vfs.AccessFlag) (bool, error) {
	_, err := gpgvfs.memfs.Stat(name)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (gpgvfs *GPGVFS) FullPathname(name string) string {
	return name
}

package gpgvfs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/blang/vfs/memfs"
	"github.com/psanford/sqlite3vfs"
	"golang.org/x/crypto/openpgp"
)

type GPGVFS struct {
	EncryptedPath string
	memfs         *memfs.MemFS
}

var ErrWrongPassword = errors.New("openpgp.ReadMessage: wrong password")

func NewGPGVFS(EncryptedPath string, EncryptionPassphrase []byte) (*GPGVFS, error) {
	memfs := memfs.Create()
	memfs.Mkdir("/data", 0777)

	fileContents, err := readEncryptedFile(EncryptedPath, EncryptionPassphrase)
	if err != nil {
		return nil, err
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(fileContents)))
	_, err = base64.StdEncoding.Decode(decoded, fileContents)
	if err != nil {
		return nil, err
	}

	// copy contents
	f, err := memfs.OpenFile("/data/test.db", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return nil, err
	}

	if n, err := f.Write([]byte(decoded)); err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
	} else if n != len(decoded) {
		fmt.Printf("Invalid write count: %d\n", n)
	}

	f.Close()

	return &GPGVFS{
		EncryptedPath: "/data/test.db",
		memfs:         memfs,
	}, nil
}

func readEncryptedFile(EncryptedPath string, EncryptionPassphrase []byte) ([]byte, error) {
	dbFile, err := os.Open(EncryptedPath)

	if err != nil {
		return nil, err
	}

	attempts := 0

	md, err := openpgp.ReadMessage(dbFile, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if attempts > 1 {
			return nil, ErrWrongPassword
		}

		attempts = attempts + 1
		return EncryptionPassphrase, nil
	}, nil)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(md.UnverifiedBody)
}

func (gpgvfs *GPGVFS) Open(name string, flags sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	f, err := gpgvfs.memfs.OpenFile(name, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, 0, err
	}

	return &gpgFile{
		Name: name,
		f:    f,
	}, flags, nil
}

func (gpgvfs *GPGVFS) Close(name string, EncryptionPassphrase []byte) error {
	fileinfo, err := gpgvfs.memfs.Stat("/data/test.db")
	if err != nil {
		return err
	}

	f, err := gpgvfs.memfs.OpenFile("/data/test.db", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	size := fileinfo.Size()
	p := make([]byte, size)
	f.ReadAt(p, 0)

	data := make([]byte, base64.StdEncoding.EncodedLen(len(p)))
	base64.StdEncoding.Encode(data, p)

	dst, err := os.Create(name)
	if err != nil {
		return err
	}

	encryptor, err := openpgp.SymmetricallyEncrypt(dst, EncryptionPassphrase, nil, nil)
	if err != nil {
		return err
	}

	_, err = encryptor.Write(data)
	if err != nil {
		return err
	}

	err = encryptor.Close()
	if err != nil {
		return err
	}

	err = dst.Close()
	if err != nil {
		return err
	}

	return nil
}

func (gpgvfs *GPGVFS) Delete(name string, dirSync bool) error {
	return gpgvfs.memfs.Remove(name)
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

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

const DecryptedPath = "/data/file.db"

type GPGVFS struct {
	memfs *memfs.MemFS
}

var ErrWrongPassword = errors.New("openpgp.ReadMessage: wrong password")

func InitializeFile(path string, password []byte) error {
	vfs, err := newGPGVFSFromContent([]byte(sqliteEmptyFile))
	if err != nil {
		return err
	}

	return vfs.Close(path, password)
}

func NewGPGVFS(path string, password []byte) (*GPGVFS, error) {
	fileContents, err := readEncryptedFile(path, password)
	if err != nil {
		return nil, err
	}

	return newGPGVFSFromContent(fileContents)
}

func newGPGVFSFromContent(contentB64Encoded []byte) (*GPGVFS, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(contentB64Encoded)))
	_, err := base64.StdEncoding.Decode(decoded, contentB64Encoded)
	if err != nil {
		return nil, err
	}

	memfs := memfs.Create()
	memfs.Mkdir("/data", 0777)
	f, err := memfs.OpenFile(DecryptedPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return nil, err
	}

	n, err := f.Write([]byte(decoded))
	if err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
		// @TODO return err
	} else if n != len(decoded) {
		fmt.Printf("Invalid write count: %d\n", n)
		// @TODO return err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return &GPGVFS{memfs: memfs}, nil
}

func readEncryptedFile(path string, password []byte) ([]byte, error) {
	dbFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	attempts := 0

	md, err := openpgp.ReadMessage(dbFile, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if attempts > 1 {
			return nil, ErrWrongPassword
		}

		attempts = attempts + 1
		return password, nil
	}, nil)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(md.UnverifiedBody)
}

func (gpgvfs *GPGVFS) Open(path string, flags sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	f, err := gpgvfs.memfs.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, 0, err
	}

	return &gpgFile{f: f}, flags, nil
}

func (gpgvfs *GPGVFS) Close(path string, EncryptionPassphrase []byte) error {
	fileinfo, err := gpgvfs.memfs.Stat(DecryptedPath)
	if err != nil {
		return err
	}

	f, err := gpgvfs.memfs.OpenFile(DecryptedPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	size := fileinfo.Size()
	p := make([]byte, size)
	_, err = f.ReadAt(p, 0)
	if err != nil {
		panic(err)
	}

	data := make([]byte, base64.StdEncoding.EncodedLen(len(p)))
	base64.StdEncoding.Encode(data, p)

	dst, err := os.Create(path)
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

	return dst.Close()
}

func (gpgvfs *GPGVFS) Delete(path string, dirSync bool) error {
	return gpgvfs.memfs.Remove(path)
}

func (gpgvfs *GPGVFS) Access(path string, flag sqlite3vfs.AccessFlag) (bool, error) {
	_, err := gpgvfs.memfs.Stat(path)

	if err != nil {
		// @TODO return err
		return false, nil
	}

	return true, nil
}

func (gpgvfs *GPGVFS) FullPathname(path string) string {
	return path
}

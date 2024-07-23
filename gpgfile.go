package gpgvfs

import (
	"io"
	"sync/atomic"

	"github.com/blang/vfs"
	"github.com/psanford/sqlite3vfs"
)

type gpgFile struct {
	lockCount int64
	f         vfs.File
}

func (tf *gpgFile) Close() error {
	return tf.f.Close()
}

func (tf *gpgFile) ReadAt(p []byte, off int64) (int, error) {
	return tf.f.ReadAt(p, off)
}

func (tf *gpgFile) WriteAt(b []byte, off int64) (int, error) {
	_, err := tf.f.Seek(off, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return tf.f.Write(b)
}

func (tf *gpgFile) Truncate(size int64) error {
	return nil
}

func (tf *gpgFile) Sync(flag sqlite3vfs.SyncType) error {
	return tf.f.Sync()
}

func (tf *gpgFile) FileSize() (int64, error) {
	cur, err := tf.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	end, err := tf.f.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = tf.f.Seek(cur, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return end, nil
}

func (tf *gpgFile) Lock(elock sqlite3vfs.LockType) error {
	if elock == sqlite3vfs.LockNone {
		return nil
	}

	atomic.AddInt64(&tf.lockCount, 1)
	return nil
}

func (tf *gpgFile) Unlock(elock sqlite3vfs.LockType) error {
	if elock == sqlite3vfs.LockNone {
		return nil
	}

	atomic.AddInt64(&tf.lockCount, -1)
	return nil
}

func (tf *gpgFile) CheckReservedLock() (bool, error) {
	count := atomic.LoadInt64(&tf.lockCount)
	return count > 0, nil
}

func (tf *gpgFile) SectorSize() int64 {
	return 0
}

func (tf *gpgFile) DeviceCharacteristics() sqlite3vfs.DeviceCharacteristic {
	return 0
}

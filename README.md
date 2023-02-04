# sqlite3-gpgvfs

Combine [github.com/blang/vfs](https://github.com/blang/vfs) and [github.com/psanford/sqlite3vfs](https://github.com/psanford/sqlite3vfs) to directly open GPG-encrypted SQLite3 databases.

Compatible with [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3).

## TODO

- sqlite file has to be wrapped in base64, otherwise it misses some bytes and `sqlite3 decrypted.db` says the db is corrupted;

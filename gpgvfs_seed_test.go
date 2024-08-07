package gpgvfs

var PASSWORD = []byte("admin123")

// echo -n 'sqlite3-gpgvfs-abc1234d' | gpg --pinentry-mode loopback --passphrase "admin123" --symmetric - | od -v -t x1 -An | xargs -n1 | awk '{printf "0x" $1 ","}'
var TEST_FILE_DECRIPTED = "sqlite3-gpgvfs-abc1234d"
var TEST_FILE_ENCRYPTED = []byte{
	0x8c, 0x0d, 0x04, 0x09, 0x03, 0x02, 0x89, 0x16, 0xbc, 0xeb, 0xc6, 0x22, 0x1b, 0xf8, 0xf9, 0xd2,
	0x4c, 0x01, 0xaa, 0x1c, 0xe5, 0x37, 0x3f, 0x31, 0xc0, 0x6a, 0x5e, 0x06, 0x92, 0x30, 0x9d, 0x09,
	0x02, 0x28, 0x9c, 0x8a, 0x32, 0x2a, 0x80, 0x84, 0xf4, 0x12, 0x9a, 0xac, 0x1d, 0xf6, 0xe8, 0x8a,
	0xfd, 0x0e, 0x73, 0xb5, 0x79, 0xf6, 0xd1, 0x01, 0xdb, 0xc2, 0x3e, 0x47, 0x36, 0xe5, 0x0b, 0x5d,
	0xd9, 0x39, 0xbf, 0x34, 0xc8, 0x4b, 0x74, 0xbf, 0xec, 0x25, 0x9f, 0xcb, 0x1b, 0x9b, 0xf7, 0x74,
	0xbb, 0xfb, 0x78, 0x54, 0x04, 0xd0, 0x0e, 0x5e, 0x03, 0x1e, 0x3e, 0xdb, 0x52,
}

/*
sqlite3 test.db <<EOF
CREATE TABLE test (id INTEGER PRIMARY KEY, name VARCHAR(255) NOT NULL);
INSERT INTO test (name) VALUES ('test1'), ('test2');
EOF

base64 test.db | gpg --pinentry-mode loopback --passphrase "admin123" --symmetric - | od -v -t x1 -An | xargs -n1 | awk '{printf "0x" $1 ","}'
*/
var TEST_DB_ENCRYPTED = []byte{
	0x8c, 0x0d, 0x04, 0x09, 0x03, 0x08, 0x14, 0xe1,
	0x8a, 0xf8, 0x36, 0x5b, 0x26, 0xdd, 0xff, 0xd2,
	0xc0, 0xa6, 0x01, 0x0b, 0x1e, 0xcc, 0xce, 0x10,
	0x0f, 0x28, 0xb9, 0x87, 0xf0, 0xf4, 0x7c, 0xbf,
	0xb2, 0x53, 0xaa, 0xed, 0x4a, 0x68, 0xd0, 0x98,
	0x09, 0x98, 0x16, 0x9a, 0x8f, 0x13, 0x58, 0x5e,
	0xf8, 0x9c, 0x95, 0xb5, 0xa8, 0xe5, 0x01, 0xbe,
	0xb6, 0xa3, 0x85, 0xad, 0x55, 0x47, 0xbd, 0x7c,
	0x33, 0xe1, 0xde, 0x49, 0x9b, 0x59, 0xe1, 0x24,
	0xf5, 0xf2, 0xc4, 0x8c, 0xad, 0xa8, 0xde, 0x68,
	0xe8, 0x69, 0xb8, 0x13, 0xf6, 0x06, 0x98, 0x5a,
	0x92, 0x08, 0x1b, 0x07, 0x1e, 0xe2, 0xaf, 0x35,
	0x49, 0x74, 0x9f, 0x7d, 0x7a, 0xbc, 0x10, 0x4f,
	0xd2, 0x38, 0xe8, 0x1e, 0x4c, 0x36, 0xa1, 0x9e,
	0xeb, 0x6a, 0x3a, 0x07, 0xae, 0x23, 0x0c, 0x77,
	0x7c, 0xbf, 0xbb, 0x94, 0x5a, 0x7c, 0x9e, 0x0f,
	0xdd, 0x2d, 0xbe, 0x87, 0x89, 0x2f, 0xc4, 0x55,
	0x75, 0xb8, 0x20, 0x52, 0x67, 0x4a, 0x2e, 0xe3,
	0x2b, 0x95, 0x5e, 0x80, 0x64, 0x23, 0xbd, 0x48,
	0x49, 0xb5, 0x95, 0xcb, 0x87, 0x68, 0x36, 0x59,
	0xb2, 0x20, 0x94, 0xd7, 0x16, 0xb6, 0x91, 0x53,
	0x93, 0x0a, 0x2f, 0xdb, 0xad, 0x26, 0x95, 0x04,
	0x15, 0x0c, 0xdd, 0x9b, 0x4f, 0x28, 0x93, 0x6c,
	0x02, 0xba, 0x3f, 0x19, 0x2c, 0x9d, 0xe4, 0x54,
	0x41, 0x45, 0x9a, 0x4a, 0x2b, 0xa0, 0x84, 0x87,
	0xc9, 0xad, 0xdf, 0x18, 0xde, 0xc2, 0xc0, 0x2d,
	0xef, 0x58, 0x29, 0x06, 0xce, 0xa1, 0xd0, 0x1a,
	0xb0, 0x46, 0x87, 0xf6, 0xb7, 0xc0, 0x74, 0x04,
	0x30, 0x87, 0x19, 0xf8, 0xdb, 0x96, 0x29, 0x78,
	0xbd, 0x55, 0x9d, 0xfc, 0x4d, 0x13, 0x3b, 0xdf,
	0xd1, 0x4d, 0xda, 0x46, 0x35, 0xaf, 0x2e, 0x85,
	0x42, 0x08, 0x4a, 0xf0, 0x30, 0xb7, 0xef, 0x36,
	0xaf, 0xa9, 0x73, 0x28, 0x0c, 0x02, 0x83, 0x4c,
	0x4f, 0x12, 0x81, 0x1c, 0x32, 0xf5, 0x43, 0x7d,
	0x5d, 0x4a, 0x5e, 0x91, 0x36, 0xc4, 0x71, 0x7f,
	0xa1, 0x1f, 0x16, 0x4b, 0x4f, 0xbc, 0x08, 0xef,
	0x60, 0x0d, 0xc0, 0x78, 0x35, 0x5b, 0x6f, 0x1c,
	0x95, 0xd0, 0x27, 0xa0, 0xd4, 0x84, 0x8e, 0xa6,
	0xd5, 0xbb, 0x71, 0x3e, 0x7e, 0xc6, 0x5d, 0x01,
	0x64, 0x38, 0xe4, 0xd0, 0xf9, 0x24, 0x3d, 0xde,
	0xa6, 0x16, 0x9f, 0x4a, 0x61, 0xa0, 0x29, 0xf0,
	0x85, 0xd0, 0x35, 0xe8, 0xc2, 0x01, 0xee, 0xb0,
	0xff, 0x69, 0xd7, 0xc6, 0x02, 0x32, 0x63, 0xbb,
	0x1d, 0xa8, 0x2a, 0x4d, 0x51, 0x62, 0xe8, 0xdb,
	0x0f, 0x34, 0x3a, 0xe0, 0xb1, 0x9d, 0x44, 0xa4,
	0x26, 0x23, 0x7c, 0x6c, 0xba, 0xb9, 0xc9, 0x53,
	0x04, 0xc6, 0x3e, 0x92, 0x66, 0xc2, 0x50, 0x24,
}

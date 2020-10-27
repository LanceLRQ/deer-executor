package judge_result

/*********
------------------------
|MAG|VER|CMP|RSZ|BSZ|CSZ| Certificate |SSZ| Signature | Result | Body
------------------------
| 2 | 1 | 1 | 4 | 4 | 2 | ... | 2 | ...
------------------------
**********/
type JudgeResultPackage struct {
	Version 		uint8				// (VER) Package Version
	CompressorType	uint8				// (CMP) Compressor type: 0-disabled; 1-gzip
	ResultSize 		uint32				// (RSZ) Result JSON Text Size
	BodySize 		uint32				// (BSZ) Result Body Size
	CertSize		uint16				// (CSZ) Public Certificate Size
	SignSize		uint16				// (SSZ) Signature Size
	Certificate		[]byte				// Public Certificate
	Signature     	[]byte		 		// Signature: SHA256(Result + Body)
	Result 			[]byte				// Result JSON
	BodyPackageFile	string				// Body package file
}

/***
------------------------
Magic | Size | FileName | Content
------------------------
  2   |  4  | (Sep: \n) |  ...
------------------------
***/
type JudgeResultPackageBody struct {
	BodyPackageFile		string
	Files 				[]struct {
		Size     		uint32
		FileName 		string
		Position 		uint32
	}
}

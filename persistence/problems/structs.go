package judge_result

import "github.com/LanceLRQ/deer-executor/persistence"

/*********
------------------------
|MAG|VER|CSZ|BSZ|PCSZ| Certificate |SSZ| Signature | Result | Body
------------------------
| 2 | 2 | 4 | 4 | 2 | ... | 2 | ...
------------------------
**********/
type ProblemPackage struct {
	Version 		uint16				// (VER) Package Version
	ConfigSize 		uint32				// (CSZ) Config JSON Text Size
	BodySize 		uint32				// (BSZ) Result Body Size
	CertSize		uint16				// (PCSZ) Public Certificate Size
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
type ProblemPackageBody struct {
	BodyPackageFile		string
	Files 				[]struct {
		Size     		uint32
		FileName 		string
		Position 		uint32
	}
}

type ProblemPersisOptions struct {
	DigitalSign    bool
	DigitalPEM     persistence.DigitalSignPEM
	OutFile        string
}

package problems

import "github.com/LanceLRQ/deer-executor/persistence"

/*********
------------------------
|MAG|VER|CMT|CSZ|BSZ|PCSZ| Certificate |SSZ| Signature | Result | Body
------------------------
| 2 | 2 | 4 | 4 | 4 | 2 | ... | 2 | ...
------------------------
**********/
type ProblemPackage struct {
	Version 		uint16				// (VER) Package Version
	CommitVersion 	uint32				// (CMT) Commit Version
	ConfigSize 		uint32				// (CSZ) Config JSON Text Size
	BodySize 		uint32				// (BSZ) Result Body Size
	CertSize		uint16				// (PCSZ) Public Certificate Size
	SignSize		uint16				// (SSZ) Signature Size
	Certificate		[]byte				// Public Certificate
	Signature     	[]byte		 		// Signature: SHA256(Result + Body)
	Configs 		[]byte				// Configs JSON
	BodyPackageFile	string				// Body package file
}

type ProblemPersisOptions struct {
	DigitalSign    bool
	DigitalPEM     persistence.DigitalSignPEM
	OutFile        string
}

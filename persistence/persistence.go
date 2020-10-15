package persistence


/*********
------------------------
|MAG|VER|RSZ|BSZ|SSZ|
------------------------
| 2 | 1 | 4 | 4 | 2 |
------------------------
**********/
type JudgeResultPackage struct {
	Version 		uint8				// (VER) Package Version
	ResultSize 		uint32				// (RSZ) Result JSON Text Size
	BodySize 		uint32				// (BSZ) Result Body Size
	SignSize		uint16				// (SSZ) Signature Size
	CertSize		uint16				// (CSZ) Public Certificate Size
	Certificate		[]byte				// Public Certificate
	Signature       []byte		 		// Signature: SHA256(Result + Body)
	Result 			[]byte				// Result JSON
	Body			[]byte				// Body Binary
}

type JudgeResultPackageBody struct {

}
//
//func persistentJudgeResult(
//	session *executor.JudgeSession,
//	rst *executor.JudgeResult,
//	certKeyFile string,
//	outFile string,
//) error {
//
//}
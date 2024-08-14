package models

type Secrets struct {
	Namespace  string
	Secret     string
	Set        bool
	Key        string
	Value      string
	Get        bool
	Delete     bool
	List       bool
	All        bool
	Operation  string
	EnvPath    string
	FillPath   string
	Modify     bool
	FileFormat string
}

var SupportedFormats = map[string]struct{}{
	"yaml": {},
	"json": {},
}

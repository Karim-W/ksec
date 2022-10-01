package models

type Secrets struct {
	Namespace string
	Secret    string
	Set       bool
	Key       string
	Value     string
	Get       bool
	Delete    bool
	List      bool
	All       bool
	Operation string
}

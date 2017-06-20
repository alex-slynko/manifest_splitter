package types

type Operation struct {
	Type  string
	Path  string
	Value interface{}
}

//[]Operation -> yaml.Marshal() -> print

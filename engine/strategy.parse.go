package engine

import (
	"encoding/json"
	"fmt"
)

type ArgTyp int

func (t ArgTyp) String() string {
	switch t {
	case 0:
		return "int"
	case 1:
		return "double"
	case 2:
		return "string"
	case 3:
		return "enum"
	}
	return ""
}

type Arg struct {
	Typ   ArgTyp
	Value interface{}
}

func (a *Arg) UnmarshalJSON(data []byte) error {
	var aa struct {
		Typ   ArgTyp
		Value interface{}
	}
	if err := json.Unmarshal(data, &aa); err != nil {
		return err
	}
	switch aa.Typ {
	case 0, 3:
		a.Value = int((aa.Value).(float64))
	default:
		a.Value = aa.Value
	}
	a.Typ = aa.Typ
	return nil
}

func ArgsTo(args []Arg) (r []interface{}) {
	for _, a := range args {
		r = append(r, a.Value)
	}
	return
}

type Rule struct {
	Path      []int
	Factor    string
	Operation int
	Args      []Arg
}

func (r Rule) String() string{
	return fmt.Sprintf("domain_path = %v, factor = %s, op_id = %d, args = %v; ", r.Path, r.Factor, r.Operation, r.Args)
}

type Strategy struct {
	EventID     int
	Rules       []Rule
	Combination string
}

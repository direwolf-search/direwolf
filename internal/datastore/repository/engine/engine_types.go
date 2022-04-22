package engine

type EType int

const (
	SQL EType = iota
	NoSQL
)

func (e EType) Int() int {
	return int(e)
}

func (e EType) String() string {
	return eTypes[e.Int()]
}

var eTypes = []string{
	"SQL",
	"NoSQL",
}

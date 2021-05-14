package sdk

// Obj is an entity in the state
type Obj struct {
	Data map[string]string
}

// StateResolver is an interface used to resolve state information
type StateResolver interface {
	GetObj2(table string, keys map[string]string) (*Obj, error)
	GetObjs2(q *Query) ([]*Obj, error)
}

const (
	AscOrder  = "asc"
	DescOrder = "desc"
)

type WhereCond string

const (
	WhereCondEqual = "equal"
)

type QueryWhere struct {
	Key   string
	Val   string
	Where WhereCond
}

type Query struct {
	Table   string
	First   uint64
	Skip    uint64
	OrderBy string
	Order   string
	Where   []QueryWhere
}

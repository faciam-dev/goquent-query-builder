package structs

type Column struct {
	Name   string
	Raw    string
	Values []interface{}
}

type Table struct {
	Name string
}

type Where struct {
	Column    string
	Condition string
	Value     []interface{}
	Operator  int
	Query     *Query
}

type WhereGroup struct {
	Conditions   []Where
	Subgroups    []WhereGroup
	Operator     int
	IsDummyGroup bool
}

type Query struct {
	Columns         *[]Column
	Table           Table
	Joins           *[]Join
	ConditionGroups *[]WhereGroup
	Conditions      *[]Where
	Limit           *Limit
	Offset          *Offset
	Order           *[]Order
	SubQuery        *[]Query
	Group           *GroupBy
	Lock            *Lock
}

type SelectQuery struct {
	Table    string
	Columns  *[]Column
	Limit    *Limit
	Offset   *Offset
	SubQuery *[]Query
	Group    *GroupBy
	Lock     *Lock
}

type InsertQuery struct {
	Table       string
	Values      map[string]interface{}
	ValuesBatch []map[string]interface{}
	Columns     []string
	Query       *Query
}

type UpdateQuery struct {
	Table  string
	Values map[string]interface{}
	Query  *Query
}

type DeleteQuery struct {
	Table string
	Query *Query
}

type Join struct {
	Name               string
	TargetNameMap      map[string]string
	SearchColumn       string
	SearchCondition    string
	SearchTargetColumn string
}

type Limit struct {
	Limit int64
}

type Offset struct {
	Offset int64
}

type Order struct {
	Column string
	IsAsc  bool
	Raw    string
}

type Orders struct {
	Orders []*Order
}

type GroupBy struct {
	Columns []string
	Having  *[]Having
}

type Having struct {
	Column    string
	Condition string
	Value     interface{}
	Operator  int
	Raw       string
}

type Lock struct {
	LockType string
}

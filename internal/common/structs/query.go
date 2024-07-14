package structs

type Column struct {
	Name   string
	Raw    string
	Values []string
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
	Limit           Limit
	Order           *[]Order
	SubQuery        *[]Query
	Group           *GroupBy
}

type Join struct {
	Name               string
	TargetNameMap      map[string]string
	SearchColumn       string
	SearchCondition    string
	SearchTargetColumn string
}

type Limit struct {
	Offset   int64
	RowCount int64
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

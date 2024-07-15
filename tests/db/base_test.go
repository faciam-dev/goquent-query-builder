package db_test

import (
	"strings"
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type QueryBuilderExpected struct {
	Expected string
	Values   []interface{}
}

func TestBaseQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		input    structs.Query
		expected QueryBuilderExpected
	}{
		{
			"Select",
			"Select",
			structs.Query{
				Columns: &[]structs.Column{
					{Name: "id"},
					{Name: "name"},
				},
			},
			QueryBuilderExpected{
				Expected: "SELECT id, name",
				Values:   nil,
			},
		},
		{
			"From",
			"From",
			structs.Query{
				Table: structs.Table{Name: "users"},
			},
			QueryBuilderExpected{
				Expected: "FROM users",
				Values:   nil,
			}},

		{
			"Where",
			"Where",
			structs.Query{
				ConditionGroups: &[]structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:    "age",
								Condition: ">",
								Value:     []interface{}{18},
								Operator:  consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
					},
				},
			},
			QueryBuilderExpected{
				Expected: " WHERE age > ?",
				Values:   []interface{}{18},
			},
		},
		{
			"WhereQuery",
			"Where",
			structs.Query{
				ConditionGroups: &[]structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:    "age",
								Condition: ">",
								Query: &structs.Query{
									Columns: &[]structs.Column{
										{Name: "id"},
									},
									Table:           structs.Table{Name: "users"},
									ConditionGroups: &[]structs.WhereGroup{},
									Conditions:      &[]structs.Where{},
									Joins:           &[]structs.Join{},
									Order:           &[]structs.Order{},
									Group:           &structs.GroupBy{},
								},
								Operator: consts.LogicalOperator_AND,
							},
							{
								Column:    "name",
								Condition: "=",
								Value:     []interface{}{"John"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
					},
					{
						Conditions: []structs.Where{
							{
								Column:    "city",
								Condition: "=",
								Value:     []interface{}{"New York"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
					},
				},
			},
			QueryBuilderExpected{
				Expected: " WHERE (age > (SELECT id FROM users) AND name = ?) AND city = ?",
				Values:   []interface{}{"John", "New York"},
			},
		},
		{
			"WhereGroup_Or",
			"WhereGroup",
			structs.Query{
				ConditionGroups: &[]structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:    "age",
								Condition: ">",
								Value:     []interface{}{18},
								Operator:  consts.LogicalOperator_AND,
							},
							{
								Column:    "name",
								Condition: "=",
								Value:     []interface{}{"John"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
						Operator: consts.LogicalOperator_AND,
					},
					{
						Conditions: []structs.Where{
							{
								Column:    "age",
								Condition: ">",
								Value:     []interface{}{18},
								Operator:  consts.LogicalOperator_AND,
							},
							{
								Column:    "name",
								Condition: "=",
								Value:     []interface{}{"John"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
						Operator: consts.LogicalOperator_OR,
					},
				},
			},
			QueryBuilderExpected{
				Expected: " WHERE (age > ? AND name = ?) OR (age > ? AND name = ?)",
				Values:   []interface{}{18, "John", 18, "John"},
			},
		},
		{
			"WhereGroup",
			"WhereGroup",
			structs.Query{
				ConditionGroups: &[]structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:    "age",
								Condition: ">",
								Value:     []interface{}{18},
								Operator:  consts.LogicalOperator_AND,
							},
							{
								Column:    "name",
								Condition: "=",
								Value:     []interface{}{"John"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
					},
					{
						Conditions: []structs.Where{
							{
								Column:    "city",
								Condition: "=",
								Value:     []interface{}{"New York"},
								Operator:  consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
					},
				},
			},
			QueryBuilderExpected{
				Expected: " WHERE (age > ? AND name = ?) AND city = ?",
				Values:   []interface{}{18, "John", "New York"},
			},
		},
		{
			"Join",
			"Join",
			structs.Query{
				Table: structs.Table{Name: "users"},
				Joins: &[]structs.Join{
					{
						Name:               "orders",
						TargetNameMap:      map[string]string{"inner": "orders"},
						SearchColumn:       "users.id",
						SearchCondition:    "=",
						SearchTargetColumn: "orders.user_id",
					},
				},
			},
			QueryBuilderExpected{
				Expected: " INNER JOIN orders ON users.id = orders.user_id",
				Values:   nil,
			},
		},
		{
			"OrderBy",
			"OrderBy",
			structs.Query{
				Order: &[]structs.Order{
					{
						Column: "name",
						IsAsc:  true,
					},
				},
			},
			QueryBuilderExpected{
				Expected: " ORDER BY name ASC",
				Values:   nil,
			},
		},
	}

	builder := db.BaseQueryBuilder{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got string
			var gotValues []interface{} = nil
			switch tt.method {
			case "Select":
				columns, values := builder.Select(tt.input.Columns, nil)
				got = got + "SELECT " + strings.Join(columns, ", ")
				gotValues = values
			case "From":
				got = builder.From(tt.input.Table.Name)
			case "Where":
				whereString, values := builder.Where(tt.input.ConditionGroups)
				got = whereString
				gotValues = values
			case "WhereGroup":
				whereString, values := builder.Where(tt.input.ConditionGroups)
				got = whereString
				gotValues = values
			case "Join":
				_, gotQuery := builder.Join(tt.input.Table.Name, tt.input.Joins)
				got = gotQuery
			case "OrderBy":
				got = builder.OrderBy(tt.input.Order)
			}
			if got != tt.expected.Expected {
				t.Errorf("expected '%s' but got '%s'", tt.expected, got)
			}

			if len(gotValues) != len(tt.expected.Values) {
				t.Errorf("expected '%v' but got '%v'", tt.expected.Values, gotValues)
			}

		})
	}
}

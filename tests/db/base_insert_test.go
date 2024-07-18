package db_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

func TestBaseInsertQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		input    *structs.InsertQuery
		expected QueryBuilderExpected
	}{
		{
			"Insert",
			"Insert",
			&structs.InsertQuery{
				Table: "users",
				Values: map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			QueryBuilderExpected{
				Expected: "INSERT INTO users (age, name) VALUES (?, ?)",
				Values:   []interface{}{"John", 30},
			},
		},
		{
			"InsertBatch",
			"InsertBatch",
			&structs.InsertQuery{
				Table: "users",
				ValuesBatch: []map[string]interface{}{
					{
						"name": "John",
						"age":  30,
					},
					{
						"name": "Mike",
						"age":  25,
					},
				},
			},
			QueryBuilderExpected{
				Expected: "INSERT INTO users (age, name) VALUES (?, ?), (?, ?)",
				Values:   []interface{}{30, "John", 25, "Mike"},
			},
		},
		{
			"InsertUsing",
			"InsertUsing",
			&structs.InsertQuery{
				Table:   "users",
				Columns: []string{"name", "age"},
				SelectQuery: &structs.Query{
					Table: structs.Table{Name: "profiles"},
					Joins: &[]structs.Join{},
					Columns: &[]structs.Column{
						{Name: "name"},
						{Name: "age"},
					},
					Conditions: &[]structs.Where{},
					ConditionGroups: &[]structs.WhereGroup{
						{
							Conditions: []structs.Where{
								{
									Column:    "age",
									Condition: ">",
									Value:     []interface{}{18},
								},
							},
							IsDummyGroup: true,
						},
					},
					Order: &[]structs.Order{},
					Lock:  &structs.Lock{},
				},
			},
			QueryBuilderExpected{
				Expected: "INSERT INTO users (name, age) SELECT name, age FROM profiles WHERE age > ?",
				Values:   []interface{}{18},
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
			case "Insert":
				got, gotValues = builder.Insert(tt.input)
			case "InsertBatch":
				got, gotValues = builder.InsertBatch(tt.input)
			case "InsertUsing":
				got, gotValues = builder.InsertUsing(tt.input)
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

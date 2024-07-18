package db_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

func TestBaseUpdateQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		input    *structs.UpdateQuery
		expected QueryBuilderExpected
	}{
		{
			"Update",
			"Update",
			&structs.UpdateQuery{
				Table: "users",
				Values: map[string]interface{}{
					"name": "Joe",
					"age":  30,
				},
				SelectQuery: &structs.Query{
					ConditionGroups: &[]structs.WhereGroup{
						{
							Conditions: []structs.Where{
								{
									Column:    "id",
									Condition: "=",
									Value:     []interface{}{1},
								},
							},
							IsDummyGroup: true,
						},
					},
					Joins: &[]structs.Join{},
				},
			},
			QueryBuilderExpected{
				Expected: "UPDATE users SET age = ?, name = ? WHERE id = ?",
				Values:   []interface{}{30, "Joe", 1},
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
			case "Update":
				got, gotValues = builder.BuildUpdate(tt.input)
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

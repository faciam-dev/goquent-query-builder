package db_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/database/postgres"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

func TestPostgreSQLQueryBuilder(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		input    structs.Query
		expected QueryBuilderExpected
	}{
		{
			"WhereFullText",
			"Where",
			structs.Query{
				ConditionGroups: []structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								FullText: &structs.FullText{
									Columns: []string{"name", "description"},
									Search:  "search",
									Options: map[string]interface{}{"mode": "websearch"},
								},
								Operator: consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
						Operator:     consts.LogicalOperator_AND,
					},
				},
			},
			QueryBuilderExpected{
				Expected: ` WHERE (to_tsvector($1, "name") || to_tsvector($2, "description")) @@ websearch_to_tsquery($3, $4)`,
				Values:   []interface{}{"english", "english", "english", "search"},
			},
		},
		{
			"WhereJsonContains",
			"Where",
			structs.Query{
				ConditionGroups: []structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:       "options->languages",
								JsonContains: &structs.JsonContains{Values: []interface{}{"en"}},
								Operator:     consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
						Operator:     consts.LogicalOperator_AND,
					},
				},
			},
			QueryBuilderExpected{
				Expected: ` WHERE ("options"->'languages')::jsonb @> $1`,
				Values:   []interface{}{"\"en\""},
			},
		},
		{
			"WhereJsonLength",
			"Where",
			structs.Query{
				ConditionGroups: []structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								Column:     "options->languages",
								JsonLength: &structs.JsonLength{Operator: ">", Value: 1},
								Operator:   consts.LogicalOperator_AND,
							},
						},
						IsDummyGroup: true,
						Operator:     consts.LogicalOperator_AND,
					},
				},
			},
			QueryBuilderExpected{
				Expected: ` WHERE jsonb_array_length(("options"->'languages')::jsonb) > $1`,
				Values:   []interface{}{1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			builder := postgres.NewPostgreSQLQueryBuilder()
			sb := make([]byte, 0, consts.StringBuffer_Middle_Query_Grow)

			var got string
			var gotValues []interface{} = nil
			switch tt.method {
			case "Select":
				values := builder.Select(&sb, tt.input.Columns, "", nil)
				columns := string(sb)
				got = got + "SELECT " + columns
				gotValues = values
			case "From":
				builder.From(&sb, tt.input.Table.Name)
				got = string(sb)
			case "Where":
				values, _ := builder.Where(&sb, tt.input.ConditionGroups)
				got = string(sb)
				gotValues = values
			case "WhereGroup":
				values, _ := builder.Where(&sb, tt.input.ConditionGroups)
				got = string(sb)
				gotValues = values
			case "Join":
				values := builder.Join(&sb, tt.input.Joins)
				got = string(sb)
				gotValues = values
			case "OrderBy":
				builder.OrderBy(&sb, tt.input.Order)
				got = string(sb)
			case "GroupBy":
				values := builder.GroupBy(&sb, tt.input.Group)
				got = string(sb)
				gotValues = values
			case "Limit":
				builder.Limit(&sb, tt.input.Limit)
				got = string(sb)
			case "Offset":
				builder.Offset(&sb, tt.input.Offset)
				got = string(sb)
			case "Limit_And_Offset":
				builder.Limit(&sb, tt.input.Limit)
				gotLimit := string(sb)
				sb = sb[:0]
				builder.Offset(&sb, tt.input.Offset)
				gotOffset := string(sb)
				got = gotLimit + gotOffset
			case "Lock":
				builder.Lock(&sb, tt.input.Lock)
				got = string(sb)
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

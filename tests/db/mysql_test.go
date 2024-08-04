package db_test

import (
	"strings"
	"testing"

	"github.com/faciam-dev/goquent-query-builder/database/mysql"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

func TestMySQLQueryBuilder(t *testing.T) {
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
				ConditionGroups: &[]structs.WhereGroup{
					{
						Conditions: []structs.Where{
							{
								FullText: &structs.FullText{
									Columns: []string{"name", "description"},
									Search:  "search",
									Options: map[string]interface{}{"mode": "boolean"},
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
				Expected: " WHERE MATCH (`name`, `description`) AGAINST (? IN BOOLEAN MODE)",
				Values:   []interface{}{"search"},
			},
		},
	}

	builder := mysql.NewMySQLQueryBuilder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sb := &strings.Builder{}

			var got string
			var gotValues []interface{} = nil
			switch tt.method {
			case "Select":
				values := builder.Select(sb, tt.input.Columns, "", nil)
				columns := sb.String()
				got = got + "SELECT " + columns
				gotValues = values
			case "From":
				builder.From(sb, tt.input.Table.Name)
				got = sb.String()
			case "Where":
				values := builder.Where(sb, tt.input.ConditionGroups)
				got = sb.String()
				gotValues = values
			case "WhereGroup":
				values := builder.Where(sb, tt.input.ConditionGroups)
				got = sb.String()
				gotValues = values
			case "Join":
				values := builder.Join(sb, tt.input.Joins)
				got = sb.String()
				gotValues = values
			case "OrderBy":
				builder.OrderBy(sb, tt.input.Order)
				got = sb.String()
			case "GroupBy":
				values := builder.GroupBy(sb, tt.input.Group)
				got = sb.String()
				gotValues = values
			case "Limit":
				builder.Limit(sb, tt.input.Limit)
				got = sb.String()
			case "Offset":
				builder.Offset(sb, tt.input.Offset)
				got = sb.String()
			case "Limit_And_Offset":
				builder.Limit(sb, tt.input.Limit)
				gotLimit := sb.String()
				sb.Reset()
				builder.Offset(sb, tt.input.Offset)
				gotOffset := sb.String()
				got = gotLimit + gotOffset
			case "Lock":
				got = builder.Lock(tt.input.Lock)
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

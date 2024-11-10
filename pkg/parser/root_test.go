package parser

import (
	db "BD/pkg/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BasicScenario(t *testing.T) {
	mockDb := db.NewDataBaseImpl()
	parse := ParserImpl{
		Databases: *mockDb,
	}
	emptyMockTable := db.NewTableImpl()
	abcMockTable := db.NewTableImpl()
	abcVal, _ := db.NewValue("abc", "11.11.2026T11:11:11")
	_, _ = abcMockTable.Insert("abc", abcVal)

	tests := []struct {
		scenarioName   string
		commands       []string
		errNo          []bool
		expectedResult []any
	}{
		{
			scenarioName: "create table1, insert value, get value",
			commands: []string{
				"DB select table1",
				"DB create table1",
				"DB select table1",
				"Table insert table1 abc abc 11.11.2026T11:11:11",
				"DB select table1",
				"Table get table1 abc",
			},
			errNo: []bool{
				true,
				false,
				false,
				false,
				false,
				false,
			},
			expectedResult: []any{
				"Таблица не найдена",
				true,
				*emptyMockTable,
				true,
				*abcMockTable,
				abcVal,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenarioName, func(t *testing.T) {
			for idx, cmd := range tc.commands {
				fmt.Printf("CMD: %s\t", cmd)

				res, err := parse.Parse(cmd)

				if tc.errNo[idx] {
					require.Error(t, err)
					assert.Equal(t, tc.expectedResult[idx], err.Error())

					fmt.Println("PASS")
					continue
				}

				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult[idx], res)

				fmt.Println("PASS")
			}
		})
	}
}

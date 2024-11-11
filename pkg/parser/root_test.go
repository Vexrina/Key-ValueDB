package parser

import (
	db "BD/pkg/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BasicScenario(t *testing.T) {
	emptyMockTable := db.NewTableImpl()
	abcMockTable := db.NewTableImpl()
	abcVal, _ := db.NewValue("abc", "11.11.2026T11:11:11")
	_, _ = abcMockTable.Insert("abc", abcVal)
	updatedAbcVal, _ := db.NewValue("abcd", "11.11.2027T11:11:11")

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
		{
			scenarioName: "create table1, insert value, create table2, insert same value to table2, table1=table2, drops tables",
			commands: []string{
				"DB create table1",
				"Table insert table1 abc abc 11.11.2026T11:11:11",
				"DB create table2",
				"Table insert table2 abc abc 11.11.2026T11:11:11",
				"Table get table1 abc",
				"Table get table2 abc",
				"Table delete table1 abc",
				"Table get table1 abc",
			},
			errNo: []bool{
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				true,
			},
			expectedResult: []any{
				true,
				true,
				true,
				true,
				abcVal,
				abcVal,
				true,
				"Ключа не существует",
			},
		},
		{
			scenarioName: "try add and update existing key",
			commands: []string{
				"DB create table1",
				"Table insert table1 abc abc 11.11.2026T11:11:11",
				"Table insert table1 abc abcd 11.11.2027T11:11:11",
				"Table get table1 abc",
				"Table update table1 abc abcd 11.11.2027T11:11:11",
				"Table get table1 abc",
			},
			errNo: []bool{
				false,
				false,
				true,
				false,
				false,
				false,
			},
			expectedResult: []any{
				true,
				true,
				"Такой ключ существует! Добавление невозможно",
				abcVal,
				true,
				updatedAbcVal,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenarioName, func(t *testing.T) {
			mockDb := db.NewDataBaseImpl()
			parse := ParserImpl{
				Databases: *mockDb,
			}

			for idx, cmd := range tc.commands {
				fmt.Printf("CMD:\t%s\t", cmd)

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

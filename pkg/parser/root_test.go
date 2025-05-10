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
				"DB db1 select table1",
				"DB db1 create table1",
				"DB db1 select table1",
				"Table db1 insert table1 abc abc 11.11.2026T11:11:11",
				"DB db1 select table1",
				"Table db1 get table1 abc",
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
				"DB db1 create table1",
				"Table db1 insert table1 abc abc 11.11.2026T11:11:11",
				"DB db1 create table2",
				"Table db1 insert table2 abc abc 11.11.2026T11:11:11",
				"Table db1 get table1 abc",
				"Table db1 get table2 abc",
				"Table db1 delete table1 abc",
				"Table db1 get table1 abc",
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
				"DB db1 create table1",
				"Table db1 insert table1 abc abc 11.11.2026T11:11:11",
				"Table db1 insert table1 abc abcd 11.11.2027T11:11:11",
				"Table db1 get table1 abc",
				"Table db1 update table1 abc abcd 11.11.2027T11:11:11",
				"Table db1 get table1 abc",
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
			mockDb := map[string]db.DataBaseImpl{"db1": *db.NewDataBaseImpl()}
			parse := ParserImpl{
				Databases: mockDb,
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

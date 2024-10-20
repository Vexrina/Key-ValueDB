package parser

/*
работаем с таблицами следующим образом:
- DB select TableName					|| 3 total len || 2 of arguments
- DB delete TableName					|| 3 total len || 2 of arguments
- DB create TableName					|| 3 total len || 2 of arguments
- DB rename OldTN NewTN					|| 4 total len || 3 of arguments

со значениями в таблице так:
	0	  1			2	  3    4	5(optional)
- Table delete TableName key	  		|| 4 total len || 3 of arguments
- Table insert TableName key value 	ttl	|| 5 total len || 4 of arguments
- Table get    TableName key	   		|| 4 total len || 3 of arguments
- Table update TableName key value 	ttl	|| 5 total len || 4 of arguments
- Table size   TableName		   		|| 3 total len || 2 of arguments
---
min len of command - 3 word
min len of arguments - 2 word
*/

const MINIMUM_LENGTH = 3

type Parser interface {
	Parse(command string) (any, error)
	parseDatabaseCommand(arguments []string) (any, error)
	parseTableCommand(arguments []string) (any, error)

	getDatabaseOperations() []string
}

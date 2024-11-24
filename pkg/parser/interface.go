package parser

/*
работаем с таблицами следующим образом:
- DB dbName select TableName			|| 4 total len || 3 of arguments
- DB dbName delete TableName			|| 4 total len || 3 of arguments
- DB dbName create TableName			|| 4 total len || 3 of arguments
- DB dbName rename OldTN NewTN			|| 5 total len || 4 of arguments

со значениями в таблице так:
	0	  1		2		  3    	 4	 5		6(optional)
- Table dbName delete TableName key	  			|| 5 total len || 4 of arguments
- Table dbName insert TableName key value 	ttl	|| 6 total len || 5 of arguments
- Table dbName get    TableName key	   			|| 5 total len || 4 of arguments
- Table dbName update TableName key value 	ttl	|| 6 total len || 5 of arguments
- Table dbName size   TableName		   			|| 4 total len || 3 of arguments
---
min len of command - 4 word
min len of arguments - 3 word
*/

const MINIMUM_LENGTH = 4

type Parser interface {
	Parse(command string) (any, error)
	parseDatabaseCommand(arguments []string) (any, error)
	parseTableCommand(arguments []string) (any, error)

	getDatabaseOperations() []string
}

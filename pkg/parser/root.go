package parser

import (
	db "BD/pkg/database"
	"fmt"
	"strings"
)

type ParserImpl struct {
	Databases db.DataBaseImpl
}

func (p *ParserImpl) Parse(command string) (any, error) {
	arguments := strings.Split(command, " ")
	if len(arguments) < MINIMUM_LENGTH {
		return nil, fmt.Errorf("there are not enough arguments in the command, got %d, need %d", len(arguments), MINIMUM_LENGTH)
	}
	switch arguments[0] {
	case "DB":
		return p.parseDatabaseCommand(arguments[1:])
	case "Table":
		return p.parseTableCommand(arguments[1:])
	default:
		return nil, fmt.Errorf("please specify what you want to perform the operation on, \"DB\" or \"Table\"")
	}
}

func (p *ParserImpl) parseDatabaseCommand(arguments []string) (any, error) {
	switch arguments[0] {
	case "select":
		return p.Databases.Select(arguments[1])
	case "delete":
		return p.Databases.Delete(arguments[1])
	case "create":
		table := db.NewTableImpl()
		return p.Databases.Create(arguments[1], table)
	case "rename":
		return p.Databases.Rename(arguments[1], arguments[2])
	default:
		return nil, fmt.Errorf(
			"please specify which operation you want to perform, here are the available operations: %v",
			p.getTableOperations(),
		)
	}
}

func (p *ParserImpl) parseTableCommand(arguments []string) (any, error) {
	operation := arguments[0] //1
	tableName := arguments[1] //2
	// key := arguments[2]//3
	// value := arguments[3]//4
	// ttl := arguments[4]//5
	table, err := p.Databases.Select(tableName)
	if err != nil {
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	switch operation {
	case "delete":
		return table.Delete(arguments[2])
	case "insert":
		val := db.Value{}
		// if len(arguments)==5{
		// 	val.Ttl = arguments[4] // parsing string -> time.Time
		// }
		val.Val = arguments[3]
		return table.Insert(arguments[2], val)
	case "get":
		return table.Get(arguments[2])
	case "update":
		val := db.Value{}
		// if len(arguments)==5{
		// 	val.Ttl = arguments[4] // parsing string -> time.Time
		// }
		val.Val = arguments[3]
		return table.Update(arguments[2], val)
	case "size":
		return table.Size(), nil
	case "parseTime":
		return nil, nil
	default:
		return nil, fmt.Errorf(
			"please specify which operation you want to perform, here are the available operations: %v",
			p.getTableOperations(),
		)
	}
}

func (p *ParserImpl) getTableOperations() []string {
	return []string{
		"create",
		"select",
		"delete",
		"rename",
	}
}

func (p *ParserImpl) getValueOperations() []string {
	return []string{
		"delete",
		"insert",
		"get",
		"update",
		"size",
		"parseTime",
	}
}

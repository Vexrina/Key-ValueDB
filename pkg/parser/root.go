package parser

import (
	db "BD/pkg/database"
	"fmt"
	"strings"
)

type ParserImpl struct {
	Databases map[string]db.DataBaseImpl
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
	dbName := arguments[0]

	database, ok := p.Databases[dbName]
	if !ok {
		return nil, fmt.Errorf("provided database doesnot exist")
	}

	operation := arguments[1]
	tableName := arguments[2]
	switch operation {
	case "select":
		return database.Select(tableName)
	case "delete":
		return database.Delete(tableName)
	case "create":
		table := db.NewTableImpl()
		return database.Create(tableName, *table)
	case "rename":
		return database.Rename(tableName, arguments[3])
	default:
		return nil, fmt.Errorf(
			"please specify which operation you want to perform, here are the available operations: %v",
			p.getTableOperations(),
		)
	}
}

func (p *ParserImpl) parseTableCommand(arguments []string) (any, error) {
	dbName := arguments[0]

	database, ok := p.Databases[dbName]
	if !ok {
		return nil, fmt.Errorf("provided database doesnot exist")
	}

	operation := arguments[1]
	tableName := arguments[2]
	// key := arguments[3]
	// value := arguments[3]
	// ttl := arguments[4]
	table, err := database.Select(tableName)
	if err != nil {
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	if len(arguments) < 6 {
		arguments = append(arguments, "")
	}

	switch operation {
	case "delete":
		return table.Delete(arguments[3])
	case "insert":
		val, errInsert := db.NewValue(arguments[4], arguments[5])
		if errInsert != nil {
			return nil, fmt.Errorf("db.NewValue error: %w", errInsert)
		}
		return table.Insert(arguments[3], val)
	case "get":
		return table.Get(arguments[3])
	case "update":
		val, errInsert := db.NewValue(arguments[4], arguments[5])
		if errInsert != nil {
			return nil, fmt.Errorf("db.NewValue error: %w", errInsert)
		}
		return table.Update(arguments[3], val)
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

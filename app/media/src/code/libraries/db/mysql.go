package db

import (
	"database/sql"
	"fmt"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/helper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"reflect"
	"strings"
)

func MysqlGetWhere(connection Connection, whereClause string, sort string, limit int) ([]map[string]interface{}, error) {

	database, databaseError := sql.Open("mysql", connection.UserName+":"+connection.Password+"@tcp("+connection.IpAddress+":"+connection.Port+")/"+connection.Database)

	if databaseError != nil {
		return nil, databaseError
	}

	defer database.Close()

	sqlQuery := "SELECT * FROM " + connection.Table + " WHERE " + whereClause
	rows, queryError := database.Query(sqlQuery)

	if queryError != nil {
		return nil, queryError
	}

	// be careful deferring Queries if you are using transactions
	defer rows.Close()

	var returnCollection []map[string]interface{}

	cols, _ := rows.ColumnTypes()

	pointers := make([]interface{}, len(cols))
	modelInstance := make(map[string]interface{}, len(cols))

	for index, column := range cols {
		var value interface{}

		switch column.DatabaseTypeName() {
		case "INT":
			value = new(sql.NullInt32)
		case "VARCHAR", "STRING":
			value = new(sql.NullString)
		case "TIMESTAMP", "DATETIME":
			value = new(sql.NullString)
		case "UUID":
			value = new(sql.NullString)
		default:
			value = new(interface{}) // destination must be a pointer
		}

		modelInstance[column.Name()] = value
		pointers[index] = value
	}

	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			log.Println(err)
		}
		returnCollection = append(returnCollection, modelInstance)
	}

	return returnCollection, nil
}

func MysqlCreateNew(connection Connection, model []SqlExecModel) ([]map[string]interface{}, error) {

	database, databaseError := sql.Open("mysql", connection.UserName+":"+connection.Password+"@tcp("+connection.IpAddress+":"+connection.Port+")/"+connection.Database)

	if databaseError != nil {
		return nil, databaseError
	}

	defer database.Close()

	fieldsForInsertion := buildFieldsFromModel(model)
	valuesForInsertion := buildValuesFromModel(model)

	sqlQuery := "INSERT INTO " + connection.Table + " (" + fieldsForInsertion + ") VALUES (" + valuesForInsertion + ")"
	result, execError := database.Exec(sqlQuery)

	if execError != nil {
		return nil, execError
	}

	newestId, _ := result.LastInsertId()
	newRecordQuery := connection.PrimaryKey + " = " + fmt.Sprint(newestId)

	entityCollection, error := MysqlGetWhere(connection, newRecordQuery, "ASC", 1)

	if error != nil {
		return nil, error
	}

	return entityCollection, nil
}

func buildFieldsFromModel(model []SqlExecModel) string {
	var modelFields []string

	for _, currFieldModel := range model {
		if currFieldModel.Value.IsZero() && currFieldModel.Field != "sys_row_id" {
			continue
		}
		modelFields = append(modelFields, currFieldModel.Field)
	}

	return strings.Join(modelFields, ",")
}

func buildValuesFromModel(model []SqlExecModel) string {
	var modelFields []string

	for _, currFieldModel := range model {
		if currFieldModel.Value.IsZero() && currFieldModel.Field != "sys_row_id" {
			continue
		}
		var currType string

		switch currFieldModel.Type {
		case reflect.String:
			currType = "\"" + fmt.Sprint(currFieldModel.Value) + "\""
			break
		case reflect.Int, reflect.Int32, reflect.Int64:
			currType = fmt.Sprint(currFieldModel.Value)
			break
		case reflect.Bool:
			if currFieldModel.Value.Bool() == true {
				currType = "true"
			} else {
				currType = "false"
			}
			break
		case reflect.Struct:
			switch currFieldModel.Value.Type().String() {
			case "helper.NullTime":
				currDateTime := currFieldModel.Value.Interface().(helper.NullTime)
				currType = "\"" + currDateTime.Value.Format("2006-01-02 03:04:05") + "\""
			}
			break
		default:
			switch currFieldModel.Value.Type().String() {
			case "uuid.UUID":
				if !currFieldModel.Value.IsZero() {
					currType = "\"" + currFieldModel.Value.Interface().(uuid.UUID).String() + "\""
				} else {
					currType = "\"" + uuid.New().String() + "\""
				}
			}
		}

		modelFields = append(modelFields, currType)
	}

	return strings.Join(modelFields, ",")
}

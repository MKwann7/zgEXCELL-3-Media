package db

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Connection struct {
	Table      string
	PrimaryKey string
	UuidKey    string
	IpAddress  string
	Port       string
	Database   string
	UserName   string
	Password   string
	DbType     string
}

const MySQL = "mysql"
const Postgres = "extensions"

func (connection *Connection) GetMain(tableName string, userKey string, uuidKey string) Connection {
	return Connection{
		os.Getenv("MAIN_DB_NAME") + "." + tableName,
		userKey,
		uuidKey,
		os.Getenv("MAIN_DB_HOST"),
		os.Getenv("MAIN_DB_PORT"),
		os.Getenv("MAIN_DB_NAME"),
		os.Getenv("MAIN_DB_USER"),
		os.Getenv("MAIN_DB_PASS"),
		MySQL}
}

func (connection *Connection) GetMedia(tableName string, userKey string, uuidKey string) Connection {
	return Connection{
		os.Getenv("MEDIA_DB_NAME") + "." + tableName,
		userKey,
		uuidKey,
		os.Getenv("MEDIA_DB_HOST"),
		os.Getenv("MEDIA_DB_PORT"),
		os.Getenv("MEDIA_DB_NAME"),
		os.Getenv("MEDIA_DB_USER"),
		os.Getenv("MEDIA_DB_PASS"),
		MySQL}
}

func (connection *Connection) GetTraffic(tableName string, userKey string, uuidKey string) Connection {
	return Connection{
		os.Getenv("TRAFFIC_DB_NAME") + "." + tableName,
		userKey,
		uuidKey,
		os.Getenv("TRAFFIC_DB_HOST"),
		os.Getenv("TRAFFIC_DB_PORT"),
		os.Getenv("TRAFFIC_DB_NAME"),
		os.Getenv("TRAFFIC_DB_USER"),
		os.Getenv("TRAFFIC_DB_PASS"),
		MySQL}
}

func (connection *Connection) GetNotification(tableName string, userKey string, uuidKey string) Connection {
	return Connection{
		tableName,
		userKey,
		uuidKey,
		os.Getenv("NOTIFY_DB_HOST"),
		os.Getenv("NOTIFY_DB_PORT"),
		os.Getenv("NOTIFY_DB_NAME"),
		os.Getenv("NOTIFY_DB_USER"),
		os.Getenv("NOTIFY_DB_PASS"),
		Postgres}
}

type SqlExecModel struct {
	Field string
	Type  reflect.Kind
	Value reflect.Value
}

func MakeSqlExecModel(myStructFields reflect.Type, myStructValues reflect.Value, indexField string) []SqlExecModel {

	num := myStructFields.NumField()

	var myModel []SqlExecModel

	for i := 0; i < num; i++ {
		field := myStructFields.Field(i)
		value := myStructValues.Field(i)

		currFieldName := convertToSqlField(field.Name)

		if currFieldName != indexField {
			sqlExecModel := SqlExecModel{Field: convertToSqlField(field.Name), Type: value.Kind(), Value: value}
			myModel = append(myModel, sqlExecModel)
		}
	}

	return myModel
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func convertToSqlField(fieldName string) string {
	snake := matchFirstCap.ReplaceAllString(fieldName, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

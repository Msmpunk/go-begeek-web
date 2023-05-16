package database

import (
	"GoWeb/app"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/lib/pq"
)

func Migrate(app *app.App, anyStruc interface{}) error {

	valueOfStruct := reflect.ValueOf(anyStruc)
	typeOfStruct := valueOfStruct.Type()

	tableName := typeOfStruct.Name()

	err := createTable(app, tableName)

	if err != nil {
		return err
	}

	for i := 0; i < valueOfStruct.NumField(); i++ {
		fieldtype := typeOfStruct.Field(i)
		fieldName := fieldtype.Name
		if fieldName != "ID" && fieldName != "id" {
			err := createColumn(app, tableName, fieldName, fieldtype.Type.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createTable(app *app.App, tableName string) error {

	var tableExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE c.relname ~ $1 AND pg_catalog.pg_table_is_visible(c.oid))", "^"+tableName+"$").Scan(&tableExists)

	if err != nil {
		log.Println("Error validando si existe la tabla " + tableName)
		return err
	}

	if tableExists {
		log.Println("La tabla ya existe" + tableName)
		return nil
	} else {

		sanitizedTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\"Id\" serial primary key)", tableName)

		_, err := app.Db.Query(sanitizedTableQuery)

		if err != nil {
			log.Println("Error creando la tabla:  " + tableName)
			return err
		}

		log.Println("La tabla se creo perfectamente: " + tableName)
		return nil

	}
}

func createColumn(app *app.App, tableName, columnName, columnType string) error {

	var columnExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", tableName, columnName).Scan(&columnExists)

	if err != nil {
		log.Println("Error creando la columna :  " + columnName + " En La tabla: " + tableName)
		return err
	}

	if columnExists {
		log.Println("La columna :  " + columnName + " ya existe en En La tabla: " + tableName)
		return nil
	} else {

		postgresType, err := getPostgresType(columnType)

		if err != nil {
			log.Println("Error creando la columna :  " + columnName + " En La tabla: " + tableName + "Con el tipo: " + columnType)
			return err
		}

		sanitazedTableName := pq.QuoteIdentifier(tableName)

		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS \"%s\" %s", sanitazedTableName, columnName, postgresType)

		_, err = app.Db.Query(query)

		if err != nil {
			log.Println("Error editando la columna :  " + columnName + " En La tabla: " + tableName + "Con el tipo: " + columnType)
			return err
		}

		log.Println("La Columna se creo corectamente", columnName)

		return nil

	}
}

func getPostgresType(goType string) (string, error) {
	switch goType {
	case "int", "int32", "uint", "uint32":
		return "integer", nil
	case "int64", "uint64":
		return "biginit", nil
	case "int16", "int8", "uint16", "uint8", "byte":
		return "smallint", nil
	case "string":
		return "text", nil
	case "float64":
		return "dooble precision", nil
	case "bool":
		return "boolean", nil
	case "Time":
		return "tymestamp", nil
	case "[]byte":
		return "byte", nil
	}

	return "", errors.New("Tipo no reconocido" + goType)
}

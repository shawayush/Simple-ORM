// Implimentation of Simple ORM
// by Ayush Shaw

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	_dataBaseName    = ""
	_tableName       = ""
	db               *sql.DB
	_origionalStruct []User
)

type (
	User struct { //change the attributes of the strucutre not the name
		Name     string `sql_type:"varchar(100)"`
		Id       int    `sql_type:"int"`
		UserName string `sql_type:"varchar(100)"`
		Password string `sql_type:"varchar(100)"`
	}

	UserMigration struct { //Can be used to check for migration change the attribute not the name
		Name     string `sql_type:"varchar(100)"`
		Id       int    `sql_type:"int"`
		UserName string `sql_type:"varchar(100)"`
		Password string `sql_type:"varchar(100)"`
		Address  string `sql_type:"varchar(100)"`
		Phone    int    `sql_type:"int"`
	}
)

// Delete row creates a map of which is row would
// be deleted in the Table and then the mapped
// data is passed to DeleteRowQuery which forms
// the delete function.
func (strucutre User) DeleteRow() {

	e := reflect.ValueOf(&strucutre).Elem()
	valueMap := make(map[interface{}]interface{})

	for i := 0; i < e.NumField(); i++ {

		varName := e.Type().Field(i).Name
		varValue := e.Field(i).Interface()

		valueMap[varName] = varValue
	}

	sQuery := DeleteRowQuery(valueMap)
	ExecuteQuery(sQuery, " ROW DELETION ")
}

// Creates a Query for deleting the row in the table
func DeleteRowQuery(mapData map[interface{}]interface{}) string {

	deleteQuery := " DELETE FROM " + _tableName + " WHERE (1=1) "

	for k, v := range mapData {
		deleteQuery += " AND " + fmt.Sprintf("%v", k) + " = \"" + fmt.Sprintf("%v", v) + "\" "
	}

	return deleteQuery

}

// Update function update the row given by the user.
// Intially the struct map is stored in the memory and
// the both the struct are compared and checked to
// which will be passed through UpdateRowInTableV2
// and there the query is formed for the update table.
func (newStruct User) Update() {

	status := false

	for _, origionalStruct := range _origionalStruct {

		oldE := reflect.ValueOf(&origionalStruct).Elem()
		newE := reflect.ValueOf(&newStruct).Elem()

		newMap := make(map[interface{}]interface{})
		oldMap := make(map[interface{}]interface{})

		if oldE.NumField() != newE.NumField() {
			fmt.Println("Something went wrong, please try again")
			return
		}

		lenght := oldE.NumField()
		checkIntLenght := 0
		for i := 0; i < lenght; i++ {

			oldVarName := oldE.Type().Field(i).Name
			newVarName := newE.Type().Field(i).Name

			oldVarValue := oldE.Field(i).Interface()
			newVarValue := newE.Field(i).Interface()

			oldMap[oldVarName] = oldVarValue
			newMap[newVarName] = newVarValue

		}

		if reflect.DeepEqual(oldMap, newMap) {
			//nothing to do just pass by
		} else {
			for k1, v1 := range newMap {
				for k2, v2 := range oldMap {
					if k1 == k2 {
						if v1 == v2 {
							checkIntLenght += 1
						}
					}
				}
			}
		}

		if lenght-1 == checkIntLenght {
			UpdateRowInTableV2(newStruct, origionalStruct)
			status = true
			break
		}
	}

	if !status {
		fmt.Println("Nothing to Update!")
	}
}

//Compare the two strucutres for the update that is
//the old strucutre and the new strucutre, and gets the
//map for the same, which is again used in the function
//to get the Update Query.
func UpdateRowInTableV2(newStruct User, origionalStruct User) {

	updateMap := make(map[interface{}]interface{})
	sameMap := make(map[interface{}]interface{})

	oldE := reflect.ValueOf(&origionalStruct).Elem()
	newE := reflect.ValueOf(&newStruct).Elem()

	for i := 0; i < oldE.NumField(); i++ {

		oldVarName := oldE.Type().Field(i).Name
		newVarName := newE.Type().Field(i).Name

		oldVarValue := oldE.Field(i).Interface()
		newVarValue := newE.Field(i).Interface()

		if oldVarValue == newVarValue {
			sameMap[oldVarName] = oldVarValue
		} else {
			updateMap[newVarName] = newVarValue
		}

	}

	sQuery := UpdateRowQuery(sameMap, updateMap)
	ExecuteQuery(sQuery, "Updation")
}

//Returns the query for the Update row
func UpdateRowQuery(sameMap map[interface{}]interface{}, updateMap map[interface{}]interface{}) string {

	updateQuery := " UPDATE " + _tableName + " SET "

	for k, v := range updateMap {
		updateQuery += fmt.Sprintf("%v", k) + " = \"" + fmt.Sprintf("%v", v) + "\" "
	}

	updateQuery += " WHERE (1=1) "

	for k, v := range sameMap {
		updateQuery += " AND " + fmt.Sprintf("%v", k) + " = \"" + fmt.Sprintf("%v", v) + "\" "
	}
	return updateQuery

}

// Alters the schemma for the Table
func Migration(structure interface{}) {

	query := "SELECT * FROM " + _tableName + " WHERE (1=1) "

	status := false
	var columnSlice []string
	var migrationTable []interface{}
	rows, err := db.Query(query)
	columns, err := rows.Columns()
	params1 := GetStructureAttribute(structure)
	if err != nil {
		panic(err)
	}

	for _, col := range columns {

		stng := fmt.Sprintf("%v", col)
		columnSlice = append(columnSlice, stng)
	}

	for _, param := range params1 {

		stng := fmt.Sprintf("%v", param)
		firstWordString := FirstWord(stng, 1)

		if ContainsString(columnSlice, firstWordString) {
			//nothing to append, let it pass
		} else {
			migrationTable = append(migrationTable, param)
			status = true
		}
	}

	if status == false {
		fmt.Println("Noting To migrate!")
		return
	}

	newQuery := AlterTable(migrationTable)
	newQuery = strings.TrimSuffix(newQuery, ",")

	ExecuteQuery(newQuery, "Migration")

}

//returns a query for Alter table used in Migrations function
func AlterTable(tables []interface{}) string {

	i := 0
	query := "ALTER TABLE " + _tableName
	for _, table := range tables {
		if i == 0 {
			//To add table name here
			//or just pass
			i += 1
		} else {
			query += " ADD " + fmt.Sprintf("%v", table) + ","
		}
	}

	return query

}

//check if the slice contains, particular string
func ContainsString(s []string, str string) bool {

	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//Gives the first word from a series of string
func FirstWord(value string, count int) string {

	for i := range value {
		if value[i] == ' ' {
			count -= 1
			if count == 0 {
				return value[0:i]
			}
		}
	}
	return value
}

// Fetch all the rows from the table or
// fetch with a filter or
// fetch using Lazy query method
func FetchData(vals ...interface{}) []User {

	if len(vals) > 1 {
		fmt.Println("Too many Parameters, Please remove extra parameters!")
	}

	sql := "SELECT DISTINCT * FROM " + _tableName + " WHERE (1=1)"

	for _, val := range vals {
		switch val.(type) {
		case string:
			sql += " AND " + fmt.Sprintf("%v", val)
			break
		default:
			fmt.Println("Value not in proper format!")
			break
		}
	}

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	var p User
	rows, err := db.Query(sql)
	columns, err := rows.Columns()

	var allMaps []map[string]interface{}
	var jsonString string
	listStruct := []User{}

	for rows.Next() {

		jsonString = "{"
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			panic(err)
		}

		resultMap := make(map[string]interface{})

		for i, val := range values {

			if val != nil {
				myString := val.([]uint8) // strings are stored in runes, converting them back to string
				var checkIntiger string
				status := false

				for i := 0; i < len(myString); i++ {
					if _, err := strconv.Atoi(string(myString[i])); err == nil {
						checkIntiger += string(myString[i])
						status = true
					}
				}

				if status {
					intVar, _ := strconv.Atoi(checkIntiger)
					resultMap[string(columns[i])] = intVar
				} else {
					resultMap[string(columns[i])] = string(myString)
				}

			}

		}

		for k, v := range resultMap {
			var dummyString string
			dummyString = k

			if reflect.TypeOf(v).String() == "int" {
				jsonString += "\"" + dummyString + "\"" + ":" + fmt.Sprintf("%v", v) + ", "
			} else {
				jsonString += "\"" + dummyString + "\"" + ":" + "\"" + v.(string) + "\", "
			}
		}

		jsonString += "}"
		jsonString = strings.Replace(jsonString, ", }", "}", 2)
		bytes := []byte(jsonString)
		err = json.Unmarshal(bytes, &p)
		if err != nil {
			panic(err)
		}

		allMaps = append(allMaps, resultMap)
		listStruct = append(listStruct, p)

	}

	var fetchedData []string
	_origionalStruct = listStruct

	for _, row := range allMaps {
		relevantString := fmt.Sprintf("%v", row)
		stringLength := len(relevantString) - 1
		fetchedData = append(fetchedData, relevantString[4:stringLength])
	}

	for _, row := range fetchedData {
		fmt.Println(row)
	}

	return listStruct

}

//Delete the table created in the database
func DeleteTable(table string) {

	query := "DROP TABLE " + table
	ExecuteQuery(query, "Deleteion")

}

//This Inserts data into TABLE
func InsertIntoTable(structure interface{}) {

	query := "INSERT INTO " + _tableName + " VALUES ( "

	if reflect.ValueOf(structure).Kind() == reflect.Struct {
		value := reflect.ValueOf(structure)
		numberOfFields := value.NumField()
		for i := 0; i < numberOfFields; i++ {
			query += "'" + fmt.Sprintf("%v", value.Field(i)) + "' " + ","
		}
		query += ")"
		ExecuteQuery(strings.Replace(query, ",)", ")", 2), "Insertion")
	}
}

//This Also Inserts into Table, but also checks, if the user has
//provided the right tables to be inserted also throws the
//error if Tables are not matched
//only problem is that, I have to User as data-type instead
//of interface, which limits the usage structure
func InsertIntoTableV2(structure User) {

	query := "INSERT INTO " + _tableName + " ( "
	s := reflect.ValueOf(&structure).Elem()
	typeOfs := s.Type()

	sql := "SELECT * FROM " + _tableName + " WHERE (1=1)"
	rows, err := db.Query(sql)
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	for i := 0; i < s.NumField(); i++ {
		//f := s.Field(i)

		if ContainsString(columns, typeOfs.Field(i).Name) { //ContainsString(columns, typeOfs.Field(i).Name)
			query += " " + typeOfs.Field(i).Name + ","
		} else {
			fmt.Println(typeOfs.Field(i).Name + " Column Not in Table " + _tableName)
			break
		}
	}
	query += ")"
	query = strings.Replace(query, ",)", ")", 2)
	query += " VALUES ("

	for i := 0; i < s.NumField(); i++ {

		f := s.Field(i)
		query += "'" + fmt.Sprintf("%v", f.Interface()) + "' ,"
	}
	query += ")"
	query = strings.Replace(query, ",)", ")", 2)

	ExecuteQuery(query, "Insertion")
}

//Use to Execute queries throughout the code
func ExecuteQuery(query string, operation string) {

	_, err := db.Exec("USE " + _dataBaseName)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(operation + " Operation Successful")
	}
}

// Get the properties of a strucuture, using reflect
// most useful, to get the schema from the struct used above.
func GetStructureAttribute(domain interface{}) (params []interface{}) {

	val := reflect.ValueOf(domain)

	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	structType := val.Type()
	tableName := structType.Name()
	params = append(params, tableName)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag
		fieldName := field.Name
		fieldType := tag.Get("sql_type")

		paramstring := fieldName + " " + fieldType
		params = append(params, paramstring)
	}
	return params
}

//Inserts Table into Database, also gets the name of the name of the table
func createTableInDataBase() {
	var i int
	i = 0
	query := ""
	//to check once
	params := GetStructureAttribute(&User{})

	for _, param := range params {
		if i == 0 {
			query = "CREATE TABLE IF NOT EXISTS " + fmt.Sprintf("%v", param) + " ( "
			i += 1
			_tableName = fmt.Sprintf("%v", param)
		} else {
			query += fmt.Sprintf("%v", param) + ","
		}
	}
	query += ")"
	output := strings.Replace(query, ",)", ")", 2)

	ExecuteQuery(output, "Table Insertion ")

}

//Create a Database
func createDataBase(name string) {

	_dataBaseName = name
	query := "CREATE DATABASE IF NOT EXISTS " + name

	ExecuteQuery(query, "Database Creation ")
	ChangeDatabase(name)

}

//Change to Database
func ChangeDatabase(dataBasename string) {
	_dataBaseName = dataBasename
	_, err := db.Exec("USE " + dataBasename)
	if err != nil {
		panic(err)
	}

}

//Intiliae Database connection with the controller
func InitDB(dataSourceName string) error {
	var err error

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	return db.Ping()
}

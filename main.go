package main

import "fmt"

var _dataSourceName = "shaw:password@tcp(localhost:3306)/" /// Change according to the database in your PC. USED sql here

func main() {

	//1. for initiliation of Database (keep it uncommented, as aditional check is applied not to rewite the database)
	InitDB(_dataSourceName)
	createDataBase("testdb")
	ChangeDatabase("testdb")

	//2. for creating table into the DB (keep it uncommented, as aditional check is applied not to rewite the database)
	createTableInDataBase()

	//3. Deleting table from the database, you can add your table
	//DeleteTable(_tableName)

	//4. Inserting data into the table
	/*
		person := User{
			Id:       3,
			Name:     "Ryan",
			UserName: "huga",
			Password: "hugaRyan",
		}
		//InsertIntoTable(person)        //Note: Insert before Migration works
		InsertIntoTableV2(person) //Note: Works all the time but function takes datatype User from the above
	*/
	//Note: both the function can used but InsertIntoTableV2 has a check if the right
	//data is inserted but it checks from the structure declared above.
	// if the data attributes used is right or not

	//5.Fetching From query

	// 5.1 Fetches all the data
	//FetchData()

	// 5.2 Fetches data with the filter should be in same as given
	//FetchData("Name = 'Ayush'")

	//6. Update Row
	rows := FetchData()
	row := rows[0]
	fmt.Println(row)
	row.Id = 12
	row.Update()
	//OR
	row = FetchData()[0]
	row.Id = 1
	fmt.Println(row)
	row.Update()

	//7. Update through lazy filter
	row = FetchData("Name = 'Ayush'")[0]
	fmt.Println(row)
	row.UserName = "cosmos"
	row.Update()

	//8. Delete a row
	row = FetchData()[1]
	fmt.Println(row)
	row.DeleteRow()

	//9. Migration
	//Migration(&UserMigration{})
	//Note Migartion will only add the table
}

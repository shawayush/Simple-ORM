# Simple-ORM
A implimentation of Object Relation Mapping in Golang


## Introduction

ORM stands for **Object Relational Mapping** where: 
- The **Object** part is the one with the programming Language (Here it is golang)
- The **Relational** is the Database Management System (DBMS) (Here it is MYSQL)
- **Mapping** Stands for the bridge between the two that is Object and Tables.

## Prerequisites
- Golang (1.15 or higher)
- MySQL (4.1 or higher)

## Installation
- You can install **golang** from https://go.dev/doc/install
- You can install **MYSQL** from https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/
- Additionally you need to install go-sql driver, which you can do by installing the package to your $GOPATH with the go tool from shell:
```
$ go get -u github.com/go-sql-driver/mysql
```
Make sure Git is installed on your machine and in your system's `PATH`.

## Working Around the Code

Clone the repository in your local directory.

To run the code in terminal write:
```
go run main.go ORM.go
```

Add the following commands in main.go to use the functions

Change the **_dataSourceName** at main.go and it should be in the format
Data Source Name (DSN)

`[username[:password]@][protocol[(address)]]/`

eg:"shaw:password@tcp(localhost:3306)/"

### Initialization
1. For initialization of Database (keep it uncommented, as additional check is applied not to rewrite the database)
```
InitDB(_dataSourceName)
createDataBase("testdb")
ChangeDatabase("testdb")
```
### Create Table
2. for creating table into the DB (keep it uncommented, as additional check is applied not to rewrite the database)
```
createTableInDataBase()
```
### Delete Table
3. Deleting table from the database, you can add your table name.
```
DeleteTable(_tableName)
```
### Insert Into Table
4. Inserting data into the table
NOTE: Use the struct Name (User) defined in ORM.go (You can change the attribute)
```
	person := User{          
		Id:       11,
		Name:     "Ayush",
		UserName: "Shaw",
		Password: "Shawayush",
	}
  
InsertIntoTable(person)       Note: Insert before Migration works

InsertIntoTableV2(person)     Note: Works all the time but function takes datatype User from the above
```

Note: both the function can used but InsertIntoTableV2 has a check if the right data is inserted 
but it checks from the structure declared above. If the data attributes used is right or not

### Fetch
5. Fetching From query

  - Different Fetches 
    - Fetch All Data
  
	  `FetchData()`

    - Fetches data with the Additional filter (lazyQuery)
	
	  `FetchData("Name = 'Ayush'")`
	  
### Operation on Fetched Data
6. Update row
```
		rows := FetchData()
		row := rows[0]
		row.Id = 10
		row.Update()
```
	  	
**OR**

You can directly update it from Row
```
		row := FetchData()[0]
		row.Id = 10
		row.Update()
```
7. Update row with Filter (lazyQuery)
```
		row := FetchData("Name = 'Ayush'")[0]
		row.Id = 10
		row.Update()
```
Note: You can only update one Row value at a time

8. Delete a Row
You can delete a row after fetching the row Data and using the function `DeleteRow()`
```
		row := FetchData()[0]
		row.DeleteRow()
```

### Migration

9. Migration
```	
    Migration(&UserMigration{})	 
```
Note: Use the structure defined in ORM.go (You can change the attribute)  
Note: Migration will only add the table


package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//var dbserver string
	//dbserver = "127.0.0.1"
	condb, errdb := sql.Open("mssql", "server=localhost;user id=sa;password=P@ssword;")
	if errdb != nil {
		//sprawdza tylko pierwszy parametr! Czyli mssql
		log.Fatal("Błąd spowodowany przez: ", errdb.Error())
	}

	var sqlversion string
	rows, err := condb.Query("select @@version")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err1 := rows.Scan(&sqlversion)
		if err1 != nil {
			log.Fatal(err1)
		}

		log.Println(sqlversion)
	}

	defer condb.Close()
}

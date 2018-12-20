package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = "localhost"
var user = "sa"
var port = 1433
var password = "P@ssword"
var database = "testdb"

func main() {
	cnstr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	fmt.Println(cnstr)

	iloscr, err := czytajdb(cnstr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Ilość rekordów: %d\n", iloscr)

	var noweid int64
	var errdodawania error
	noweid, errdodawania = dodaj(cnstr, "Jas", "Fasola")
	if errdodawania != nil {
		log.Fatal(errdodawania)
	}
	log.Printf("NowododanyID to: %d", noweid)

	var zaktualizowaneid int
	var erraktu error

	zaktualizowaneid, erraktu = aktualizuj(cnstr, 1, "Dupa", "Jasiu")
	if erraktu != nil {
		log.Fatal(erraktu)
	}
	log.Printf("Zaktualizowany ID to: %d", zaktualizowaneid)
	//defer db.Close()
}

func czytajdb(cnstr string) (int, error) {
	var err error
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Fatal("Błąd sterownika: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {
		log.Fatal("Problem puli połączenia: ", err.Error())
	}

	sqlq := "select * from TestSchema.Employees"
	rows, errq := db.QueryContext(ctx, sqlq)
	if errq != nil {
		log.Fatal("Problem zapytania: ", errq.Error())
	}
	count := 0
	for rows.Next() {
		var id int
		var name, loc string

		errqw := rows.Scan(&id, &name, &loc)
		if errqw != nil {
			log.Fatal("Błąd w wyniku zapytania", errqw.Error())
			return -1, errqw
		}
		fmt.Printf("id: %d, name: %s, location: %s\n", id, name, loc)
		count++
	}

	defer rows.Close()
	return count, errq
}

func dodaj(cnstr string, name string, location string) (int64, error) {
	var err error
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Fatal("Błąd sterownika: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {
		log.Fatal("Problem puli połączenia: ", err.Error())
	}

	var zapins = "insert into TestSchema.Employees (name, location) values (@name,@location); select @@identity;"

	skladnia, err := db.Prepare(zapins)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()
	//var name, location string
	row := skladnia.QueryRowContext(ctx,
		sql.Named("name", name),
		sql.Named("location", location))

	var noweid int64
	err = row.Scan(&noweid)
	if err != nil {
		log.Fatal("Błąd w wyniku zapytania: ", err.Error())
	}
	log.Println("Gicior")
	return noweid, err
}

func aktualizuj(cnstr string, id int, name string, location string) (int, error) {
	var err error
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Fatal("Błąd sterownika: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {
		log.Fatal("Problem puli połączenia: ", err.Error())
	}
	var aktualizacja = "update Testschema.Employees set name=@name, location=@location where id=@id"

	skladnia, err := db.Prepare(aktualizacja)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()

	newakt := fmt.Sprintf("update Testschema.Employees set name='%s', location='%s' where id=%d;", name, location, id)

	result, errctx := skladnia.ExecContext(ctx, newakt,
		sql.Named("name", name),
		sql.Named("location", location),
		sql.Named("id", id),
	)

	if errctx != nil {
		log.Fatal("Problem z execcontext: ", errctx.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatal("Panika ", err.Error())
	}

	log.Println("Gicior")
	return id, err
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = getEnv("SQL_SERVER", "localhost")
var user = getEnv("SQL_USER", "sa")
var port = getEnv("SQL_PORT", "1433")
var password = getEnv("SQL_PASSWORD", "P@ssword")
var database = getEnv("SQL_DATABASE", "testdb")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	cnstr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)
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

	zaktualizowaneid, erraktu = aktualizuj(cnstr, 3, "Pupa", "Jasiu")
	if erraktu != nil {
		log.Fatal(erraktu)
	}
	log.Printf("Zaktualizowany ID to: %d", zaktualizowaneid)

	var usunieteid int
	var errusun error

	usunieteid, errusun = usun(cnstr, 3)
	if errusun != nil {
		log.Fatal(errusun)
	}
	log.Printf("Usunięty ID to: %d", usunieteid)

	//defer db.Close()
}

func czytajdb(cnstr string) (int, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
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
	connclose(db)
	return count, errq
}

func dodaj(cnstr string, name string, location string) (int64, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
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
	connclose(db)
	return noweid, err
}

func aktualizuj(cnstr string, id int, name string, location string) (int, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	var aktualizacja = "update Testschema.Employees set name=@name, location=@location where id=@id"

	skladnia, err := db.Prepare(aktualizacja)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()

	newakt := fmt.Sprintf("update Testschema.Employees set name='@name', location='@location' where id=@id;")
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
	connclose(db)
	return id, err

}

func usun(cnstr string, id int) (int, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}

	statement := fmt.Sprintf("delete from Testschema.Employees where id=@id")

	skladnia, err := db.Prepare(statement)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()

	result, errctx := skladnia.ExecContext(ctx, statement,
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
	connclose(db)
	return id, err

}

func connopen(cnstr string) (*sql.DB, context.Context, error) {
	var err error
	db, err := sql.Open("sqlservers", cnstr)

	if err != nil {
		log.Println("Błąd sterownika: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {
		log.Println("Problem puli połączenia: ", err.Error())
	}

	return db, ctx, err
}

func connclose(db *sql.DB) {
	defer db.Close()
}

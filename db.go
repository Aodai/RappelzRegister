package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var db *sql.DB

func connectToDB() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Database)
	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected to the database!\n")
}

func userExists(user string) (bool, error) {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		return false, err
	}
	tsql := fmt.Sprintf("SELECT * FROM dbo.Accounts WHERE login_name='%s';", user)
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		count++
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func createUser(user, pass, ip string) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateUser: DB is null")
		log.Fatal(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tsql := "INSERT INTO dbo.Accounts (login_name, password, ip) VALUES (@Username, @Password, @IP);"
	stmt, err := db.Prepare(tsql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(
		ctx,
		sql.Named("Username", user),
		sql.Named("Password", pass),
		sql.Named("IP", ip))
}

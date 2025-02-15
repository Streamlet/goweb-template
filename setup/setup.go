package setup

import (
	"fmt"
	"goweb/setup/tables"
	"path/filepath"
	"strings"

	"github.com/Streamlet/gosql"
	"golang.org/x/term"
)

func InteractiveSetup(c *gosql.Connection, db string) {
	println("Initialize database...")
	var adminName, nickName, password, confirmPassword string
	print("Please enter admin name: ")
	fmt.Scanln(&adminName)
	print("Please enter nick name: ")
	fmt.Scanln(&nickName)
	for password == "" || confirmPassword == "" {
		print("Please enter admin password: ")
		input, err := term.ReadPassword(0)
		println("")
		if err != nil {
			println(err.Error())
			return
		}
		password = string(input)

		print("Please confirm admin password: ")
		input, err = term.ReadPassword(0)
		println("")
		if err != nil {
			println(err.Error())
			return
		}
		confirmPassword = string(input)
		if confirmPassword != password {
			fmt.Println("Passwords not match.")
			password = ""
			confirmPassword = ""
		}
	}
	if err := install(c, db, adminName, nickName, password); err != nil {
		println(err.Error())
		return
	}
	println("Setup accomplished. Please restart service without '--setup' flag.")
}

func install(c *gosql.Connection, db, adminName, nickName, password string) error {
	if _, err := c.Update("CREATE DATABASE IF NOT EXISTS " + db); err != nil {
		return err
	}
	if _, err := c.Update("USE " + db); err != nil {
		return err
	}
	if err := createTables(c); err != nil {
		return err
	}
	return nil
}

func createTables(c *gosql.Connection) error {
	files, err := tables.CreateTableSqls.ReadDir(".")
	if err != nil {
		return err
	}
	for _, file := range files {
		table := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		_, err = c.Update("DROP TABLE IF EXISTS " + table)
		if err != nil {
			return err
		}
		sql, err := tables.CreateTableSqls.ReadFile(file.Name())
		if err != nil {
			return err
		}
		_, err = c.Update(string(sql))
		if err != nil {
			return err
		}
	}
	return nil
}

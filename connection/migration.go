package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	databaseUrl := "postgres://postgres:admin@localhost:5432/project-web"

	var err error

	//context akan buat backgroundnya jalan terus , kalo gapake context nanti sekali eksekusi doang
	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}

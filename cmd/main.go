package main

import (
	"database/sql"
	"gosql/cmd/server/app"
	"gosql/pkg/customers"
	"log"
	"net"
	"net/http"
)

package main

import (
"database/sql"
_ "github.com/jackc/pgx/v4/stdlib"
"log"
"net"
"net/http"
"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func execute(host string, port string, dsn string) (err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()

	mux := http.NewServeMux()
	customerSvg := customers.NewService(db)
	server := app.NewServer(mux, customerSvg)
	server.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}


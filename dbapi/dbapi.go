/******************************************************************************
dbapi provides an HTTP based universal api that can be used to access the
realtime database as well as supported relational databases and nosql
databases. Its primary use is for the HMI, but can also be used with command
line tools that utilize an HTTP client.

All data exchanges utilize JSON.

the dbapi is accessed through the RTAP http server
******************************************************************************/

package dbapi

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DBapi provides the universal interface.
type DBapi interface {
	Init(conn_str string) error
	Get(schema_name string, table_name string, columns *[]string,
		where string, orderby *[]string, out *[]any) error
}

var scada dbapi_rt
var hmi dbapi_rt
var scada_config dbapi_sql
var hmi_config dbapi_sql

var DBs []DBapi

func init() {

	scada.Init("scada")

	DBs = append(DBs, &scada)
	DBs = append(DBs, &hmi)
	DBs = append(DBs, &scada_config)
	DBs = append(DBs, &hmi_config)

	/*
		var out []any
			for i, v := range DBs {

				err := v.Get("", "", nil, "", nil, &out)
				if err == nil {
					fmt.Println("Error getting row")
				} else {
					fmt.Printf("Values: %d %s\n", i, v)
				}
			}
	*/

}

func sendError(w http.ResponseWriter, req *http.Request, status_code int, message string) {
	// Get the current time to use in the headers.
	t := time.Now()

	// Set the headers
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Date", fmt.Sprintf(t.Format("Mon, 2 Jan 2006 15:04:05 MST"))) //TODO: verify this with wireshark
	w.WriteHeader(status_code)
	w.Write([]byte(message))

}

func DBapiHandler(w http.ResponseWriter, req *http.Request) {

	// Parse the dbid, schema, and table. the dbid uses a map to a db_desciptor
	// structure that holds all the meta data required to access the database.
	path := strings.TrimRight(req.URL.Path, "/") //remove trailing '/', if any
	splits := strings.Split(path, "/")
	if len(splits) != 5 {
		sendError(w, req, 400, "Invalid path. Must be in the form: '/dbpai/[db id]/[schema name]/[table_name]'")
		return
	}
	dbid := splits[2]
	schema := splits[3]
	table := splits[4]

	// Get the current time to use in the headers.
	t := time.Now()

	// Set the headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Date", fmt.Sprint(t.Format("Mon, 2 Jan 2006 15:04:05 MST")))

	fmt.Fprintf(w, "{\n")

	fmt.Fprintf(w, "\"request\":{")
	fmt.Fprintf(w, "\"RemoteAddr\":\"%s\",", req.RemoteAddr)
	fmt.Fprintf(w, "\"Path\":\"%s\",", req.URL.Path)

	fmt.Fprintf(w, "\"splits\":\"%d\",", len(splits))

	fmt.Fprintf(w, "\"Method\":\"%s\",", req.Method)
	fmt.Fprintf(w, "\"Query\":\"%s\",", req.URL.Query())
	fmt.Fprintf(w, "\"dbid\":\"%s\",", dbid)
	fmt.Fprintf(w, "\"schema\":\"%s\",", schema)
	fmt.Fprintf(w, "\"table\":\"%s\"", table)
	fmt.Fprintf(w, "}\n")

	fmt.Println(req.URL.RawQuery)
	fmt.Println(req.URL.Query())

	fmt.Fprintf(w, "}")

}

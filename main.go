//package main
//
//import (
//	_ "admin/models"
//	_ "context"
//	"database/sql"
//	"fmt"
//	_ "github.com/jackc/pgx/v4/pgxpool"
//	"html/template"
//	"log"
//	"net/http"
//	_ "reflect"
//	"strconv"
//
//	_ "github.com/lib/pq" // PostgreSQL драйвер
//)
//
//type PageData struct {
//	CurrentPage int
//	TotalPages  int
//	Records     any
//	TableName   string
//	Columns     []string // Добавлено поле для колонок
//}
//
//const (
//	recordsPerPage = 5
//)
//
//var db *sql.DB
//
//// Функция для подключения к базе данных
//func connectToDatabase() {
//	var err error
//	dsn := "postgres://postgres:radpass@localhost:5432/freeradius?sslmode=disable"
//	db, err = sql.Open("postgres", dsn)
//	if err != nil {
//		log.Fatalf("Error connecting to the database: %v\n", err)
//	}
//	if err := db.Ping(); err != nil {
//		log.Fatalf("Database is not reachable: %v\n", err)
//	}
//}
//
//// Основная функция
//func main() {
//	connectToDatabase()
//	defer db.Close()
//
//	// Создание пользовательских функций
//	funcMap := template.FuncMap{
//		"sub": func(a, b int) int { return a - b },
//		"add": func(a, b int) int { return a + b },
//	}
//
//	// Предзагрузка шаблонов
//	indexTmpl := template.Must(template.New("index").Parse(indexTemplate))
//	tableTmpl := template.Must(template.New("table").Funcs(funcMap).Parse(tableTemplate))
//	addTmpl := template.Must(template.New("add").Parse(addTemplate))
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		tables := []string{"radacct", "radcheck", "radgroupcheck", "radgroupreply", "radreply", "radusergroup", "radpostauth", "nas", "nasreload"}
//		if err := indexTmpl.Execute(w, tables); err != nil {
//			http.Error(w, "Template execution error", http.StatusInternalServerError)
//		}
//	})
//
//	http.HandleFunc("/table/", func(w http.ResponseWriter, r *http.Request) {
//		table := r.URL.Path[len("/table/"):]
//
//		// Подсчет общего количества записей
//		var totalRecords int
//		err := db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)).Scan(&totalRecords)
//		if err != nil {
//			http.Error(w, "Error counting records", http.StatusInternalServerError)
//			return
//		}
//
//		// Подсчет количества страниц
//		totalPages := (totalRecords + recordsPerPage - 1) / recordsPerPage
//		page := 1
//		if p := r.URL.Query().Get("page"); p != "" {
//			if pInt, err := strconv.Atoi(p); err == nil {
//				page = pInt
//			}
//		}
//
//		if page < 1 {
//			page = 1
//		} else if page > totalPages {
//			page = totalPages
//		}
//
//		// Получение данных
//		offset := (page - 1) * recordsPerPage
//		rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %s LIMIT $1 OFFSET $2`, table), recordsPerPage, offset)
//		if err != nil {
//			http.Error(w, "Database query error", http.StatusInternalServerError)
//			return
//		}
//		defer rows.Close()
//
//		// Чтение данных
//		var records []map[string]any
//		cols, _ := rows.Columns()
//		for rows.Next() {
//			values := make([]any, len(cols))
//			pointers := make([]any, len(cols))
//			for i := range values {
//				pointers[i] = &values[i]
//			}
//			if err := rows.Scan(pointers...); err != nil {
//				http.Error(w, "Error reading record", http.StatusInternalServerError)
//				return
//			}
//
//			record := make(map[string]any)
//			for i, col := range cols {
//				record[col] = values[i]
//			}
//			records = append(records, record)
//		}
//
//		// Отображение данных
//		data := PageData{
//			CurrentPage: page,
//			TotalPages:  totalPages,
//			Records:     records,
//			TableName:   table,
//		}
//		if err := tableTmpl.Execute(w, data); err != nil {
//			http.Error(w, "Template execution error", http.StatusInternalServerError)
//		}
//	})
//
//	http.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
//		table := r.URL.Path[len("/add/"):]
//
//		if r.Method == http.MethodGet {
//			// Получаем список всех столбцов таблицы
//			rows, err := db.Query(fmt.Sprintf(`SELECT column_name FROM information_schema.columns WHERE table_name = '%s'`, table))
//			if err != nil {
//				http.Error(w, "Error retrieving column names", http.StatusInternalServerError)
//				return
//			}
//			defer rows.Close()
//
//			var columns []string
//			for rows.Next() {
//				var column string
//				if err := rows.Scan(&column); err != nil {
//					http.Error(w, "Error reading column names", http.StatusInternalServerError)
//					return
//				}
//				columns = append(columns, column)
//			}
//
//			// Отображение формы для добавления записи
//			data := PageData{
//				TableName: table,
//				Columns:   columns,
//			}
//			if err := addTmpl.Execute(w, data); err != nil {
//				http.Error(w, "Template execution error", http.StatusInternalServerError)
//			}
//		} else if r.Method == http.MethodPost {
//			r.ParseForm()
//
//			columns := []string{}
//			values := []any{}
//			placeholders := []string{}
//			for key, value := range r.Form {
//				columns = append(columns, key)
//				values = append(values, value[0])
//				placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
//			}
//
//			// Выполнение вставки записи
//			query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, join(columns, ","), join(placeholders, ","))
//			_, err := db.Exec(query, values...)
//			if err != nil {
//				http.Error(w, fmt.Sprintf("Error inserting record: %v", err), http.StatusInternalServerError)
//				return
//			}
//			http.Redirect(w, r, fmt.Sprintf("/table/%s", table), http.StatusSeeOther)
//		}
//	})
//
//	port := ":8080"
//	fmt.Printf("Server is running at http://localhost%s\n", port)
//	if err := http.ListenAndServe(port, nil); err != nil {
//		log.Fatalf("Server error: %v\n", err)
//	}
//}
//
//// Вспомогательная функция для объединения строк
//func join(strings []string, sep string) string {
//	result := ""
//	for i, s := range strings {
//		if i > 0 {
//			result += sep
//		}
//		result += s
//	}
//	return result
//}
//
//// Шаблон для главной страницы
//const indexTemplate = `<!DOCTYPE html>
//<html>
//<head>
// <title>Tables</title>
//</head>
//<body>
// <h1>Tables</h1>
// <ul>
//     {{range .}}
//     <li><a href="/table/{{.}}">{{.}}</a> | <a href="/add/{{.}}">Add Record</a></li>
//     {{end}}
// </ul>
//</body>
//</html>
//`
//
//// Шаблон для страницы таблицы
//const tableTemplate = `<!DOCTYPE html>
//<html>
//<head>
//  <title>{{.TableName}}</title>
//</head>
//<body>
//  <h1>Table: {{.TableName}}</h1>
//  <table border="1">
//      <thead>
//          <tr>
//              {{range $key, $value := (index .Records 0)}}
//              <th>{{$key}}</th>
//              {{end}}
//          </tr>
//      </thead>
//      <tbody>
//          {{range .Records}}
//          <tr>
//              {{range $key, $value := .}}
//              <td>{{$value}}</td>
//              {{end}}
//          </tr>
//          {{end}}
//      </tbody>
//  </table>
//  <div>
//      {{if gt .CurrentPage 1}}
//      <a href="?page={{sub .CurrentPage 1}}">Previous</a>
//      {{end}}
//      Page {{.CurrentPage}} of {{.TotalPages}}
//      {{if lt .CurrentPage .TotalPages}}
//      <a href="?page={{add .CurrentPage 1}}">Next</a>
//      {{end}}
//  </div>
//</body>
//</html>
//`
//
//// Шаблон для добавления записи
//const addTemplate = `<!DOCTYPE html>
//<html>
//<head>
//  <title>Add Record</title>
//</head>
//<body>
//  <h1>Add Record to {{.TableName}}</h1>
//  <form method="POST">
//      {{range .Columns}}
//      <p><label>{{.}}: <input type="text" name="{{.}}"></label></p>
//      {{end}}
//      <p><button type="submit">Add</button></p>
//  </form>
//</body>
//</html>
//`
//
////
////package main
////
////import (
////	"context"
////	"fmt"
////	"html/template"
////	"log"
////	"net/http"
////	"reflect"
////	_ "time"
////
////	"admin/models"
////	"github.com/jackc/pgx/v4/pgxpool"
////)
////
////const (
////	dbURL = "postgres://postgres:radpass@localhost:5432/freeradius?sslmode=disable"
////)
////
////var conn *pgxpool.Pool
////
////func main() {
////	var err error
////	conn, err = pgxpool.Connect(context.Background(), dbURL)
////	if err != nil {
////		log.Fatalf("Unable to connect to database: %v\n", err)
////	}
////	defer conn.Close()
////
////	http.HandleFunc("/radacct", handleTable("radacct", func() interface{} { return &models.Radacct{} }, "templates/radacct.html"))
////	http.HandleFunc("/radcheck", handleTable("radcheck", func() interface{} { return &models.Radcheck{} }, "templates/radcheck.html"))
////	http.HandleFunc("/radgroupcheck", handleTable("radgroupcheck", func() interface{} { return &models.Radgroupcheck{} }, "templates/radgroupcheck.html"))
////	http.HandleFunc("/radgroupreply", handleTable("radgroupreply", func() interface{} { return &models.Radgroupreply{} }, "templates/radgroupreply.html"))
////	http.HandleFunc("/radreply", handleTable("radreply", func() interface{} { return &models.Radreply{} }, "templates/radreply.html"))
////	http.HandleFunc("/radusergroup", handleTable("radusergroup", func() interface{} { return &models.Radusergroup{} }, "templates/radusergroup.html"))
////	http.HandleFunc("/radpostauth", handleTable("radpostauth", func() interface{} { return &models.Radpostauth{} }, "templates/radpostauth.html"))
////	http.HandleFunc("/nas", handleTable("nas", func() interface{} { return &models.Nas{} }, "templates/nas.html"))
////	http.HandleFunc("/nasreload", handleTable("nasreload", func() interface{} { return &models.Nasreload{} }, "templates/nasreload.html"))
////
////	fmt.Println("Server started at :8080")
////	log.Fatal(http.ListenAndServe(":8080", nil))
////}
////
////func handleTable(tableName string, modelFunc func() interface{}, tmpl string) http.HandlerFunc {
////	return func(w http.ResponseWriter, r *http.Request) {
////		query := fmt.Sprintf("SELECT * FROM %s", tableName)
////		rows, err := conn.Query(context.Background(), query)
////		if err != nil {
////			http.Error(w, "Unable to fetch records", http.StatusInternalServerError)
////			log.Printf("Unable to fetch records from %s: %v\n", tableName, err)
////			return
////		}
////		defer rows.Close()
////
////		records := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(modelFunc()).Elem()), 0, 0).Interface()
////		for rows.Next() {
////			elem := modelFunc()
////			values := reflect.ValueOf(elem).Elem()
////
////			fields := make([]interface{}, values.NumField())
////			for i := 0; i < values.NumField(); i++ {
////				if values.Field(i).Addr().Interface() != nil {
////					fields[i] = values.Field(i).Addr().Interface()
////				}
////			}
////
////			if err := rows.Scan(fields...); err != nil {
////				http.Error(w, "Unable to read record", http.StatusInternalServerError)
////				log.Printf("Unable to read record for table %s: %v\n", tableName, err)
////				return
////			}
////
////			records = reflect.Append(reflect.ValueOf(records), reflect.ValueOf(elem).Elem()).Interface()
////		}
////
////		if err := renderTemplate(w, tmpl, records); err != nil {
////			http.Error(w, "Unable to render template", http.StatusInternalServerError)
////			log.Printf("Unable to render template for table %s: %v\n", tableName, err)
////		}
////	}
////}
////
////func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
////	t, err := template.ParseFiles(tmpl)
////	if err != nil {
////		return err
////	}
////	return t.Execute(w, data)
////}

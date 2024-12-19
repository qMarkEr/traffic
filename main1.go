//package main
//
//import (
//	"admin/models" // Замените на ваш путь импорта
//	"database/sql"
//	"fmt"
//	_ "github.com/lib/pq" // Подключаем драйвер PostgreSQL
//	"html/template"
//	"log"
//	"net/http"
//	"reflect"
//)
//
//var db *sql.DB
//
//// Универсальная функция для чтения данных из базы данных и их заполняющий слайс
//func fetchData(query string, dest interface{}) error {
//	rows, err := db.Query(query)
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//
//	// Получаем слайс, в который будем записывать данные
//	val := reflect.ValueOf(dest)
//	if val.Kind() != reflect.Ptr || val.IsNil() {
//		return fmt.Errorf("destination must be a non-nil pointer")
//	}
//
//	// Получаем слайс из интерфейса
//	sliceVal := val.Elem()
//	if sliceVal.Kind() != reflect.Slice {
//		return fmt.Errorf("destination must be a slice")
//	}
//
//	// Получаем тип элемента слайса
//	elemType := sliceVal.Type().Elem()
//
//	// Прокачиваем через каждую строку
//	for rows.Next() {
//		// Создаем новый элемент для структуры
//		elemPtr := reflect.New(elemType).Interface()
//
//		// Заполняем его с помощью scan
//		// Нужно сделать передачу структуры через интерфейс
//		err := rows.Scan(getStructFields(elemPtr)...)
//		if err != nil {
//			return err
//		}
//
//		// Добавляем в слайс
//		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(elemPtr).Elem()))
//	}
//
//	if err := rows.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Вспомогательная функция для получения полей структуры для метода Scan
//func getStructFields(i interface{}) []interface{} {
//	val := reflect.ValueOf(i).Elem()
//	fields := make([]interface{}, val.NumField())
//
//	for i := 0; i < val.NumField(); i++ {
//		fields[i] = val.Field(i).Addr().Interface()
//	}
//
//	return fields
//}
//
//func init() {
//	// Подключение к базе данных PostgreSQL
//	var err error
//	connStr := "postgres://postgres:radpass@localhost:5432/freeradius?sslmode=disable"
//	db, err = sql.Open("postgres", connStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func main() {
//	http.HandleFunc("/", showTables)  // Обработчик для главной страницы
//	http.ListenAndServe(":8080", nil) // Запуск HTTP сервера
//}
//
//func showTables(w http.ResponseWriter, r *http.Request) {
//	// Запрос для получения данных из таблицы radacct
//	query := "SELECT RadAcctId, AcctSessionId, AcctUniqueId, UserName, Realm, NASIPAddress, AcctStartTime, AcctStopTime, AcctInputOctets, AcctOutputOctets FROM radacct LIMIT 10"
//
//	// Создаем слайс для хранения данных из базы данных
//	var radacctList []models.Radacct
//
//	// Вызов универсальной функции для получения данных
//	err := fetchData(query, &radacctList)
//	if err != nil {
//		http.Error(w, "Error fetching data from database", http.StatusInternalServerError)
//		log.Println(err)
//		return
//	}
//
//	// Создание шаблона HTML
//	tmpl := template.Must(template.New("tables").Parse(`
//		<!DOCTYPE html>
//		<html lang="en">
//		<head>
//			<meta charset="UTF-8">
//			<meta name="viewport" content="width=device-width, initial-scale=1.0">
//			<title>FreeRADIUS Tables</title>
//		</head>
//		<body>
//			<h1>Radacct Table Data</h1>
//			<table border="1">
//				<tr>
//					<th>RadAcctId</th>
//					<th>AcctSessionId</th>
//					<th>AcctUniqueId</th>
//					<th>UserName</th>
//					<th>Realm</th>
//					<th>NASIPAddress</th>
//					<th>AcctStartTime</th>
//					<th>AcctStopTime</th>
//					<th>AcctInputOctets</th>
//					<th>AcctOutputOctets</th>
//				</tr>
//				{{range .}}
//				<tr>
//					<td>{{.RadAcctId}}</td>
//					<td>{{.AcctSessionId}}</td>
//					<td>{{.AcctUniqueId}}</td>
//					<td>{{.UserName}}</td>
//					<td>{{.Realm}}</td>
//					<td>{{.NASIPAddress}}</td>
//					<td>{{.AcctStartTime}}</td>
//					<td>{{.AcctStopTime}}</td>
//					<td>{{.AcctInputOctets}}</td>
//					<td>{{.AcctOutputOctets}}</td>
//				</tr>
//				{{end}}
//			</table>
//		</body>
//		</html>
//	`))
//
//	// Отправка данных в шаблон
//	if err := tmpl.Execute(w, radacctList); err != nil {
//		http.Error(w, "Error rendering template", http.StatusInternalServerError)
//		log.Println(err)
//	}
//}

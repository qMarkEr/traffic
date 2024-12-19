//package main
//
//import (
//	"fmt"
//	"html/template"
//	"log"
//	"net/http"
//)
//
//// Структура для данных, передаваемых в шаблон
//type PageData struct {
//	Records [][]string // Двумерный массив для хранения данных
//}
//
//func handleTable(w http.ResponseWriter, r *http.Request) {
//	// Пример двумерного массива
//	records := [][]string{
//		{"ID", "Name", "Age"},
//		{"1", "Alice", "30"},
//		{"2", "Bob", "25"},
//		{"3", "Charlie", "35"},
//	}
//
//	for _, row := range records {
//		// Перебор элементов в строках
//		for _, value := range row {
//			// Выводим каждый элемент
//			fmt.Printf("%s\t", value)
//		}
//		// Печатаем новую строку после каждой строки массива
//		fmt.Println()
//	}
//
//	// Передача данных в шаблон
//	data := PageData{
//		Records: records,
//	}
//
//	// Создание и выполнение шаблона
//	tmpl, err := template.New("table").Parse(`
//	<!DOCTYPE html>
//	<html lang="ru">
//	<head>
//		<meta charset="UTF-8">
//		<meta name="viewport" content="width=device-width, initial-scale=1.0">
//		<title>Двумерная таблица</title>
//		<style>
//			table {
//				width: 100%;
//				border-collapse: collapse;
//			}
//			th, td {
//				padding: 10px;
//				border: 1px solid #ddd;
//			}
//			th {
//				background-color: #4CAF50;
//				color: white;
//			}
//			tr:nth-child(even) {
//				background-color: #f2f2f2;
//			}
//		</style>
//	</head>
//	<body>
//		<h1>Таблица с данными</h1>
//		<table>
//			<tr>
//				{{range .Records}}
//  			<tr>
//					{{range .}}
//					<td>{{.}}</td>
//					{{end}}
//			</tr>
//{{end}}
//		</table>
//	</body>
//	</html>
//	`)
//	if err != nil {
//		http.Error(w, "Ошибка при создании шаблона", http.StatusInternalServerError)
//		return
//	}
//
//	// Выполнение шаблона и передача данных
//	err = tmpl.Execute(w, data)
//	if err != nil {
//		http.Error(w, "Ошибка при рендеринге шаблона", http.StatusInternalServerError)
//	}
//}
//
//func main() {
//	http.HandleFunc("/table", handleTable)
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

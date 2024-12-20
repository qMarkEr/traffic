package main

import (
	_ "admin/models"
	"bufio"
	_ "context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	_ "encoding/hex"
	"fmt"
	_ "github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	_ "reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	recordsPerPage = 5
)

var (
	sessions     = make(map[string]string) // Карта сессий: sessionID -> username
	sessionsLock sync.Mutex
)

type PageData struct {
	CurrentPage int
	TotalPages  int
	Records     any
	TableName   string
	Columns     []string // Добавлено поле для колонок
}

var db *sql.DB

const (
	indexPath = "templates/index.html"
	addPath   = "templates/add.html"
	tablePath = "templates/table.html"
)

// Функция для подключения к базе данных
func connectToDatabase() {
	var err error
	dsn := "postgres://postgres:radpass@localhost:5432/freeradius?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v\n", err)
	}
}

// Функция для отображения главной страницы
func handleIndex(w http.ResponseWriter, r *http.Request) {
	tables := []string{"radacct", "radcheck", "radgroupcheck", "radgroupreply", "radreply", "radusergroup", "radpostauth", "nas", "nasreload"}
	indexTmpl := template.Must(template.New("index.html").ParseFiles(indexPath))
	if err := indexTmpl.Execute(w, tables); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// Функция для отображения таблицы
func handleTable(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Path[len("/table/"):]

	// Подсчет общего количества записей
	var totalRecords int
	err := db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)).Scan(&totalRecords)
	if err != nil {
		http.Error(w, "Error counting records", http.StatusInternalServerError)
		return
	}

	// Подсчет количества страниц
	totalPages := (totalRecords + recordsPerPage - 1) / recordsPerPage
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil {
			page = pInt
		}
	}

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	// Получение данных
	offset := (page - 1) * recordsPerPage
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %s LIMIT $1 OFFSET $2`, table), recordsPerPage, offset)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Чтение данных
	var records [][]any
	cols, _ := rows.Columns()
	for rows.Next() {
		values := make([]any, len(cols))
		pointers := make([]any, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			http.Error(w, "Error reading record", http.StatusInternalServerError)
			return
		}

		record := make([]any, len(cols))
		for i, _ := range cols {
			var value any
			switch v := values[i].(type) {
			case []byte:
				value = string(v) // Если значение — []byte, преобразуем в строку
			case nil:
				value = "NULL" // Если значение NULL, заменяем его на строку "NULL"
			default:
				value = v // В остальных случаях оставляем как есть
			}
			record[i] = value
			//println(col, values[i])
		}
		records = append(records, record)
	}
	//for i := range records {
	//	for j, col := range cols {
	//		// Преобразование значения в читабельный формат
	//		var value any
	//		switch v := records[i][j].(type) {
	//		case []byte:
	//			value = string(v) // Если значение — []byte, преобразуем в строку
	//		case nil:
	//			value = "NULL" // Если значение NULL, заменяем его на строку "NULL"
	//		default:
	//			value = v // В остальных случаях оставляем как есть
	//		}
	//		var value1 any
	//		switch v := records1[i][col].(type) {
	//		case []byte:
	//			value1 = string(v) // Если значение — []byte, преобразуем в строку
	//		case nil:
	//			value1 = "NULL" // Если значение NULL, заменяем его на строку "NULL"
	//		default:
	//			value1 = v // В остальных случаях оставляем как есть
	//		}
	//
	//		// Логирование для отладки
	//		//fmt.Printf("Val: %s, Val1: %s\n", value, value1)
	//	}
	//}
	//for _, row := range records {
	//	// Перебор элементов в строках
	//	for _, value := range row {
	//		// Выводим каждый элемент
	//		fmt.Printf("%s\t", value)
	//	}
	//	// Печатаем новую строку после каждой строки массива
	//	fmt.Println()
	//}
	//fmt.Println("---------------------------------------")
	//

	data := PageData{
		CurrentPage: page,
		TotalPages:  totalPages,
		Records:     records,
		TableName:   table,
		Columns:     cols,
	}

	tableTmpl := template.Must(template.New("table.html").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}).ParseFiles(tablePath))
	if err := tableTmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// Функция для отображения формы добавления записи
func handleAdd(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Path[len("/add/"):]

	if r.Method == http.MethodGet {
		// Получаем список всех столбцов таблицы
		rows, err := db.Query(fmt.Sprintf(`SELECT column_name FROM information_schema.columns WHERE table_name = '%s'`, table))
		if err != nil {
			http.Error(w, "Error retrieving column names", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var columns []string
		for rows.Next() {
			var column string
			if err := rows.Scan(&column); err != nil {
				http.Error(w, "Error reading column names", http.StatusInternalServerError)
				return
			}
			columns = append(columns, column)
		}

		// Отображение формы для добавления записи
		data := PageData{
			TableName: table,
			Columns:   columns,
		}
		addTmpl := template.Must(template.New("add.html").ParseFiles(addPath))
		if err := addTmpl.Execute(w, data); err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		columns := []string{}
		values := []any{}
		placeholders := []string{}
		for key, value := range r.Form {
			columns = append(columns, key)
			values = append(values, value[0])
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
		}

		// Выполнение вставки записи
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, join(columns, ","), join(placeholders, ","))
		_, err := db.Exec(query, values...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting record: %v", err), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/table/%s", table), http.StatusSeeOther)
	}
}

// / Функция для обработки удаления записи
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	table := r.URL.Query().Get("table")
	id := r.URL.Query().Get("id")

	// Логирование параметров для диагностики
	fmt.Printf("Table: %s, ID: %s\n", table, id)

	if table == "" || id == "" {
		http.Error(w, "Missing table or id parameter", http.StatusBadRequest)
		return
	}
	if table == "radacct" {
		// Формирование запроса для удаления записи
		query := fmt.Sprintf("DELETE FROM %s WHERE radacctid = $1", table)
		_, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting record: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// Формирование запроса для удаления записи
		query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
		_, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting record: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Перенаправление обратно на страницу таблицы
	http.Redirect(w, r, fmt.Sprintf("/table/%s", table), http.StatusSeeOther)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	const loginTemplatePath = "templates/login.html"

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("login.html").ParseFiles(loginTemplatePath))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Пример проверки логина и пароля
		if username == "admin" && password == "admin" {
			sessionID := generateSessionID()

			// Сохраняем сессию
			sessionsLock.Lock()
			sessions[sessionID] = username
			sessionsLock.Unlock()

			// Устанавливаем куки с идентификатором сессии
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	sessionID := cookie.Value

	sessionsLock.Lock()
	_, exists := sessions[sessionID]
	sessionsLock.Unlock()

	return exists
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) && r.URL.Path != "/login" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate session ID: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		sessionID := cookie.Value

		// Удаляем сессию
		sessionsLock.Lock()
		delete(sessions, sessionID)
		sessionsLock.Unlock()

		// Удаляем куки
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0), // Устанавливаем истекшее время
			HttpOnly: true,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Вспомогательная функция для объединения строк
func join(strings []string, sep string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp string
	IP        string
	Service   string
	Level     string
	Message   string
}

// TemplateData holds data for rendering templates
type TemplateData struct {
	LogFiles     []string
	LogEntries   []LogEntry
	SelectedFile string
}

// parseLogFile parses the log file at the given path
func parseLogFile(filePath string) ([]LogEntry, error) {
	var entries []LogEntry

	// Regular expression to parse log lines
	logPattern := regexp.MustCompile(`^(?P<timestamp>[\d\-T:+]+)\s+(?P<ip>\d{1,3}(\.\d{1,3}){3})\s+\[(?P<service>[a-zA-Z0-9]+)\.(?P<level>[a-z]+)\]\s+(?P<message>.+)$`)

	// Open the log file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := logPattern.FindStringSubmatch(line)
		if match != nil {
			entries = append(entries, LogEntry{
				Timestamp: match[1],
				IP:        match[2],
				Service:   match[3],
				Level:     match[4],
				Message:   match[5],
			})
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return entries, nil
}

// getLogFiles retrieves a list of log files in the given directory
func getLogFiles(logDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log") {
			files = append(files, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error reading log directory: %w", err)
	}

	return files, nil
}

// logListHandler serves the list of log files
func logListHandler(w http.ResponseWriter, r *http.Request) {
	logDir := "/var/log/syslog/" // Change to your log directory
	files, err := getLogFiles(logDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve log files: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse and execute the HTML template for log selection
	tmpl, err := template.ParseFiles("templates/log_template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load template: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare template data with the list of log files
	data := TemplateData{
		LogFiles: files,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
	}
}

// logViewerHandler serves the content of a selected log file
func logViewerHandler(w http.ResponseWriter, r *http.Request) {
	logDir := "/var/log/syslog/" // Change to your log directory
	logFile := r.URL.Query().Get("file")

	if logFile == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fullPath := filepath.Join(logDir, logFile)
	entries, err := parseLogFile(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse log file: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse and execute the HTML template for log viewer
	tmpl, err := template.ParseFiles("templates/log_template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load template: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare template data with the log entries
	data := TemplateData{
		LogEntries:   entries,
		SelectedFile: logFile,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	connectToDatabase()
	defer db.Close()

	// Подключаем обработчики
	http.HandleFunc("/logs/", requireAuth(logListHandler))
	http.HandleFunc("/logs/view/", requireAuth(logViewerHandler))
	http.HandleFunc("/", requireAuth(handleIndex))
	http.HandleFunc("/table/", requireAuth(handleTable))
	http.HandleFunc("/delete/", requireAuth(handleDelete))
	http.HandleFunc("/add/", requireAuth(handleAdd))
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", requireAuth(handleLogout))

	port := ":8080"
	fmt.Printf("Server is running at http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

// radtest shs 1111 localhost 0 testing123
// echo "User-Name=shs, User-Password=1111" | radclient -x localhost:1812 auth testing123
// echo "User-Name=shs, Acct-Status-Type=Start" | radclient -x localhost acct testing123
// echo "User-Name = 'shs', Acct-Status-Type = Start, NAS-IP-Address = 192.168.1.1, Acct-Session-Id = 'session123'" | radclient -x localhost acct testing123
// http://172.22.87.209:8080/login

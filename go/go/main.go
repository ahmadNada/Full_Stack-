package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"	
)

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://ahmad:postgres@localhost/checksApp?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully Connected")
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

type Check struct {
	ID   string
	Amount  string
	BankName string
	DueDate  string
	Status string
}

func main() {
	http.HandleFunc("/", showChecks)
	http.HandleFunc("/add", addCheck)
	http.HandleFunc("/pay", payCheck)
	http.HandleFunc("/filter", filterChecks)
	http.ListenAndServe(":8080", nil)
}


func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}


func showChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM checks")
	
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()
	chks := make([]Check, 0)
	for rows.Next() {
		chk := Check{}
		err := rows.Scan(&chk.ID, &chk.Amount, &chk.BankName, &chk.DueDate, &chk.Status)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		chks = append(chks, chk)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "index.gohtml", chks)
}


func addCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	chk := Check{}
	chk.ID = r.FormValue("id")
	chk.Amount = r.FormValue("amount")
	chk.BankName = r.FormValue("bankName")
	chk.DueDate = r.FormValue("dueDate")
	if chk.ID == "" || chk.Amount == "" || chk.BankName == "" || chk.DueDate == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	var err error
	_, err = db.Exec("INSERT INTO checks (ID, amount, bankName, dueDate,status) VALUES ($1, $2, $3, $4,'Not Paid')", chk.ID, chk.Amount, chk.BankName, chk.DueDate)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	  http.Redirect(w, r, "http://localhost:8080/", 301)
}


func filterChecks(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	if r.FormValue("filter")=="Not Paid"{
	rows, err := db.Query("SELECT * FROM checks WHERE status='Not Paid'")
	defer rows.Close()
	chks := make([]Check, 0)
	for rows.Next() {
		chk := Check{}
		err := rows.Scan(&chk.ID, &chk.Amount, &chk.BankName, &chk.DueDate, &chk.Status)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		chks = append(chks, chk)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpl.ExecuteTemplate(w, "index.gohtml", chks)}
	if r.FormValue("filter")=="Paid"{
		rows, err := db.Query("SELECT * FROM checks WHERE status='Paid'")
		defer rows.Close()
		chks := make([]Check, 0)
		for rows.Next() {
			chk := Check{}
			err := rows.Scan(&chk.ID, &chk.Amount, &chk.BankName, &chk.DueDate, &chk.Status) // order matters
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
				}
			chks = append(chks, chk)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
			}
			tpl.ExecuteTemplate(w, "index.gohtml", chks)
			}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if r.FormValue("filter")=="showAll"{
		http.Redirect(w, r, "http://localhost:8080/", 301)
		}
}


func payCheck(w http.ResponseWriter, r *http.Request){
	var err error
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE checks SET status = 'Paid' WHERE id="+id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "http://localhost:8080/", 301)
	}




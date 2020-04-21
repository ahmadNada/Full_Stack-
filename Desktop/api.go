package main

import(
"encoding/json"
"log"
"net/http"
"github.com/gorilla/mux"
)


type Check struct{
ID string `json:"id"`
Bank string `json:"bank"`
Amount string `json:"amount"`
Status string `json:"status"`
}


var checks []Check

func getChecks(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
json.NewEncoder(w).Encode(checks)
}


func getPaid(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
var paid []Check
for _, item := range checks {
if item.Status == "Paid"{
paid = append(paid , item)
json.NewEncoder(w).Encode(paid)
return
}

}
}


func getNotPaid(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
var notPaid []Check
for _, item := range checks {
if item.Status == "Not Paid"{
notPaid = append(notPaid , item)
json.NewEncoder(w).Encode(notPaid)
return 
}

}
}


func createCheck(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
var check Check
_ = json.NewDecoder(r.Body).Decode(&check)
checks = append (checks , check)
json.NewEncoder(w).Encode(check)
}


func deleteCheck(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
params := mux.Vars(r)
for index , item := range checks {
if item.ID == params["id"]{
checks = append(checks[:index],checks[index+1:]...)
break
}
}
json.NewEncoder(w).Encode(checks)
}


func payCheck(w http.ResponseWriter,r *http.Request){
w.Header().Set("content-type","application/json")
params:=mux.Vars(r)
for index, item := range checks {
if item.ID == params["id"]{
checks[index].Status="Paid"
json.NewEncoder(w).Encode(checks)
return
}
}}

func main(){
r := mux.NewRouter()
checks= append(checks, Check{ID: "1",Bank:"FNB",Amount:"342534",Status: "Not Paid"})
checks= append(checks, Check{ID: "2",Bank:"ABC",Amount:"243",Status: "Paid"})
checks= append(checks, Check{ID: "3",Bank:"XYZ",Amount:"542",Status: "Paid"})

r.HandleFunc("/api/checks",getChecks).Methods("GET")
r.HandleFunc("/api/checks/paid",getPaid).Methods("GET")
r.HandleFunc("/api/checks/notPaid",getNotPaid).Methods("GET")
r.HandleFunc("/api/checks",createCheck).Methods("POst")
r.HandleFunc("/api/checks/{id}",deleteCheck).Methods("Delete")
r.HandleFunc("/api/checks/{id}",payCheck).Methods("PUT")

log.Fatal(http.ListenAndServe(":8888",r))
}



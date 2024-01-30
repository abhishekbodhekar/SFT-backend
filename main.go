package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func main() {

	h := handler{
		db: InitDB(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", h.AddNewCustomer)
	mux.HandleFunc("/login", h.Login)

	fmt.Println("starting server on port : 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err.Error())
	}

}

func InitDB() *gorm.DB {
	dbURL := "postgres://admin:food-tracker@123@food-tracker-pg-service.default.svc.cluster.local:5432/admin"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&customers{})

	return db
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var customer customers
	err = json.Unmarshal(body, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	retrivedCust := customers{}

	result := h.db.First(&retrivedCust, "email = ?", customer.Email)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := HttpResp{
			Error:  "either email id not found or some unknown error",
			Result: false,
		}

		respByte, _ := json.Marshal(resp)

		w.Write(respByte)
		return
	}

	if retrivedCust.Passsword != GetMD5Hash(customer.Passsword) {
		w.WriteHeader(http.StatusBadRequest)
		resp := HttpResp{
			Error:  "passwords do not match",
			Result: false,
		}

		respByte, _ := json.Marshal(resp)

		w.Write(respByte)
		return
	}

	resp := HttpResp{
		Error:   "",
		Result:  true,
		Message: retrivedCust,
	}

	respByte, _ := json.Marshal(resp)

	w.Write(respByte)

}

func (h handler) AddNewCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	var customer customers
	err = json.Unmarshal(body, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	customer.Passsword = GetMD5Hash(customer.Passsword)

	result := h.db.Create(customer)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := HttpResp{
			Error:  result.Error.Error(),
			Result: false,
		}

		respByte, _ := json.Marshal(resp)

		w.Write(respByte)
		return
	}

	resp := HttpResp{
		Error:  "",
		Result: true,
	}

	respByte, _ := json.Marshal(resp)

	w.Write(respByte)

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

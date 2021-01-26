package main

import (
	"io"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"encoding/json"
	"strconv"
)

type TodoItemModel struct {
	Id int `gorm:"primary_key"`
	Description string
	Completed bool
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
        // Get URL parameter from mux
        vars := mux.Vars(r)
        id, _ := strconv.Atoi(vars["id"])

        // Test if the TodoItem exist in DB
        err := GetItemByID(id)
        if err == false {
               w.Header().Set("Content-Type", "application/json")
               io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
       } else {
               completed, _ := strconv.ParseBool(r.FormValue("completed"))
               log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
			   todo := &TodoItemModel{}
			   fmt.Println(todo)
			   fmt.Println("after test")
			   db.First(&todo, id)
			   fmt.Println(todo)
               todo.Completed = completed
               db.Save(&todo)
               w.Header().Set("Content-Type", "application/json")
                io.WriteString(w, `{"updated": true}`)
	   }
	}

func Healthz(w http.ResponseWriter, r *http.Request) {

	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"alive": true}`)
}

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")


func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	id, _ :=strconv.Atoi(r.FormValue("id"))

	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
	todo := &TodoItemModel{Description: description, Completed: false,Id: id}
	db.Create(&todo)
	//fmt.Printf("This is a test %v:",todo)
	result := db.Last(&todo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
	

}

func GetItemByID(Id int) bool {
	       todo := &TodoItemModel{}
	       result := db.First(&todo, Id)
	       if result.Error != nil{
	               log.Warn("TodoItem not found in database")
	               return false
	       }
		   return true
		}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
			       log.Info("Get completed TodoItems")
			       completedTodoItems := GetTodoItems(true)
			       w.Header().Set("Content-Type", "application/json")
			       json.NewEncoder(w).Encode(completedTodoItems)
			}
			
func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
			       log.Info("Get Incomplete TodoItems")
			       IncompleteTodoItems := GetTodoItems(false)
			       w.Header().Set("Content-Type", "application/json")
			       json.NewEncoder(w).Encode(IncompleteTodoItems)
			}
			 
func GetTodoItems(completed bool) interface{} {
			       var todos []TodoItemModel
			       TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
			       return TodoItems
			}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
				       // Get URL parameter from mux
				       vars := mux.Vars(r)
				       id, _ := strconv.Atoi(vars["id"])
				
				       // Test if the TodoItem exist in DB
				       err := GetItemByID(id)
				       if err == false {
				               w.Header().Set("Content-Type", "application/json")
				               io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
				       } else {
				               log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
				               todo := &TodoItemModel{}
				               db.First(&todo, id)
				               db.Delete(&todo)
				               w.Header().Set("Content-Type", "application/json")
				                io.WriteString(w, `{"deleted": true}`)
				       }
				}
			 

func main(){

	defer db.Close()

	//db.Debug().DropTableIfExists(&TodoItemModel{})
    db.Debug().AutoMigrate(&TodoItemModel{})
	
	log.Info("Starting Todolist API server")
	router := mux.NewRouter()
	router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
    router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo",CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)
	
}
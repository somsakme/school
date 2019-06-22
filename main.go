package main

import (
	"fmt"
	"net/http"		
	"github.com/gin-gonic/gin"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

func pingHandler(c *gin.Context){
	response:=gin.H{"message":"This is ping GET",}
	c.JSON(http.StatusOK,response)
}
type Student struct{
	Name string  `json:"name"`
	ID   int     `json:"student_id"`
}
var students=map[int]Student{
	1:Student{Name:"Anuchito",ID:1},
}

func postStudentHandler(c *gin.Context){
//receive -> Student{....}
s:=Student{}
fmt.Printf("befor bind % #v\n",s)
if err:=c.ShouldBindJSON(&s); err!=nil {
	c.JSON(http.StatusBadRequest,err)
	return
}
fmt.Printf("After bind % #v\n",s)

//add Student ->map ss
id:=len(students)
id++
s.ID=id
students[id]=s
// response
	c.JSON(http.StatusOK,s)
}

func getStudentHandler(c *gin.Context){
	ss := []Student{}
	for _,s := range students{
		ss=append(ss,s)
	}

	c.JSON(http.StatusOK,ss)
}

func pingPostHandler(c *gin.Context){
	response:=gin.H{"message":"This is ping POST",}
	c.JSON(http.StatusOK,response)
}
type Todo struct{
   ID int
   Title string
   Status string
}

func getTodosHandler(c *gin.Context){
	 db, err:=sql.Open("postgres",os.Getenv("DATABASE_URL"))
	 if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error open database":err.Error()})
		return
	 }
	 
	 stmt, _:=db.Prepare("SELECT id,title,status FROM todos")
	 rows, _:=stmt.Query()

	 todos := []Todo{}
	 for rows.Next(){
		t:=Todo{}
		err:=rows.Scan(&t.ID,&t.Title,&t.Status)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		todos=append(todos,t) 
	 }

	 c.JSON(http.StatusOK,todos)
} 


func main(){
	r:=gin.Default()
	r.GET("/api/todos",getTodosHandler)

	r.Run(":1234")
}


//	r.GET("/ping",pingHandler)
//	r.POST("/ping",pingPostHandler)

//	r.GET("/students",getStudentHandler)
//	r.POST("/students",postStudentHandler)

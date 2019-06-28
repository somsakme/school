package todo
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"school/database"
	_ "github.com/lib/pq"
)


type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type Todohandler struct {}

//func getDBConn() (*sql.DB, error) {
//	db, err := sql.Open("postgres", "postgres://kpplnkjj:owliREOn5WEm14dcPOipnHv33pUn6U_p@john.db.elephantsql.com:5432/kpplnkjj")
//	if err != nil {
//		return nil, err
//	}
//	return db, nil
//}

func (Todohandler) GetTodosHandler(c *gin.Context) {
	//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error open database =": err.Error()})
	//		return
	//	}
	db, err := database.GetDBConn()
	defer db.Close()
	fmt.Println("test-1-", os.Getenv("DATABASE_URL"))

	stmt, err := db.Prepare("SELECT id,title,status FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}

	//	fmt.Println("test-2-")

	rows, _ := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}
	fmt.Println("test-3-")

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}

	c.JSON(http.StatusOK, todos)
}

func (Todohandler)  GetTodosByIdHandler(c *gin.Context) {
	id := c.Param("id")
	//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error open database =": err.Error()})
	//		return
	//	}
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, _ := db.Prepare("SELECT id,title,status FROM todos WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	rows, _ := stmt.Query(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}

	todos := []Todo{}
	t := Todo{}
	for rows.Next() {

		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, t)
}

var todos = map[int]Todo{
	//	5: Todo{title: "test6",status:"active", ID: 5},
}

func (Todohandler) PostTodosHandler(c *gin.Context) {
	//receive -> todos{....}
	t := Todo{}
	fmt.Printf("befor post bind % #v\n", t)
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("After post bind % #v\n", t)

	//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error open database =": err.Error()})
	//		return
	//	}
	db, err := database.GetDBConn()
	defer db.Close()

	query := `
		INSERT INTO todos (title,status) VALUES ($1,$2) RETURNING id
		`
	var id int
	row := db.QueryRow(query, t.Title, t.Status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("Can't scan id", id)
	}
	fmt.Println("insert sucsess id:", id)
	t.ID = id
	c.JSON(http.StatusCreated, t)
}

func (Todohandler) DeleteTodosByIDHanderler(c *gin.Context) {
	id := c.Param("id")
	//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error open database =": err.Error()})
	//		return
	//	}
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, _ := db.Prepare("DELETE FROM todos WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	_, err = stmt.Query(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Query =": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (Todohandler) PutTodosByIDHanderler(c *gin.Context) {
	id := c.Param("id")
	t := Todo{}
	fmt.Printf("befor bind % #v\n", t)
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("After bind % #v\n", t)

	//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error open database =": err.Error()})
	//		return
	//	}
	db, err := database.GetDBConn()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE todos SET title=$2,status=$3 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare statment =": err.Error()})
		return
	}
	_, err = stmt.Query(id, t.Title, t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error update =": err.Error()})
		return
	}
	t.ID, err = strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Conv =": err.Error()})
		return
	}
	c.JSON(200, t)
}

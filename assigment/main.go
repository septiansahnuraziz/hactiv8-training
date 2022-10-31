package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"assigment/docs"
	"assigment/models"
)

var todos = []models.Todo{
	{ID: 1, Name: "Coding"},
	{ID: 2, Name: "Sleeping"},
}

const baseURL = "0.0.0.0:8000"

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	router := gin.Default()
	router.GET("/hello", HelloWorld)
	router.POST("/todo", CreateTodo)
	router.GET("/todo/:id", GetTodo)
	router.PUT("/todo/:id", UpdateTodo)

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1:8000"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(baseURL)
}

// hello godoc
// @Summary      get hallo world
// @Description get hallo world
// @Tags        assigment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  response.Response
// @Router       /hello [get]
func HelloWorld(c *gin.Context) {

	word := "Hello World"
	c.JSON(200, word)
}

// todo godoc
// @Summary      create todo
// @Description get todo world
// @Tags        assigment
// @Accept       json
// @Produce      json
// @Param        id   body      models.Todo  true  "id"
// @Success      200  {object}  response.Response
// @Router       /todo/ [post]
func CreateTodo(c *gin.Context) {

	var reqTodo models.Todo

	if errBind := c.BindJSON(&reqTodo); errBind != nil {
		c.JSON(400, "bad request")

	}

	todos = append(todos, reqTodo)

	c.JSON(200, gin.H{"message": "success create todo", "data": todos})

}

// todo godoc
// @Summary      get todo
// @Description get todo
// @Tags        assigment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  response.Response
// @Router       /todo/{id} [get]
func GetTodo(c *gin.Context) {
	id := c.Param("id")

	i, _ := strconv.Atoi(id)
	for _, a := range todos {
		if a.ID == i {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func UpdateTodo(c *gin.Context) {
	ids := c.Param("id")
	var reqTodo models.UpdateTodo
	if errBind := c.BindJSON(&reqTodo); errBind != nil {
		c.JSON(400, "bad request")

	}

	id, _ := strconv.Atoi(ids)
	for i := 0; i < len(todos); i++ {
		attr := &todos[i]
		if attr.ID == id {
			attr.Name = reqTodo.Name
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "success update todo"})
}

func DeleteTodo(c *gin.Context) {
	ids := c.Param("id")
	var reqTodo models.UpdateTodo
	if errBind := c.BindJSON(&reqTodo); errBind != nil {
		c.JSON(400, "bad request")

	}

	id, _ := strconv.Atoi(ids)
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "success delete todo"})
}

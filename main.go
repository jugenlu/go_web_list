package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Todo module
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {
	// 创建数据库
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	// 模型绑定
	DB.AutoMigrate(&Todo{})
	// 初始化一个路由
	engine := gin.Default()
	// 告诉gin模板里的静态文件去哪里找
	engine.Static("/static", "static")
	// 告诉gin模板在哪里
	engine.LoadHTMLFiles("./templates/index.html")
	// 获得这个html页面
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	v1Group := engine.Group("v1")
	{
		// add
		v1Group.POST("/todo", func(context *gin.Context) {
			// 1. 从请求中取出数据
			var todo Todo
			context.Bind(&todo)
			// 2. 存入数据库
			err := DB.Create(&todo).Error
			// 3. 返回响应
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error" : err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todo)
			}
		})
		// view all
		v1Group.GET("/todo", func(context *gin.Context) {
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"error":err.Error()})
			} else {
				context.JSON(http.StatusOK, todoList)
			}
		})
		// view one
		v1Group.GET("/todo/:id", func(context *gin.Context) {
		})
		// modify
		v1Group.PUT("/todo/:id", func(context *gin.Context) {

			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err!=nil{
				context.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			context.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err!= nil{
				context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			}else{
				context.JSON(http.StatusOK, todo)
			}
		})
		// del
		v1Group.DELETE("/todo/:id", func(context *gin.Context) {

			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error;err!=nil{
				context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			}else{
				context.JSON(http.StatusOK, gin.H{id:"deleted"})
			}
		})


	}
	engine.Run(":9090")
}

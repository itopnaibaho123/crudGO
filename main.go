package main

import (
	"database/sql"
	_ "database/sql"
	"os"
	"simple-api/auth"
	"simple-api/middleware"

	"net/http"
	"reflect"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

type newStudent struct {
	Student_id       uint64 `json: "student_id" binding:"required"`
	Student_name     string `json: "student_name" binding:"required"`
	Student_age      uint64 `json: "student_age" binding:"required"`
	Student_address  string `json: "student_address" binding:"required"`
	Student_phone_no string `json: "student_phone_no" binding:"required"`
}

func rowToStruct(rows *sql.Rows, dest interface{} ) error{
	destv:= reflect.ValueOf(dest).Elem()

	args:= make([]interface{}, destv.Type().Elem().NumField())

	for rows.Next(){
		rowp:=reflect.New(destv.Type().Elem())
		rowv:= rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}
		if err:=rows.Scan(args...); err!=nil{
			return err
		}

		destv.Set(reflect.Append(destv,rowv))
	}
	return nil
}
func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []newStudent


	db.Find(&newStudent )
	c.JSON(http.StatusOK, gin.H{
		"data": newStudent,
		"message":"Success Find All",
	
	})

}

func getHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent

	studentId := c.Param("student_id")
	// data:=newStudent{Student_id: id}
	if db.Find(&newStudent, "student_id=?",studentId).RecordNotFound(){
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil",
		"data" : newStudent,
	})
}

func postHandler(c *gin.Context, db *gorm.DB) {

	var newStudent newStudent
	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil",
		"data": newStudent,
	})
	
}
func delHandler(c *gin.Context, db *gorm.DB){

	studentId:= c.Param("student_id")
	var newStudent newStudent
	db.Delete(&newStudent, "student_id=?", studentId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Menghapus",
	})
}


func putHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent
	studentId := c.Param("student_id")
	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound(){
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	var reqStudent = newStudent
	c.Bind(&reqStudent)
	db.Model(&newStudent).Update(reqStudent)
	c.JSON(http.StatusOK, gin.H{
		"message": "success update",
		"data": reqStudent,
	})

}

func setRouter() *gin.Engine {
	errEnv := godotenv.Load(".env")
	if(errEnv!=nil){
		log.Fatal("Error Load ENV")
	}

	conn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})
	r.GET("/student",middleware.AuthValid, func(ctx *gin.Context) {
		getAllHandler(ctx,db)
	})
	r.GET("/student/:student_id", middleware.AuthValid, func(ctx *gin.Context) {
		getHandler(ctx,db)
	})
	r.PUT("/student/:student_id",middleware.AuthValid,func(ctx *gin.Context) {
		putHandler(ctx,db)
	})
	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		delHandler(ctx,db)
	})

	r.POST("/login",auth.LoginHandler)
// 
	return r
}

func Migrate(db *gorm.DB){
	db.AutoMigrate(&newStudent{})
	
	data:=newStudent{}
		if db.Find(&data).RecordNotFound(){
			fmt.Print("============================== Run Seeder User =====================")
			seederUser(db)
		}
	
}

func seederUser(db *gorm.DB){
	data := newStudent{
		Student_id: 1,
		Student_name: "JOKO",
		Student_age: 3,
		Student_address: "Jl.Johar",
		Student_phone_no: "082213713500",
	}

	db.Create(data);
}
func main() {
	r := setRouter()
	r.Run(":8080")
}

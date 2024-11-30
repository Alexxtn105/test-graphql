package main

import (
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
	_ "gorm.io/driver/sqlite"
	"log"
)

//type Blog struct {
//	gorm.Model
//	Title   string `gorm:"size:255"`
//	Content string `gorm:"type:text"`
//}

// 1. Переопределяем структуру данных (шаг 1 в README)

// Blog переопределяем структуру (старая закомментирована выше), представляющую данные
type Blog struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 2. Создаем типы GraphQL (шаг 2 в README)

// createBlogType Возвращает объект GraphQL для нашей структуры Blog
func createBlogType() *graphql.Object {
	// Возвращаем объект GraphQL
	return graphql.NewObject(
		// Описываем конфигурацию возвращаемого объекта с помощью метода ObjectConfig
		graphql.ObjectConfig{
			Name: "Blog",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"content": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
}

// 3. Определяем схему GraphQL (шаг 3 в README)

// queryType Метод определяет тип запроса для сервера GraphQL, возвращает объект GraphQL
func queryType(blogType *graphql.Object) *graphql.Object {
	// Этот метод определяет структуру и поведение запросов, обрабатываемых сервером
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query", // имя
			Fields: graphql.Fields{ // описываем запросы, которые могут быть исполнены сервером
				"blogs": &graphql.Field{
					Type: graphql.NewList(blogType),
					//описание функции получения данных
					Resolve: func(p graphql.ResolveParams) (any, error) {
						var blogs []Blog
						rows, err := db.Query("SELECT id, title, content FROM blogs")
						if err != nil {
							return nil, err
						}
						//обязательно закрываем
						defer rows.Close()
						//бежим по строкам
						for rows.Next() {
							var b Blog
							// Используем метод scan для получения данных в переменную b
							if err := rows.Scan(&b.ID, &b.Title, &b.Content); err != nil {
								return nil, err
							}
							// пишем полученные данные в переменную blogs
							blogs = append(blogs, b)
						}
						return blogs, nil
					},
				},
			},
		},
	)
}

// db База данных
var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./storage/storage.db")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	initDB()
	fmt.Println("Hello, World!")
	//r:=gin.
}

package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"time"
)

//region 1. Переопределяем структуру данных (шаг 1 в README)

//type Blog struct {
//	gorm.Model
//	Title   string `gorm:"size:255"`
//	Content string `gorm:"type:text"`
//}

// Blog переопределяем структуру (старая закомментирована выше), представляющую данные
type Blog struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

//endregion

//region 2. Создаем типы GraphQL (шаг 2 в README)

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

//endregion

//region 3. Определяем схему GraphQL (шаг 3 в README)

// queryType Метод определяет тип запроса для сервера GraphQL, возвращает объект GraphQL
func queryType(blogType *graphql.Object) *graphql.Object {
	// Этот метод определяет структуру и поведение запросов, обрабатываемых сервером
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query", // имя
			Fields: graphql.Fields{ // описываем запросы, которые могут быть исполнены сервером
				"blogs": &graphql.Field{
					Type: graphql.NewList(blogType), // тип возвращаемого значения
					Args: graphql.FieldConfigArgument{ // для использования аргументов в запросе GraphQL, например limit и offset
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) { //описание функции получения данных
						// здесь можно прочитать аргументы, указанные в разделе Args
						// читаем аргумент limit
						limit, _ := p.Args["limit"].(int)
						if limit <= 0 || limit > 20 {
							limit = 10
						}
						// читаем аргумент offset
						offset, _ := p.Args["offset"].(int)
						if offset < 0 {
							offset = 0
						}

						return getBlogs(limit, offset)
					},
				},
			},
		},
	)
}

// getBlogs Функция возвращает результат запроса
func getBlogs(limit int, offset int) ([]Blog, error) {
	var blogs []Blog
	DB.Raw("SELECT id, title, content FROM blogs LIMIT ? OFFSET ?", limit, offset).Scan(&blogs)
	return blogs, nil

	//region // Вариант для обычного sql

	//var blogs []Blog
	//rows, err := db.Query("SELECT id, title, content FROM blogs")
	//
	//if err != nil {
	//	return nil, err
	//}
	////обязательно закрываем
	//defer rows.Close()
	////бежим по строкам
	//for rows.Next() {
	//	var b Blog
	//	// Используем метод scan для получения данных в переменную b
	//	if err := rows.Scan(&b.ID, &b.Title, &b.Content); err != nil {
	//		return nil, err
	//	}
	//	// пишем полученные данные в переменную blogs
	//	blogs = append(blogs, b)
	//}

	//endregion

}

//endregion

//region База данных

// db База данных
//var db *sql.DB

var DB *gorm.DB

func initDB() {
	var err error
	//db, err = sql.Open("sqlite", "./storage/storage.db")
	//if err != nil {
	//	panic("Failed to connect DB!")
	//}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, //
			LogLevel:                  logger.Info, // уровень логирования
			IgnoreRecordNotFoundError: true,        // игнорировать ErrRecordNotFound для логгера
			Colorful:                  true,        // расцветка
		},
	)

	database, err := gorm.Open(sqlite.Open("storage/storage.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("Failed to connect DB!")
	}
	DB = database
}
func dbMigrate() {
	err := DB.AutoMigrate(
		&Blog{},
	)
	if err != nil {
		log.Printf("error on automigrate: %v", err)
		return
	}
}

//endregion

func main() {

	initDB()
	dbMigrate()

	// создаем тип блога
	blogType := createBlogType()

	// создаем схему сервера GraphQL
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType(blogType),
		})

	if err != nil {
		log.Fatalf("failed to create schema, error: %v", err)
	}

	// Шаг 4 из README - создание сервера GraphQL (с использованием graph-ql/h) - на странице https://github.com/graphql-go/handler
	// маршрут должен быть таким:
	// пишем хэндлер:
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true, // чтобы json выводился красивее
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)

}

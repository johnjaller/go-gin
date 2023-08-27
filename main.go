package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)


func main() {
	db, err := gorm.Open(sqlite.Open("db/data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	  }
	  db.AutoMigrate(&Note{})

	r := gin.Default()
	r.Use(cors.Default())
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		var notes []Note
		db.Find(&notes)
		c.HTML(
            http.StatusOK,
            "index.html",
            gin.H{
                "title": "Todo App",
				"notes":notes,
                // "url":   "./web.json",
            },
        )
	})
	r.GET("/notes", func(c *gin.Context) {
	var notes []Note
		db.Find(&notes)
		c.JSON(http.StatusOK, gin.H{
			"data":notes,
			
		})
	})
	r.POST("/notes", func(c *gin.Context) {
		var form Note
		 c.ShouldBindJSON(&form)
			db.Create(&form)
		c.JSON(http.StatusOK, gin.H{
			"data":form,
			
		})
	})
	r.PUT("/notes/:ID", func(c *gin.Context) {
		var uri uri
		var note Note
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		if err := c.ShouldBindJSON(&note); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		var result Note
		db.Model(&result).Where("id = ?", uri.ID).Updates(map[string]interface{}{"id":uri.ID,"text": note.Text, "is_completed": note.IsCompleted})
		c.JSON(http.StatusOK, gin.H{
			"note":result,
			
		})
	})
	r.DELETE("/notes/:ID", func(c *gin.Context) {
		var uri Note
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		db.Delete(&uri)
		c.JSON(http.StatusOK, gin.H{
			"deletedId":uri.ID,
			
		})
	})
	// r.GET("/ping", func(c *gin.Context) {
	
	// 	jsonFile, err := os.Open("data.json")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	parsed, _ := ioutil.ReadAll(jsonFile)
	// 	var parsedData Data

	// 	json.Unmarshal(parsed, &parsedData)
	// 	defer jsonFile.Close()
	// 	fmt.Println(jsonFile)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 		"data":    parsedData.Data,
	// 		"params":c.Query("test"),
	// 	})
	// })
	// r.GET("/pong", func(c *gin.Context) {
	// 	newData := Data{
	// 		Data: []todo{{Title: "content 1"}, {Title: "content 3"}},
	// 	}
	// 	file, _ := json.Marshal(newData)

	// 	_ = ioutil.WriteFile("data.json", file, 0644)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "ping",
	// 	})
	// })
	r.Run(":8004")
}

type Data struct {
	Data []todo `json:"data"`
}
type todo struct {
	Title string `json:"title"`
}
type uri struct {
	ID uint 
}
type Note struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Text string `json:"text"`
	IsCompleted bool `json:"isCompleted"`
  }
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "math/rand"
  "net/http"
  "time"
)

const (
  ServerToken    = "abcde"
  mainServiceUrl = "http://127.0.0.1:8000/animals/update_animal_async/"
)

func main() {
  r := InitRoutes()
  r.Run(":9000")
}

var animalsInfo = map[string]string{
  "Гепард":    "это быстроходное и грациозное хищное животное, которое отличается от остальных своими черными пятнами на светло-желтой шерсти.",
  "Синий кит": "самое крупное существо на Земле, обладающее изумительными размерами и уникальным синим оттенком кожи.",
}

type Animal struct {
  Animal_id          int    `json:"animal_id"`
  An_name            string `json:"an_name"`
  Animal_description string `json:"animal_description"`
}

func InitRoutes() *gin.Engine {
  r := gin.Default()

  r.PUT("/", func(c *gin.Context) {
    // условная проверка авторизации
    token := c.GetHeader("Server-Token")
    if token != ServerToken {
      c.AbortWithStatusJSON(http.StatusForbidden, "неверный токен авторизации")
      return
    }

    var request Animal

    if err := c.BindJSON(&request); err != nil {
      c.AbortWithStatusJSON(http.StatusBadRequest, "неверный формат данных")
      return
    }
	request.Animal_description = animalsInfo[request.An_name]
    body, err := json.Marshal(request)
    if err != nil {
      fmt.Println(err)
    }

    // асинхронное вычисление
    go func() {
      res := performTask()
      if res {
        request.Animal_description = animalsInfo[request.An_name]
      } else {
      }

      client := &http.Client{}
      req, err := http.NewRequest("PUT", mainServiceUrl, bytes.NewBuffer(body))
      if err != nil {
        fmt.Println("Error creating request:", err)
        return
      }

      req.Header.Set("Content-Type", "application/json")
      req.Header.Set("Server-Token", ServerToken)

      _, err = client.Do(req)
      if err != nil {
        fmt.Println("Error sending request:", err)
        return
      }
    }()

    c.JSON(http.StatusOK, gin.H{"message": "заявка принята в работу"})
  })

  return r
}

func performTask() bool {
  // Задержка от 6 до 10 секунд
  delay := rand.Intn(2) + 4
  time.Sleep(time.Duration(delay) * time.Second)
	return true
  // Генерируем случайное число от [0;4]
  // Если число меньше 3, возвращаем true (успех), иначе - false (неудача)
  // 60% - успех
}
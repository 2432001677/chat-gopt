package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/2432001677/chat-gopt/db"
	"github.com/2432001677/chat-gopt/gpt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Req struct {
	Question string `json:"question"`
}

func Ask(c *gin.Context) {
	ip := c.Request.Header.Get("Authorization")
	if ip == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusUnauthorized,
			"answer": "",
			"err":    "未授权",
		})
		return
	}
	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusBadRequest,
			"answer": "",
			"err":    err.Error(),
		})
		return
	}
	if req.Question == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusBadRequest,
			"answer": "",
			"err":    "问题不能为空",
		})
		return
	}
	now := time.Now()
	answer, err := gpt.AskMe(ip, req.Question)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusInternalServerError,
			"answer": "",
			"err":    err.Error(),
		})
		return
	}
	qa := db.Qa{
		Ip:       ip,
		Question: req.Question,
		Answer:   answer,
		Time:     now,
	}
	database := db.GetMongo()
	if _, err = database.Collection("qa").InsertOne(context.Background(), qa); err != nil {
		fmt.Println(now, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"answer": answer,
		"err":    "",
	})
}

type QaRes struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Date     string `json:"date"`
}

func History(c *gin.Context) {
	ip := c.Request.Header.Get("Authorization")
	if ip == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusUnauthorized,
			"answer": "",
			"err":    "未授权",
		})
		return
	}
	database := db.GetMongo()
	cursor, err := database.Collection("qa").Find(context.Background(), bson.D{{Key: "ip", Value: ip}})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusBadRequest,
			"answer": "",
			"err":    err.Error(),
		})
		return
	}
	defer cursor.Close(context.Background())

	var qas []db.Qa
	if err = cursor.All(context.Background(), &qas); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusBadRequest,
			"answer": "",
			"err":    err.Error(),
		})
		return
	}
	res := make([]QaRes, len(qas))
	for i, v := range qas {
		res[i].Question = v.Question
		res[i].Answer = v.Answer
		res[i].Date = v.Time.Format("2006-01-02 15:04:05")
	}
	c.JSON(http.StatusOK, res)
}

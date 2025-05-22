package handlers

import (
	"alert-handler/models"
	"context"
	"alert-handler/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	users *mongo.Collection
}

func NewAuthHandler(client *mongo.Client) *AuthHandler {
	return &AuthHandler{
		users: client.Database("alerts").Collection("users"),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var dbUser models.User
	err := h.users.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&dbUser)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// In production: Use bcrypt to compare hashed passwords
	if dbUser.Password != user.Password {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
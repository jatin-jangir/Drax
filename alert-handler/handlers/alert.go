package handlers

import (
	"alert-handler/models"
	"alert-handler/utils"
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlertHandler struct {
	collection *mongo.Collection
}

func NewAlertHandler(client *mongo.Client) *AlertHandler {
	return &AlertHandler{
		collection: client.Database("alerts").Collection("alerts"),
	}
}

func (h *AlertHandler) HandleAlert(c *gin.Context) {
	// Define the exact structure of Alertmanager webhook payload
	var webhookPayload struct {
		Alerts []struct {
			Status      string            `json:"status"`
			Labels      map[string]string `json:"labels"`
			Annotations map[string]string `json:"annotations"`
			StartsAt    time.Time         `json:"startsAt"`
			EndsAt      time.Time         `json:"endsAt"`
		} `json:"alerts"`
	}

	if err := c.ShouldBindJSON(&webhookPayload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid alert format", "details": err.Error()})
		return
	}

	// Process each alert in the payload
	for _, alert := range webhookPayload.Alerts {
		// Create a proper Alert document
		alertDoc := bson.M{
			"alertname":   alert.Labels["alertname"],
			"status":      alert.Status,
			"severity":    alert.Labels["severity"],
			"instance":    alert.Labels["instance"],
			"summary":     alert.Annotations["summary"],
			"description": alert.Annotations["description"],
			"startsAt":    alert.StartsAt,
			"endsAt":      alert.EndsAt,
			"createdAt":   time.Now(),
		}

		// Insert into MongoDB
		_, err := h.collection.InsertOne(context.Background(), alertDoc)
		if err != nil {
			log.Printf("Failed to insert alert: %v", err)
			continue // Continue processing other alerts even if one fails
		}
	}

	c.JSON(200, gin.H{"status": "success", "processed_alerts": len(webhookPayload.Alerts)})
}

func (h *AlertHandler) GetAlerts(c *gin.Context) {
	// Verify JWT token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header required"})
		return
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header must be Bearer <token>"})
		return
	}

	tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	_, err := utils.ValidateToken(tokenString)
	// log.Printf(tokenString);
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token: " + err.Error() + " " + tokenString})
		return
	}

	var filter models.AlertFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(400, gin.H{"error": "Invalid filter parameters"})
		return
	}

	// Build MongoDB query
	query := bson.M{}
	if filter.AlertName != "" {
		query["alertname"] = filter.AlertName
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.Severity != "" {
		query["severity"] = filter.Severity
	}
	if filter.Instance != "" {
		query["instance"] = filter.Instance
	}
	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"alertname": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"instance": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"summary": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}
	if !filter.StartTime.IsZero() || !filter.EndTime.IsZero() {
		timeFilter := bson.M{}
		if !filter.StartTime.IsZero() {
			timeFilter["$gte"] = filter.StartTime
		}
		if !filter.EndTime.IsZero() {
			timeFilter["$lte"] = filter.EndTime
		}
		query["startsAt"] = timeFilter
	}

	// Build sort
	sort := bson.D{}
	if filter.SortBy != "" {
		order := 1
		if filter.SortOrder == -1 {
			order = -1
		}
		sort = append(sort, bson.E{Key: filter.SortBy, Value: order})
	} else {
		// Default sort by newest first
		sort = append(sort, bson.E{Key: "startsAt", Value: -1})
	}

	// Execute query
	ctx := context.Background()
	cursor, err := h.collection.Find(ctx, query, options.Find().SetSort(sort))
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	var alerts []models.Alert
	if err = cursor.All(ctx, &alerts); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode alerts"})
		return
	}

	c.JSON(200, gin.H{"alerts": alerts})
}

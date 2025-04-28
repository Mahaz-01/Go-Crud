package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"gin-crud/ent"
	"gin-crud/ent/item"
	"gin-crud/internal/models"

	"github.com/gin-gonic/gin"
)

// GetItems retrieves all items
func GetItems(c *gin.Context) {
	// Log the user making the request using JWT data from context
	username, _ := c.Request.Context().Value("username").(string)
	email, _ := c.Request.Context().Value("email").(string)
	log.Printf("GetItems called by user: %s (%s)", username, email)

	ctx := context.Background()
	items, err := models.Client.Item.
		Query().
		All(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
		log.Printf("Error retrieving items: %v", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItem retrieves an item by its ID
func GetItem(c *gin.Context) {
	// Log the user making the request using JWT data from context
	username, _ := c.Request.Context().Value("username").(string)
	email, _ := c.Request.Context().Value("email").(string)
	log.Printf("GetItem called by user: %s (%s)", username, email)

	idStr := c.Param("id")

	// Convert the string ID to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}


	
	// Query the item by the int64 ID using models.Client
	ctx := context.Background()
	item, err := models.Client.Item.
		Query().
		Where(item.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item"})
		log.Printf("Error retrieving item: %v", err)
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem creates a new item
func CreateItem(c *gin.Context) {
	// Log the user making the request using JWT data from context
	username, _ := c.Request.Context().Value("username").(string)
	email, _ := c.Request.Context().Value("email").(string)
	log.Printf("CreateItem called by user: %s (%s)", username, email)

	var newItem struct {
		Name        string `json:"name" binding:"required"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if newItem.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	ctx := context.Background()
	createdItem, err := models.Client.Item.
		Create().
		SetName(newItem.Name).
		SetPrice(newItem.Price).
		SetDescription(newItem.Description).
		Save(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		log.Printf("Error creating item: %v", err)
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

// UpdateItem updates an existing item
func UpdateItem(c *gin.Context) {
	// Log the user making the request using JWT data from context
	username, _ := c.Request.Context().Value("username").(string)
	email, _ := c.Request.Context().Value("email").(string)
	log.Printf("UpdateItem called by user: %s (%s)", username, email)

	idStr := c.Param("id")

	// Convert the string ID to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedItem struct {
		Name        string `json:"name" binding:"required"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updatedItem.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	ctx := context.Background()
	updated, err := models.Client.Item.
		UpdateOneID(id).
		SetName(updatedItem.Name).
		SetPrice(updatedItem.Price).
		SetDescription(updatedItem.Description).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		log.Printf("Error updating item: %v", err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteItem deletes an item by ID
func DeleteItem(c *gin.Context) {
	// Log the user making the request using JWT data from context
	username, _ := c.Request.Context().Value("username").(string)
	email, _ := c.Request.Context().Value("email").(string)
	log.Printf("DeleteItem called by user: %s (%s)", username, email)

	idStr := c.Param("id")

	// Convert the string ID to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()
	err = models.Client.Item.
		DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		log.Printf("Error deleting item: %v", err)
		return
	}

	c.Status(http.StatusNoContent)
}

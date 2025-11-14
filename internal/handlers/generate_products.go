package handlers

import (
	"PROJECTTEST/internal/database"
	"PROJECTTEST/internal/models"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var productWords = []string{
	"Ultra", "Soft", "Premium", "Classic", "Eco", "Cotton", "Winter",
	"Summer", "Air", "Flex", "Urban", "Vintage", "Pro", "Slim",
	"Relax", "Active", "Daily", "Warm", "Fresh", "Heavy", "Light",
	"Sport", "Dry", "Original", "Modern", "Basic", "Bright", "Clean",
}

var clothingCategories = []string{
	"T-Shirt", "Hoodie", "Sweatshirt", "Pants", "Shorts", "Jacket",
	"Sneakers", "Cap", "Socks", "Backpack", "Shirt",
}

const productImage = "https://pangaia.com/cdn/shop/files/DNA_Oversized_T-Shirt_-Summit_Blue-1.png?crop=center&height=1999&v=1755260238&width=1500"

// ---------- CORE GENERATOR ----------

func generateRandomProduct() models.Product {
	rand.Seed(time.Now().UnixNano())

	// генерируем 4 случайных слова
	name := productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))]

	// описание — 6 слов
	description := productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))] + " " +
		productWords[rand.Intn(len(productWords))]

	// случайная категория
	category := clothingCategories[rand.Intn(len(clothingCategories))]

	// Цена — рандом от 0 до 200000, кратная 10000
	price := float64((rand.Intn(21)) * 10000)

	return models.Product{
		ID:          bson.NewObjectID(),
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		ImageURL:    productImage,
		CreatedAt:   time.Now(),
	}
}

func insertManyProducts(n int) error {
	products := make([]interface{}, n)
	for i := 0; i < n; i++ {
		products[i] = generateRandomProduct()
	}

	_, err := database.ProductsColl.InsertMany(context.TODO(), products)
	return err
}

// ---------- HANDLERS ----------

func Generate100Products(w http.ResponseWriter, r *http.Request) {
	err := insertManyProducts(100)
	if err != nil {
		http.Error(w, "Failed to generate products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "100 products generated successfully",
	})
}

func Generate1000Products(w http.ResponseWriter, r *http.Request) {
	err := insertManyProducts(1000)
	if err != nil {
		http.Error(w, "Failed to generate products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "1000 products generated successfully",
	})
}

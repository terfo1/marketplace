package handlers

import (
	"PROJECTTEST/internal/database"
	"PROJECTTEST/internal/helpers"
	"PROJECTTEST/internal/middleware"
	"PROJECTTEST/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	w.Write([]byte("HII"))
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if payload.Username == "" || payload.Password == "" {
		helpers.RespondError(w, http.StatusBadRequest, "username and password required")
		return
	}

	// check if username already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := database.UsersColl.CountDocuments(ctx, bson.M{"username": payload.Username})
	if count > 0 {
		helpers.RespondError(w, http.StatusBadRequest, "username already taken")
		return
	}

	// hash password
	hash, err := helpers.HashPassword(payload.Password)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	// create user
	user := models.User{
		Username:     payload.Username,
		PasswordHash: hash,
		Email:        "", // можно удалить поле из структуры, если не нужно
		CreatedAt:    time.Now(),
	}

	res, err := database.UsersColl.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	user.ID = res.InsertedID.(bson.ObjectID)

	// JWT
	token, err := middleware.CreateToken(user.ID)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "could not create token")
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	// ❗ ищем по username вместо email
	err := database.UsersColl.FindOne(ctx, bson.M{"username": payload.Username}).Decode(&user)
	if err != nil {
		helpers.RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// проверяем пароль
	if err := helpers.CheckPassword(user.PasswordHash, payload.Password); err != nil {
		helpers.RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// создаём токен
	token, err := middleware.CreateToken(user.ID)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "could not create token")
		return
	}



	log.Printf("USER FROM DB: %+v\n", user)   // 👈 добавляем
	log.Println("HASH:", user.PasswordHash)   // 👈 добавляем
	log.Println("PASS:", payload.Password)    // 👈 добавляем
	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
func MeHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := middleware.GetUserIDFromToken(r)
    if err != nil {
        helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err = database.UsersColl.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        helpers.RespondError(w, http.StatusNotFound, "user not found")
        return
    }

    helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
        "user": user,
    })
}


func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageStr := q.Get("page")
	limitStr := q.Get("limit")

	// Default values
	page := int64(1)
	limit := int64(50)

	// Parse page and limit from query parameters
	if pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 64)
		if err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 64)
		if err == nil && l > 0 {
			limit = l
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip for pagination
	skip := (page - 1) * limit

	// Use Find with pagination options
	cursor, err := database.ProductsColl.Find(ctx, bson.M{}, options.Find().SetLimit(limit).SetSkip(skip))
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}

	var items []models.Product
	if err := cursor.All(ctx, &items); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, items)
}

func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Product
	if err := database.ProductsColl.FindOne(ctx, bson.M{"_id": objID}).Decode(&p); err != nil {
		helpers.RespondError(w, http.StatusNotFound, "product not found")
		return
	}
	helpers.RespondJSON(w, http.StatusOK, p)
}

func ProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	cat := mux.Vars(r)["category"]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := database.ProductsColl.Find(ctx, bson.M{"category": cat})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	var items []models.Product
	cursor.All(ctx, &items)
	helpers.RespondJSON(w, http.StatusOK, items)
}

func ProductsSearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	pageStr := r.URL.Query().Get("page")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ------- BUILD FILTER -------
	filter := bson.M{}
	andParts := []bson.M{}

	// search by name or description
	if q != "" {
		andParts = append(andParts, bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": q, "$options": "i"}},
				{"description": bson.M{"$regex": q, "$options": "i"}},
			},
		})
	}

	// filter by category
	if category != "" {
			parts := strings.Split(category, ",")
			orList := []bson.M{}
			for _, c := range parts {
					clean := strings.TrimSpace(c)
					if clean != "" {
							orList = append(orList, bson.M{
									"category": bson.M{"$regex": "^" + clean + "$", "$options": "i"},
							})
					}
			}
			if len(orList) > 0 {
					andParts = append(andParts, bson.M{"$or": orList})
			}
	}
	// price >= min_price
	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			andParts = append(andParts, bson.M{"price": bson.M{"$gte": minPrice}})
		}
	}

	// price <= max_price
	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			andParts = append(andParts, bson.M{"price": bson.M{"$lte": maxPrice}})
		}
	}

	// combine all AND filters
	if len(andParts) > 0 {
		filter["$and"] = andParts
	}

	// ------- PAGINATION -------
	limit := int64(12)
	page := int64(1)

	if pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}

	skip := (page - 1) * limit

	// ------- TOTAL COUNT -------
	total, err := database.ProductsColl.CountDocuments(ctx, filter)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}

	// ------- FIND ITEMS -------
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := database.ProductsColl.Find(ctx, filter, opts)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}

	var items []models.Product
	if err := cursor.All(ctx, &items); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}

	// ------- RESPONSE -------
	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  page,
	})
}

func PostInteraction(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var payload struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	pid, err := bson.ObjectIDFromHex(payload.ProductID)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	// action from path
	action := strings.TrimPrefix(r.URL.Path, "/api/interactions/")
	action = strings.TrimSuffix(action, "/")
	if action == "like" {
		// toggle: if like exists -> remove, else insert
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"user_id": userID, "product_id": pid, "action_type": "like"}
		exist := database.InteractionsColl.FindOne(ctx, filter)
		if exist.Err() == nil {
			// remove like
			_, _ = database.InteractionsColl.DeleteOne(ctx, filter)
			helpers.RespondJSON(w, http.StatusOK, map[string]string{"status": "unliked"})
			return
		}
		// insert like
		it := models.Interaction{UserID: userID, ProductID: pid, ActionType: "like", Timestamp: time.Now()}
		_, err := database.InteractionsColl.InsertOne(ctx, it)
		if err != nil {
			helpers.RespondError(w, http.StatusInternalServerError, "db error")
			return
		}
		helpers.RespondJSON(w, http.StatusCreated, map[string]string{"status": "liked"})
		return
	} else if action == "view" || action == "purchase" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		it := models.Interaction{UserID: userID, ProductID: pid, ActionType: action, Timestamp: time.Now()}
		_, err := database.InteractionsColl.InsertOne(ctx, it)
		if err != nil {
			helpers.RespondError(w, http.StatusInternalServerError, "db error")
			return
		}
		helpers.RespondJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
		return
	}
	helpers.RespondError(w, http.StatusBadRequest, "unknown action")
}

func GetUserInteractionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := database.InteractionsColl.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	var items []models.Interaction
	cursor.All(ctx, &items)
	helpers.RespondJSON(w, http.StatusOK, items)
}

// Recommendations: simple collaborative filtering
// 1) find products the current user has liked or purchased
// 2) find other users who liked/purchased those products
// 3) get other products those users liked/purchased, count frequency
// 4) exclude products user already interacted with
// 5) return top N

func RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// step 1: user's positive products (like/purchase)
	cursor, err := database.InteractionsColl.Find(ctx, bson.M{"user_id": userID, "action_type": bson.M{"$in": []string{"like", "purchase"}}})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	var myInter []models.Interaction
	if err := cursor.All(ctx, &myInter); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	if len(myInter) == 0 {
		// cold-start: return newest products
		cur2, _ := database.ProductsColl.Find(ctx, bson.M{}, options.Find().SetLimit(func() int64 { n := 10; return int64(n) }()))
		var newest []models.Product
		cur2.All(ctx, &newest)
		helpers.RespondJSON(w, http.StatusOK, newest)
		return
	}
	myProductIDs := make([]bson.ObjectID, 0, len(myInter))
	for _, it := range myInter {
		myProductIDs = append(myProductIDs, it.ProductID)
	}
	// step 2: find other users who interacted with those products (like/purchase)
	otherCursor, err := database.InteractionsColl.Find(ctx, bson.M{"product_id": bson.M{"$in": myProductIDs}, "action_type": bson.M{"$in": []string{"like", "purchase"}}})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	var others []models.Interaction
	otherCursor.All(ctx, &others)
	userSet := map[bson.ObjectID]struct{}{}
	for _, o := range others {
	if len(userSet) == 0 {
			// fallback: продукты той же категории
			first := myInter[0]

			var prod models.Product
			_ = database.ProductsColl.FindOne(ctx, bson.M{"_id": first.ProductID}).Decode(&prod)

			// берем 10 товаров той же категории
			cur, _ := database.ProductsColl.Find(ctx, bson.M{
					"category": prod.Category,
			}, options.Find().SetLimit(10))

			var rec []models.Product
			cur.All(ctx, &rec)

			helpers.RespondJSON(w, http.StatusOK, rec)
			return
	}

		userSet[o.UserID] = struct{}{}
	}
	if len(userSet) == 0 {
		helpers.RespondJSON(w, http.StatusOK, []models.Product{})
		return
	}
	otherUsers := make([]bson.ObjectID, 0, len(userSet))
	for u := range userSet {
		otherUsers = append(otherUsers, u)
	}
	// step 3: get products these users liked/purchased
	cursor3, err := database.InteractionsColl.Find(ctx, bson.M{"user_id": bson.M{"$in": otherUsers}, "action_type": bson.M{"$in": []string{"like", "purchase"}}})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "db error")
		return
	}
	var rels []models.Interaction
	cursor3.All(ctx, &rels)
	countByProduct := map[bson.ObjectID]int{}
	for _, r := range rels {
		countByProduct[r.ProductID]++
	}
	// exclude products user already interacted with
	exclude := map[bson.ObjectID]struct{}{}
	for _, id := range myProductIDs {
		exclude[id] = struct{}{}
	}
	// build list sorted by count
	type kv struct {
		id    bson.ObjectID
		score int
	}
	kvList := []kv{}
	for pid, c := range countByProduct {
		if _, ok := exclude[pid]; ok {
			continue
		}
		kvList = append(kvList, kv{pid, c})
	}
	// simple sort
	for i := 0; i < len(kvList); i++ {
		for j := i + 1; j < len(kvList); j++ {
			if kvList[j].score > kvList[i].score {
				kvList[i], kvList[j] = kvList[j], kvList[i]
			}
		}
	}
	// take top 10
	limit := 10
	resProducts := []models.Product{}
	for i, kvv := range kvList {
		if i >= limit {
			break
		}
		var p models.Product
		if err := database.ProductsColl.FindOne(ctx, bson.M{"_id": kvv.id}).Decode(&p); err == nil {
			resProducts = append(resProducts, p)
		}
	}
	helpers.RespondJSON(w, http.StatusOK, resProducts)
}

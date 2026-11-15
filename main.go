package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"

	_ "backend-shoes-store/docs"
)

//go:embed db/products.sql
var productsSQL string

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Handler struct {
	DB *sql.DB
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GenerateJWT(userID int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Secret dibaca saat dipanggil, setelah .env dimuat
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest

	// 1. Parse dan validasi request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 2. Check apakah email sudah terdaftar
	var existingUserID int

	err := h.DB.QueryRowContext(
		ctx.Request.Context(),
		"SELECT id FROM users WHERE email = $1",
		req.Email,
	).Scan(&existingUserID)

	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Email already registered",
		})
		return
	}

	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check user",
		})
		return
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	// 4. Insert user ke database
	var userID int

	query := `
        INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err = h.DB.QueryRowContext(
		ctx.Request.Context(),
		query,
		req.Name,
		req.Email,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to register user",
		})
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Register success",
		"data": gin.H{
			"id":    userID,
			"name":  req.Name,
			"email": req.Email,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user using email and password, then generate JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{} "Login success with JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid email or password"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest

	// 1. Parse dan validasi request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 2. Cari user berdasarkan email
	var (
		userID         int
		name           string
		email          string
		hashedPassword string
	)

	query := `
        SELECT id, name, email, password
        FROM users
        WHERE email = $1
    	`

	err := h.DB.QueryRowContext(
		ctx.Request.Context(),
		query,
		req.Email,
	).Scan(
		&userID,
		&name,
		&email,
		&hashedPassword,
	)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve user",
		})
		return
	}

	// 3. Verifikasi password menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	// 4. Generate JWT token
	token, err := GenerateJWT(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"data": gin.H{
			"id":    userID,
			"name":  name,
			"email": email,
		},
		"token": token,
	})
}

// @title Shoes Store API
// @version 1.0
// @description Backend API for Shoes Store
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "andreasuryatanaya"),
		getEnv("DB_NAME", "shoes_store_db"),
	)
	if pw := os.Getenv("DB_PASSWORD"); pw != "" {
		dsn += " password=" + pw
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(productsSQL); err != nil {
		log.Fatalf("failed to migrate products schema: %v", err)
	}

	h := &Handler{DB: db}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Shoes Store API. Open /swagger/index.html for the API documentation.",
			"docs":    "/swagger/index.html",
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.GET("/products", h.GetAllProducts)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}

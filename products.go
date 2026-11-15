package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  int       `json:"category_id"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ProductListResponse struct {
	Message    string     `json:"message"`
	Data       []Product  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

const (
	defaultLimit = 10
	maxLimit     = 100
)

// GetAllProducts godoc
// @Summary List products
// @Description Get all products with search, filter by category, and pagination
// @Tags products
// @Produce json
// @Param search query string false "Search by product name or brand (case-insensitive)"
// @Param category query string false "Filter by category name (e.g. Running) or category id"
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} ProductListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (h *Handler) GetAllProducts(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "page must be a positive integer"})
		return
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a positive integer"})
		return
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	search := strings.TrimSpace(ctx.Query("search"))
	category := strings.TrimSpace(ctx.Query("category"))

	var (
		conds []string
		args  []any
	)
	if search != "" {
		args = append(args, "%"+search+"%")
		conds = append(conds, fmt.Sprintf("(p.name ILIKE $%d OR p.brand ILIKE $%d)", len(args), len(args)))
	}
	if category != "" {
		if id, err := strconv.Atoi(category); err == nil {
			args = append(args, id)
			conds = append(conds, fmt.Sprintf("c.id = $%d", len(args)))
		} else {
			args = append(args, category)
			conds = append(conds, fmt.Sprintf("LOWER(c.name) = LOWER($%d)", len(args)))
		}
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	const from = " FROM products p JOIN categories c ON c.id = p.category_id"

	var total int
	if err := h.DB.QueryRow("SELECT COUNT(*)"+from+where, args...).Scan(&total); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to count products"})
		return
	}

	listArgs := append(append([]any{}, args...), limit, (page-1)*limit)
	rows, err := h.DB.Query(
		"SELECT p.id, p.name, p.brand, p.description, p.price, p.stock, p.category_id, c.name, p.created_at"+
			from+where+
			fmt.Sprintf(" ORDER BY p.id LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2),
		listArgs...,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch products"})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.Category, &p.CreatedAt); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read products"})
			return
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read products"})
		return
	}

	ctx.JSON(http.StatusOK, ProductListResponse{
		Message: "success",
		Data:    products,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

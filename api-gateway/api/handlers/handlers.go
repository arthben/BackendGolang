package handlers

import (
	"net/http"
	"sync"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/arthben/BackendGolang/api-gateway/api/middlewares"
	"github.com/arthben/BackendGolang/api-gateway/api/openweather"
	"github.com/arthben/BackendGolang/api-gateway/api/rideindego"
	"github.com/arthben/BackendGolang/api-gateway/docs"
	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
	"github.com/gin-gonic/gin"
)

type fetchResult struct {
	HTTPCode int
	Error    error
}

type Handlers struct {
	router      *gin.Engine
	rideindego  *rideindego.Service
	openweather *openweather.Service
}

func NewHandler(db database.DBService, cfg *config.EnvParams) *Handlers {

	return &Handlers{
		router:      gin.New(),
		rideindego:  rideindego.NewService(db),
		openweather: openweather.NewService(db, cfg),
	}
}

func (h *Handlers) BuildHandler() (http.Handler, error) {
	h.setupMiddlewares()
	h.setupSwagger()

	apiv1 := h.router.Group("/api/v1")
	{
		apiv1.POST("/indego-data-fetch-and-store-it-db", h.FetchAndStoreIndego)
		apiv1.GET("/stations", h.FindSpecifTime)
		apiv1.GET("/stations/:kioskId", h.FindKioskWithTime)
	}

	return http.Handler(h.router), nil
}

func (h *Handlers) setupMiddlewares() {
	h.router.Use(gin.Recovery())
	h.router.Use(middlewares.AuthMiddleware())
	h.router.Use(middlewares.CorsMiddleware())
}

func (h *Handlers) setupSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (h *Handlers) fetchIndego() fetchResult {
	httCode, err := h.rideindego.FetchAndStore()
	return fetchResult{HTTPCode: httCode, Error: err}
}

func (h *Handlers) fetchWeather() fetchResult {
	httpCode, err := h.openweather.FetchAndStore()
	return fetchResult{HTTPCode: httpCode, Error: err}
}

// FetchAndStoreIndego godoc
// @Summary Store data from Indego
// @Description.markdown data-fetch
// @Tags API
// @Produce json
// @Param Authorization header string true "Bearer secret_token_static"
// @Router /api/v1/indego-data-fetch-and-store-it-db [post]
func (h *Handlers) FetchAndStoreIndego(c *gin.Context) {
	var (
		waitGroup sync.WaitGroup
		chIndego  = make(chan fetchResult)
		chWeather = make(chan fetchResult)
	)

	waitGroup.Add(3)

	go func() {
		waitGroup.Wait()
		close(chIndego)
		close(chWeather)
	}()

	go func() {
		defer waitGroup.Done()
		result := h.fetchIndego()
		chIndego <- result
	}()

	go func() {
		defer waitGroup.Done()
		result := h.fetchWeather()
		chWeather <- result
	}()

	_ = <-chWeather
	resIndego := <-chIndego

	if resIndego.Error != nil {
		c.AbortWithStatusJSON(resIndego.HTTPCode, gin.H{"error": resIndego.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Fetch and store success"})
}

// FindKioskWithTime godoc
// @Summary Snapshot of all stations at a specified time
// @Description.markdown stations
// @Tags API
// @Produce json
// @Param Authorization header string true "Bearer secret_token_static"
// @Param at            query  string true "ex: 2019-09-01T10:00:00Z"
// @Router /api/v1/stations [get]
func (h *Handlers) FindKioskWithTime(c *gin.Context) {
	kioskId := c.Param("kioskId")
	q := c.Query("at")
	if len(q) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Parameter"})
		return
	}

	at, err := time.Parse(time.RFC3339, q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	lastUpdate, jsonIndego, httpCode, err := h.rideindego.Find(at, kioskId)
	if err != nil {
		c.AbortWithStatusJSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	jsonWeather, httpCode, err := h.openweather.Find(at, "")
	if err != nil {
		c.AbortWithStatusJSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"at":       lastUpdate,
		"stations": jsonIndego,
		"weather":  jsonWeather,
	})
}

// FindSpecifTime godoc
// @Summary Snapshot of one station at a specific time
// @Description.markdown stationsKioskId
// @Tags API
// @Produce json
// @Param Authorization header string true "Bearer secret_token_static"
// @Param at            query  string true "ex: 2019-09-01T10:00:00Z"
// @Param kioskId       path   string true "ex: 3005"
// @Router /api/v1/stations/{kioskId} [get]
func (h *Handlers) FindSpecifTime(c *gin.Context) {
	q := c.Query("at")
	if len(q) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Parameter"})
		return
	}

	at, err := time.Parse(time.RFC3339, q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	_, jsonIndego, httpCode, err := h.rideindego.Find(at, "")
	if err != nil {
		c.AbortWithStatusJSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	jsonWeather, httpCode, err := h.openweather.Find(at, "")
	if err != nil {
		c.AbortWithStatusJSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"at":       q,
		"stations": jsonIndego,
		"weather":  jsonWeather,
	})
}

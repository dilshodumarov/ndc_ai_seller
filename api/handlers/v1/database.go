package v1

import (
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/database"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type databaseRoutes struct {
	handlers.BaseHandler
	log       *zap.Logger
	cfg       *config.Config
	enforcer  *casbin.CachedEnforcer
	dbUsecase database.Database
}

// NewDatabaseRoutes creates new routes for database
func NewDatabaseRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &databaseRoutes{
		log:       option.Logger,
		cfg:       option.Config,
		enforcer:  option.Enforcer,
		dbUsecase: option.Database,
	}

	db := apiV1Group.Group("/database")
	{
		db.POST("/create", r.CreateDatabase)
		db.GET("/get/:id", r.GetDatabase)
		db.PUT("/update/:id", r.UpdateDatabase)
		db.DELETE("/delete/:id", r.DeleteDatabase)
		db.GET("/list", r.ListDatabases)
	}
}

// CreateDatabase godoc
// @Summary Create a new database record
// @Description Create a new database entry with given details
// @Tags DATABASE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param database body entity.CreateDatabaseRequest true "Database details"
// @Success 201 {object} status_http.Response{data=string} "Created successfully with GUID"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /database/create [post]
func (r *databaseRoutes) CreateDatabase(c *gin.Context) {
	var req entity.CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	id, err := r.dbUsecase.Create(c, &req)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.Created, id)
}

// GetDatabase godoc
// @Summary Get a database record by ID
// @Description Get database details by GUID
// @Tags DATABASE
// @Produce json
// @Security BearerAuth
// @Param id path string true "Database GUID"
// @Success 200 {object} status_http.Response{data=entity.Database} "Database details"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /database/get/{id} [get]
func (r *databaseRoutes) GetDatabase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	db, err := r.dbUsecase.GetByID(c, id)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if db == nil {
		r.handleResponse(c, status_http.NotFound, "database not found")
		return
	}

	r.handleResponse(c, status_http.OK, db)
}

// UpdateDatabase godoc
// @Summary Update an existing database record
// @Description Update database details by GUID
// @Tags DATABASE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Database GUID"
// @Param database body entity.UpdateDatabaseRequest true "Updated database details"
// @Success 200 {object} status_http.Response{data=string} "Updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /database/update/{id} [put]
func (r *databaseRoutes) UpdateDatabase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	var req entity.UpdateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	req.Guid = id

	if err := r.dbUsecase.Update(c, &req); err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "updated successfully")
}

// DeleteDatabase godoc
// @Summary Delete a database record by ID
// @Description Delete database entry by GUID
// @Tags DATABASE
// @Produce json
// @Security BearerAuth
// @Param id path string true "Database GUID"
// @Success 200 {object} status_http.Response{data=string} "Deleted successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /database/delete/{id} [delete]
func (r *databaseRoutes) DeleteDatabase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	if err := r.dbUsecase.Delete(c, id); err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "deleted successfully")
}

// ListDatabases godoc
// @Summary List all databases
// @Description Get list of all database entries
// @Tags DATABASE
// @Produce json
// @Security BearerAuth
// @Success 200 {object} status_http.Response{data=[]entity.Database} "List of databases"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /database/list [get]
func (r *databaseRoutes) ListDatabases(c *gin.Context) {
	list, err := r.dbUsecase.List(c)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, list)
}

func (h *databaseRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 400:
	default:
		h.log.Error(
			"response",
			zap.Int("code", status.Code),
			zap.String("status", status.Status),
			zap.Any("description", status.Description),
			zap.Any("data", data),
			zap.Any("custom_message", status.CustomMessage),
		)
	}

	c.JSON(status.Code, status_http.Response{
		Status:        status.Status,
		Description:   status.Description,
		Data:          data,
		CustomMessage: status.CustomMessage,
	})
}

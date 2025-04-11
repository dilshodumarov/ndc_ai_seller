package v1

// import (
// 	"sugurta/internal/entity"
// 	"sugurta/internal/pkg/config"
// 	"sugurta/internal/pkg/helper"
// 	"sugurta/internal/pkg/logger"
// 	"sugurta/internal/pkg/redis"
// 	"sugurta/internal/usecase"

// 	"context"
// 	"time"

// 	status_http "sugurta/internal/controller/http/http_status"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// )

// // roleRoutes defines the role handler structure.
// type roleRoutes struct {
// 	useCase   usecase.Auth
// 	log       logger.Interface
// 	validator *validator.Validate
// 	inMemory  redis.InMemoryStorageI
// 	cfg       *config.Config
// }

// func NewRoleRoutes(apiV1Group *gin.RouterGroup, u usecase.Auth, cfg *config.Config, log logger.Interface, inMemory redis.InMemoryStorageI) {
// 	r := &roleRoutes{
// 		useCase:   u,
// 		log:       log,
// 		validator: validator.New(validator.WithRequiredStructEnabled()),
// 	}

// 	roleGroup := apiV1Group.Group("/role")
// 	{
// 		roleGroup.POST("", r.createRole)
// 		roleGroup.GET("", r.getlistRoles)
// 		roleGroup.GET("/:id", r.getRoleByID)
// 		roleGroup.PUT("", r.updateRole)
// 		roleGroup.DELETE("/:id", r.deleteRole)
// 	}
// }

// // @Router /role [post]
// // @Summary Create role
// // @Description Create a new role
// // @Tags ROLE
// // @Accept json
// // @Produce json
// // @Param data body entity.CreateRoleRequest true "Role data"
// // @Success 201 {object} status_http.Response{data=entity.RoleResponse}
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// func (r *roleRoutes) createRole(c *gin.Context) {
// 	var req entity.CreateRoleRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}

// 	if err := r.validator.Struct(req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	err := r.useCase.CreateRole(ctx, req)
// 	if err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.Created, "Successfully created")
// }

// // @Router /role/{id} [get]
// // @Summary Get role by ID
// // @Description Retrieve a role by its ID
// // @Tags ROLE
// // @Accept json
// // @Produce json
// // @Param id path string true "Role ID"
// // @Success 200 {object} status_http.Response{data=entity.Role}
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 404 {object} status_http.Response{data=string} "Not Found"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// func (r *roleRoutes) getRoleByID(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		r.handleResponse(c, status_http.BadRequest, "Role ID is required")
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	role, err := r.useCase.GetRoleByID(ctx, id)
// 	if err != nil {
// 		r.handleResponse(c, status_http.NotFound, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.OK, role)
// }

// // @Router /role [get]
// // @Summary Get list of roles
// // @Description Retrieve a list of all roles
// // @Tags ROLE
// // @Accept json
// // @Produce json
// // @Param page query int false "Page"
// // @Param limit query int false "Limit"
// // @Success 200 {object} status_http.Response{data=[]entity.RoleResponse}
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// func (r *roleRoutes) getlistRoles(c *gin.Context) {
// 	page, limit := helper.GetPaginationParams(c)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	roles, err := r.useCase.GetRoles(ctx, page, limit)
// 	if err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.OK, roles)
// }

// // @Router /role [put]
// // @Summary Update role
// // @Description Update role by ID
// // @Tags ROLE
// // @Accept json
// // @Produce json
// // @Param data body entity.UpdateRoleRequest true "Role data"
// // @Success 200 {object} status_http.Response{data=string} "Role updated"
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 404 {object} status_http.Response{data=string} "Not Found"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// func (r *roleRoutes) updateRole(c *gin.Context) {

// 	var req entity.UpdateRoleRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}

// 	if err := r.validator.Struct(req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	err := r.useCase.UpdateRole(ctx, req)
// 	if err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.OK, "Successfully updated")
// }

// // @Router /role/{id} [delete]
// // @Summary Delete role
// // @Description Delete role by ID
// // @Tags ROLE
// // @Accept json
// // @Produce json
// // @Param id path string true "Role ID"
// // @Success 200 {object} status_http.Response{data=string} "Role deleted successfully"
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 404 {object} status_http.Response{data=string} "Not Found"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// func (r *roleRoutes) deleteRole(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		r.handleResponse(c, status_http.BadRequest, "Role ID is required")
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	err := r.useCase.DeleteRole(ctx, id)
// 	if err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.OK, "Role deleted successfully")
// }

// func (h *roleRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
// 	switch code := status.Code; {
// 	case code < 400:
// 	default:
// 		h.log.Error(
// 			"response",
// 			logger.Int("code", status.Code),
// 			logger.String("status", status.Status),
// 			logger.Any("description", status.Description),
// 			logger.Any("data", data),
// 			logger.Any("custom_message", status.CustomMessage),
// 		)
// 	}

// 	c.JSON(status.Code, status_http.Response{
// 		Status:        status.Status,
// 		Description:   status.Description,
// 		Data:          data,
// 		CustomMessage: status.CustomMessage,
// 	})
// }

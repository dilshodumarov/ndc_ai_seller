package v1

// import (
// 	"sugurta/internal/pkg/config"
// 	"sugurta/internal/pkg/logger"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"strings"
// 	"time"

// 	status_http "sugurta/internal/controller/http/http_status"

// 	"github.com/gin-gonic/gin"
// )

// func (h *HandlerV1) AuthMiddleware(cfg config.Config) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var (
// 			res         = &auth.V2HasAccessUserRes{}
// 			ok          bool
// 			bearerToken = c.GetHeader("Authorization")
// 			app_id      = c.GetHeader("X-API-KEY")
// 			strArr      = strings.Split(bearerToken, " ")
// 		)

// 		if len(strArr) < 1 && (strArr[0] != "Bearer" && strArr[0] != "API-KEY") {
// 			h.log.Error("---ERR->Unexpected token format")
// 			_ = c.AbortWithError(http.StatusForbidden, errors.New("token error: wrong format"))
// 			return
// 		}

// 		switch strArr[0] {
// 		case "Bearer":
// 			res, ok = h.hasAccess(c)
// 			if !ok {
// 				h.log.Error("---ERR->AuthMiddleware->hasNotAccess-->")
// 				c.Abort()
// 				return
// 			}

// 			var (
// 				resourceId    = c.GetHeader("Resource-Id")
// 				environmentId = c.GetHeader("Environment-Id")
// 				projectId     = c.Query("project-id")
// 			)

// 			if res.ProjectId != "" {
// 				projectId = res.ProjectId
// 			}
// 			if res.EnvId != "" {
// 				environmentId = res.EnvId
// 			}

// 			c.Set("resource_id", resourceId)
// 			c.Set("environment_id", environmentId)
// 			c.Set("project_id", projectId)
// 			c.Set("user_id", res.UserIdAuth)
// 		case "API-KEY":
// 			if app_id == "" {
// 				h.log.Error("--AuthMiddleware--", logger.Any("error", "invalid api-key method"))
// 				h.handleResponse(c, status_http.Unauthorized, "The request requires an user authentication.")
// 				c.Abort()
// 				return
// 			}

// 			var (
// 				appIdKey, resourceAppIdKey = app_id, app_id + "resource"

// 				err      error
// 				apiJson  []byte
// 				apikeys  = &auth.GetRes{}
// 				resource = &pb.GetResourceByEnvIDResponse{}
// 			)

// 			var appWaitkey = config.CACHE_WAIT + "-appID"
// 			_, appIdOk := h.cache.Get(appWaitkey)
// 			if !appIdOk {
// 				h.cache.Add(appWaitkey, []byte(appWaitkey), config.REDIS_KEY_TIMEOUT)
// 			}

// 			if appIdOk {
// 				ctx, cancel := context.WithTimeout(c.Request.Context(), config.REDIS_WAIT_TIMEOUT)
// 				defer cancel()

// 				for {
// 					appIdBody, ok := h.cache.Get(appIdKey)
// 					if ok {
// 						apiJson = appIdBody
// 						err = json.Unmarshal(appIdBody, &apikeys)
// 						if err != nil {
// 							h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 							c.Abort()
// 							return
// 						}
// 					}

// 					if apikeys.AppId != "" {
// 						break
// 					}

// 					if ctx.Err() == context.DeadlineExceeded {
// 						break
// 					}

// 					time.Sleep(config.REDIS_SLEEP)
// 				}
// 			}

// 			if apikeys.AppId == "" {
// 				apikeys, err = h.authService.ApiKey().GetEnvID(
// 					c.Request.Context(), &auth.GetReq{Id: app_id},
// 				)
// 				if err != nil {
// 					h.handleResponse(c, status_http.BadRequest, err.Error())
// 					c.Abort()
// 					return
// 				}

// 				apiJson, err = json.Marshal(apikeys)
// 				if err != nil {
// 					h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 					c.Abort()
// 					return
// 				}

// 				go func() {
// 					h.cache.Add(appIdKey, apiJson, config.REDIS_TIMEOUT)
// 				}()
// 			}

// 			var resourceWaitKey = config.CACHE_WAIT + "-resource"
// 			_, resourceOk := h.cache.Get(resourceWaitKey)
// 			if !resourceOk {
// 				h.cache.Add(resourceWaitKey, []byte(resourceWaitKey), config.REDIS_KEY_TIMEOUT)
// 			}

// 			if resourceOk {
// 				ctx, cancel := context.WithTimeout(c.Request.Context(), config.REDIS_WAIT_TIMEOUT)
// 				defer cancel()

// 				for {
// 					resourceBody, ok := h.cache.Get(resourceAppIdKey)
// 					if ok {
// 						if err = json.Unmarshal(resourceBody, &resource); err != nil {
// 							h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 							c.Abort()
// 							return
// 						}
// 					}

// 					if resource.Resource != nil {
// 						break
// 					}

// 					if ctx.Err() == context.DeadlineExceeded {
// 						break
// 					}

// 					time.Sleep(config.REDIS_SLEEP)
// 				}
// 			}

// 			if resource.Resource == nil {
// 				resource, err = h.companyServices.Resource().GetResourceByEnvID(
// 					c.Request.Context(),
// 					&pb.GetResourceByEnvIDRequest{
// 						EnvId: apikeys.GetEnvironmentId(),
// 					},
// 				)
// 				if err != nil {
// 					h.handleResponse(c, status_http.BadRequest, err.Error())
// 					c.Abort()
// 					return
// 				}

// 				go func() {
// 					resourceBody, err := json.Marshal(resource)
// 					if err != nil {
// 						h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 						return
// 					}
// 					h.cache.Add(resourceAppIdKey, resourceBody, config.REDIS_TIMEOUT)
// 				}()
// 			}

// 			data := make(map[string]interface{})

// 			if err = json.Unmarshal(apiJson, &data); err != nil {
// 				h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 				c.Abort()
// 				return
// 			}

// 			resourceBody, err := json.Marshal(resource)
// 			if err != nil {
// 				h.handleResponse(c, status_http.BadRequest, "cant get auth info")
// 				return
// 			}

// 			c.Set("auth", models.AuthData{Type: "API-KEY", Data: data})
// 			c.Set("resource_id", resource.GetResource().GetId())
// 			c.Set("environment_id", apikeys.GetEnvironmentId())
// 			c.Set("project_id", apikeys.GetProjectId())
// 			c.Set("resource", string(resourceBody))
// 		default:
// 			if !strings.Contains(c.Request.URL.Path, "api") {
// 				err := errors.New("error invalid authorization method")
// 				h.log.Error("--AuthMiddleware--", logger.Error(err))
// 				h.handleResponse(c, status_http.BadRequest, err.Error())
// 				c.Abort()
// 			} else {
// 				h.log.Error("--AuthMiddleware--", logger.Any("error", "invalid authorization method"))
// 				h.handleResponse(c, status_http.Unauthorized, "The request requires an user authentication.")
// 				c.Abort()
// 			}

// 		}
// 		c.Set("Auth", res)

// 		c.Next()
// 	}
// }

package v1

import (
	"fmt"
	"sugurta/api/handlers"
	"sugurta/internal/pkg/config"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InstagramRoutes struct {
	handlers.BaseHandler
	log      *zap.Logger
	cfg      *config.Config
	enforcer *casbin.CachedEnforcer
}

// NewAuthRoutes creates a new auth routes controller
func NewInstagramRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &InstagramRoutes{
		log:      option.Logger,
		cfg:      option.Config,
		enforcer: option.Enforcer,
	}

	instagram := apiV1Group.Group("/instagram")
	{
		instagram.POST("/login", r.InstagramLogin)

	}
}

// InstagramLogin godoc
// @Summary Instagram login redirect
// @Description Redirect user to Instagram OAuth page
// @Security BearerAuth
// @Tags INSTAGRAM
// @Produce json
// @Success 302 {string} string "Redirect"
// @Router /instagram/login [post]
func (b *InstagramRoutes) InstagramLogin(c *gin.Context) {
	fmt.Println("Instagram login requested")

	authURL := "https://www.instagram.com/oauth/authorize?client_id=700909965624963&redirect_uri=https://dilshodforever.uz/v1/business/oauth/callback&response_type=code&scope=instagram_business_basic,instagram_business_manage_messages,instagram_business_manage_comments,instagram_business_content_publish,instagram_business_manage_insights"

	fmt.Println("Redirecting to:", authURL)

	c.Redirect(302, authURL)
}

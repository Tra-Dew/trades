package trades

import (
	"github.com/Tra-Dew/trades/pkg/core"
	"github.com/gin-gonic/gin"
)

// Controller ...
type Controller struct {
	authenticate *core.Authenticate
	service      Service
}

// NewController ...
func NewController(authenticate *core.Authenticate, service Service) Controller {
	return Controller{
		authenticate: authenticate,
		service:      service,
	}
}

// RegisterRoutes ...
func (c *Controller) RegisterRoutes(r *gin.RouterGroup) {
	inventory := r.Group("/inventory-write")
	{
		inventory.Use(
			c.authenticate.Middleware(),
		)

		// inventory.DELETE("", c.delete)
	}
}

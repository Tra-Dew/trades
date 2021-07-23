package trades

import (
	"net/http"

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

		inventory.POST("", c.post)
		inventory.POST(":id", c.accept)
		inventory.GET("", c.get)
		inventory.GET(":id", c.getByID)
	}
}

func (c *Controller) post(ctx *gin.Context) {
	req := new(CreateTradeOfferRequest)
	correlationID := ctx.GetString("X-Correlation-ID")
	userID := ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	res, err := c.service.Create(ctx, userID, correlationID, req)

	if err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *Controller) accept(ctx *gin.Context) {
	correlationID := ctx.GetString("X-Correlation-ID")
	userID := ctx.GetString("user_id")
	id := ctx.Param("id")

	if err := c.service.Accept(ctx, userID, correlationID, id); err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) get(ctx *gin.Context) {
	req := new(GetTradeOffersRequest)
	userID := ctx.GetString("user_id")

	if err := ctx.ShouldBindQuery(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	res, err := c.service.Get(ctx, userID, req)

	if err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) getByID(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	id := ctx.Param("id")

	res, err := c.service.GetByID(ctx, userID, id)

	if err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

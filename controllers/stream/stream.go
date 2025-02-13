package stream

import (
	"github.com/gin-gonic/gin"
)

// ProviderController is a controller type for provider endpoints
type StreamController struct{}

// NewProviderController creates a new instance of ProviderController with injected services
func NewStreamController() *StreamController {
	return &StreamController{}
}

func (ctrl *StreamController) QuicknodeLinkedAddressHook(ctx *gin.Context) {

}

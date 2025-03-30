package middleware

import "github.com/gin-gonic/gin"

type CommonMiddleware interface {
	Handle() gin.HandlerFunc
}

type MiddlewareManager struct {
	commonMiddlewares []CommonMiddleware
}

func NewMiddlewareManager(commonMiddlewares ...CommonMiddleware) *MiddlewareManager {
	return &MiddlewareManager{commonMiddlewares: commonMiddlewares}
}

func (m *MiddlewareManager) CommonHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, middleware := range m.commonMiddlewares {
			middleware.Handle()(c)
		}
	}
}

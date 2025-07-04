package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	pdb "go-utils/plugin/dbx"
)

// SwitchingDB sets the database context
func SwitchingDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Determine the database name based on request headers or other criteria
		dbName := c.GetHeader(pdb.CtxDBKey)
		if dbName == "" {
			c.Next()

		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, pdb.CtxDBKey, dbName)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

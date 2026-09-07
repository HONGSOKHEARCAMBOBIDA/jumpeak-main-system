package share

import "github.com/gin-gonic/gin"

type PaginationResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func ResponseSuccess(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"success": message})
}

func RespondDate(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}

func ResponsePagination(c *gin.Context, code int, data interface{}, pagination interface{}) {
	c.JSON(code, PaginationResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

const (
	Updated = "Updated"
	Created = "Created"
)

package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type User struct {
}

func (u *User) Register(c *gin.Context) {
	fmt.Println(11111111111111111)
	c.JSON(200, "ok")
}

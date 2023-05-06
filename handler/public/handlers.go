package public

import (
	"final/structs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context) {
	fmt.Println(ctx.Value("user_info"))
	val, ok := ctx.Value("user_info").(structs.UserRet)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, Bad Unauthorized User!",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %s!", val.Name),
		})
	}
}

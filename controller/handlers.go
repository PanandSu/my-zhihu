package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle404(r *gin.Context) {
	r.HTML(http.StatusNotFound, "404.html", nil)
}

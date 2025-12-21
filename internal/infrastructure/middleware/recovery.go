package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v\n", err)

				resp := utils.ErrorResponse(
					"INTERNAL_ERROR",
					"Beklenmeyen bir hata oluştu",
					fmt.Sprintf("%v", err),
				)
				utils.Return(w, http.StatusInternalServerError, resp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

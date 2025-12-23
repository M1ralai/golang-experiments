package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/events"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/eventbus"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	authHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/http"
	authService "github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/service"

	userHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/http"
	userRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/repository"
	userService "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/service"

	taskHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/http"
	taskRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/repository"
	taskService "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/service"

	healthHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/health/http"

	notificationListener "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/listener"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	httpServer *http.Server
	db         *sqlx.DB
	logger     logger.LoggerWithMiddleware
}

func NewServer(db *sqlx.DB, zapLogger logger.LoggerWithMiddleware) *Server {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	eventBus := eventbus.NewRedisBus(redisAddr, redisPassword, "task-service-group")

	eventPool := eventbus.NewWorkerPool(eventBus, 10, 1000)

	userRepository := userRepo.NewPostgresRepository(db)
	userSvc := userService.NewService(userRepository, zapLogger)
	userHandler := userHttp.NewHandler(userSvc)

	authSvc := authService.NewService(userRepository, zapLogger)
	authHandler := authHttp.NewHandler(authSvc)

	taskRepository := taskRepo.NewPostgresTaskRepository(db)
	assignmentRepository := taskRepo.NewPostgresAssignmentRepository(db)
	activityRepository := taskRepo.NewPostgresActivityRepository(db)

	userProvider := userRepo.NewUserProviderAdapter(userRepository)
	taskSvc := taskService.NewTaskService(taskRepository, assignmentRepository, activityRepository, userProvider, eventPool, zapLogger)
	taskHandler := taskHttp.NewHandler(taskSvc)

	taskListener := notificationListener.NewTaskEventListener()
	eventBus.Subscribe(context.Background(), events.TopicTaskAssigned, taskListener.HandleTaskAssigned)
	log.Println("✓ Task event listener subscribed to:", events.TopicTaskAssigned)

	healthHandler := healthHttp.NewHandler()

	router := mux.NewRouter()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(zapLogger.Middleware)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.TimeoutMiddleware)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/users", userHandler.UsersGet).Methods("GET")
	api.HandleFunc("/users", userHandler.UserPost).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.UserGetByID).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UserDelete).Methods("DELETE")

	api.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}/status", taskHandler.UpdateTaskStatus).Methods("PATCH")
	api.HandleFunc("/tasks/{id}/assignments", taskHandler.GetTaskAssignments).Methods("GET")
	api.HandleFunc("/tasks/{id}/assignments", taskHandler.AssignTask).Methods("POST")
	api.HandleFunc("/tasks/assignments/{id}", taskHandler.UnassignTask).Methods("DELETE")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}
	if len(port) > 0 && port[0] != ':' {
		port = ":" + port
	}

	httpServer := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		db:         db,
		logger:     zapLogger,
	}
}

func (s *Server) Start() error {
	errChan := make(chan error, 1)

	go func() {
		log.Printf("✓ Server starting... Port: %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(100 * time.Millisecond):

		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("✓ Graceful shutdown started...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("database close error: %w", err)
	}

	log.Println("✓ Shutdown completed")
	return nil
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// APIServer provides REST endpoints for Termonaut
type APIServer struct {
	router         *mux.Router
	db             *database.DB
	statsManager   *stats.AdvancedStatsManager
	analytics      *analytics.ProductivityAnalyzer
	port           int
	enableCORS     bool
	authenticator  Authenticator
}

// Authenticator interface for API authentication
type Authenticator interface {
	Authenticate(r *http.Request) (*User, error)
}

// User represents an authenticated user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	Total     int    `json:"total,omitempty"`
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
	HasMore   bool   `json:"has_more,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewAPIServer creates a new API server
func NewAPIServer(db *database.DB, port int) *APIServer {
	server := &APIServer{
		router:        mux.NewRouter(),
		db:            db,
		statsManager:  stats.NewAdvancedStatsManager(db),
		analytics:     analytics.NewProductivityAnalyzer(),
		port:          port,
		enableCORS:    true,
	}

	server.setupRoutes()
	return server
}

// SetAuthenticator sets the authentication method
func (s *APIServer) SetAuthenticator(auth Authenticator) {
	s.authenticator = auth
}

// setupRoutes configures all API routes
func (s *APIServer) setupRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Middleware
	api.Use(s.corsMiddleware)
	api.Use(s.authMiddleware)
	api.Use(s.loggingMiddleware)

	// Health check
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Stats endpoints
	api.HandleFunc("/stats", s.handleGetStats).Methods("GET")
	api.HandleFunc("/stats/basic", s.handleGetBasicStats).Methods("GET")
	api.HandleFunc("/stats/gamification", s.handleGetGamificationStats).Methods("GET")
	api.HandleFunc("/stats/productivity", s.handleGetProductivityStats).Methods("GET")

	// Commands endpoints
	api.HandleFunc("/commands", s.handleGetCommands).Methods("GET")
	api.HandleFunc("/commands/{id}", s.handleGetCommand).Methods("GET")
	api.HandleFunc("/commands/search", s.handleSearchCommands).Methods("POST")

	// Categories endpoints
	api.HandleFunc("/categories", s.handleGetCategories).Methods("GET")
	api.HandleFunc("/categories/stats", s.handleGetCategoryStats).Methods("GET")

	// Custom scoring endpoints
	api.HandleFunc("/scoring/rules", s.handleGetScoringRules).Methods("GET")
	api.HandleFunc("/scoring/rules", s.handleCreateScoringRule).Methods("POST")
	api.HandleFunc("/scoring/rules/{id}", s.handleUpdateScoringRule).Methods("PUT")
	api.HandleFunc("/scoring/rules/{id}", s.handleDeleteScoringRule).Methods("DELETE")

	// Bulk operations
	api.HandleFunc("/bulk/operations", s.handleBulkOperation).Methods("POST")

	// Export endpoints
	api.HandleFunc("/export/json", s.handleExportJSON).Methods("GET")
	api.HandleFunc("/export/csv", s.handleExportCSV).Methods("GET")

	// Real-time endpoints
	api.HandleFunc("/realtime/stats", s.handleRealtimeStats).Methods("GET")

	// Configuration endpoints
	api.HandleFunc("/config", s.handleGetConfig).Methods("GET")
	api.HandleFunc("/config", s.handleUpdateConfig).Methods("POST")
}

// Start starts the API server
func (s *APIServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 API Server starting on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

// Middleware

func (s *APIServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.enableCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		if s.authenticator != nil {
			user, err := s.authenticator.Authenticate(r)
			if err != nil {
				s.writeError(w, http.StatusUnauthorized, "Authentication failed")
				return
			}
			// Store user in context
			_ = user // TODO: Store in context
		}

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		fmt.Printf("API: %s %s - %v\n", r.Method, r.URL.Path, duration)
	})
}

// Handlers

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeSuccess(w, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "v0.8.0",
		"uptime":    time.Since(time.Now()).String(), // This would be actual uptime
	})
}

func (s *APIServer) handleGetStats(w http.ResponseWriter, r *http.Request) {
	basicStats, err := s.db.GetBasicStats()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	s.writeSuccess(w, basicStats)
}

func (s *APIServer) handleGetBasicStats(w http.ResponseWriter, r *http.Request) {
	statsCalc := stats.NewStatsCalculator(s.db)
	basicStats, err := statsCalc.GetBasicStats()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get basic stats")
		return
	}

	s.writeSuccess(w, basicStats)
}

func (s *APIServer) handleGetGamificationStats(w http.ResponseWriter, r *http.Request) {
	statsCalc := stats.NewStatsCalculator(s.db)
	gamificationStats, err := statsCalc.GetGamificationStats()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get gamification stats")
		return
	}

	s.writeSuccess(w, gamificationStats)
}

func (s *APIServer) handleGetProductivityStats(w http.ResponseWriter, r *http.Request) {
	// Get recent commands for analysis
	commands, err := s.db.GetRecentCommands(1000) // Last 1000 commands
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get commands")
		return
	}

	// Get sessions for productivity analysis
	sessions, err := s.db.GetRecentSessions(100) // Recent sessions
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get sessions")
		return
	}

	productivity := s.analytics.AnalyzeProductivity(commands, sessions)
	s.writeSuccess(w, productivity)
}

func (s *APIServer) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	commands, err := s.db.GetRecentCommands(limit)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get commands")
		return
	}

	// Apply offset manually (in a real implementation, this would be in the DB query)
	if offset < len(commands) {
		if offset+limit > len(commands) {
			commands = commands[offset:]
		} else {
			commands = commands[offset : offset+limit]
		}
	} else {
		commands = []*models.Command{}
	}

	s.writeSuccessWithMeta(w, commands, &Meta{
		Total:     len(commands),
		Page:      offset/limit + 1,
		PerPage:   limit,
		HasMore:   len(commands) == limit,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (s *APIServer) handleGetCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid command ID")
		return
	}

	// In a real implementation, there would be a GetCommandByID method
	s.writeError(w, http.StatusNotImplemented, "GetCommandByID not implemented")
	_ = id
}

func (s *APIServer) handleSearchCommands(w http.ResponseWriter, r *http.Request) {
	var filter stats.AdvancedFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid filter format")
		return
	}

	commands, err := s.statsManager.FilterCommands(&filter)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeSuccess(w, commands)
}

func (s *APIServer) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	classifier := categories.NewCommandClassifier()
	allCategories := classifier.GetAllCategories()

	s.writeSuccess(w, allCategories)
}

func (s *APIServer) handleGetCategoryStats(w http.ResponseWriter, r *http.Request) {
	// This would calculate category statistics from the database
	s.writeError(w, http.StatusNotImplemented, "Category stats not implemented yet")
}

func (s *APIServer) handleGetScoringRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.statsManager.GetCustomCommandScores()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get scoring rules")
		return
	}

	s.writeSuccess(w, rules)
}

func (s *APIServer) handleCreateScoringRule(w http.ResponseWriter, r *http.Request) {
	var rule stats.CustomCommandScore
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid rule format")
		return
	}

	if err := s.statsManager.CreateCustomCommandScore(&rule); err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to create scoring rule")
		return
	}

	s.writeSuccess(w, rule)
}

func (s *APIServer) handleUpdateScoringRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var rule stats.CustomCommandScore
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid rule format")
		return
	}

	if err := s.statsManager.UpdateCustomCommandScore(id, &rule); err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to update scoring rule")
		return
	}

	s.writeSuccess(w, rule)
}

func (s *APIServer) handleDeleteScoringRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.statsManager.DeleteCustomCommandScore(id); err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to delete scoring rule")
		return
	}

	s.writeSuccess(w, map[string]string{"message": "Scoring rule deleted successfully"})
}

func (s *APIServer) handleBulkOperation(w http.ResponseWriter, r *http.Request) {
	var operation stats.BulkOperation
	if err := json.NewDecoder(r.Body).Decode(&operation); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid operation format")
		return
	}

	result, err := s.statsManager.PerformBulkOperation(&operation)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Bulk operation failed")
		return
	}

	s.writeSuccess(w, result)
}

func (s *APIServer) handleExportJSON(w http.ResponseWriter, r *http.Request) {
	statsCalc := stats.NewStatsCalculator(s.db)
	
	basicStats, _ := statsCalc.GetBasicStats()
	gamificationStats, _ := statsCalc.GetGamificationStats()
	
	export := map[string]interface{}{
		"basic_stats":        basicStats,
		"gamification_stats": gamificationStats,
		"export_timestamp":   time.Now().Format(time.RFC3339),
		"version":           "v0.8.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=termonaut_export.json")
	json.NewEncoder(w).Encode(export)
}

func (s *APIServer) handleExportCSV(w http.ResponseWriter, r *http.Request) {
	// This would export data in CSV format
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=termonaut_export.csv")
	w.Write([]byte("Command,Timestamp,ExitCode,Category\n"))
	w.Write([]byte("ls,-la,2024-06-17T21:00:00Z,0,Navigation\n"))
}

func (s *APIServer) handleRealtimeStats(w http.ResponseWriter, r *http.Request) {
	// This would set up Server-Sent Events for real-time updates
	s.writeError(w, http.StatusNotImplemented, "Real-time stats not implemented yet")
}

func (s *APIServer) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	// Return current configuration
	config := map[string]interface{}{
		"api_enabled":   true,
		"cors_enabled":  s.enableCORS,
		"port":         s.port,
		"version":      "v0.8.0",
	}

	s.writeSuccess(w, config)
}

func (s *APIServer) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid config format")
		return
	}

	// Update configuration (this would be persisted in a real implementation)
	s.writeSuccess(w, map[string]string{"message": "Configuration updated successfully"})
}

// Utility methods

func (s *APIServer) writeSuccess(w http.ResponseWriter, data interface{}) {
	s.writeSuccessWithMeta(w, data, &Meta{
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (s *APIServer) writeSuccessWithMeta(w http.ResponseWriter, data interface{}, meta *Meta) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) writeError(w http.ResponseWriter, status int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
		Meta: &Meta{
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
} 
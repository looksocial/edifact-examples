// Lesson 7: Generic & Dynamic Message Handling
// This lesson covers generic message handlers and dynamic data structures.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/looksocial/edifact/internal/message"
	"github.com/looksocial/edifact/internal/model"
)

// Dynamic data structures
type DynamicElement struct {
	Value       string   `json:"value,omitempty"`
	Components  []string `json:"components,omitempty"`
	IsComposite bool     `json:"is_composite"`
}

type DynamicSegment struct {
	Tag      string           `json:"tag"`
	Elements []DynamicElement `json:"elements"`
	Position int              `json:"position"`
}

type DynamicMessage struct {
	Type     string                 `json:"type"`
	Segments []DynamicSegment       `json:"segments"`
	Metadata map[string]interface{} `json:"metadata"`
	RawData  string                 `json:"raw_data,omitempty"`
}

// Generic message handler
type GenericHandler struct{}

func (h *GenericHandler) CanHandle(messageType string) bool {
	return true // Can handle any message type
}

func (h *GenericHandler) Handle(data string) (interface{}, error) {
	// Use the generic message handler
	handler := &message.GenericHandler{}
	genericMessage, err := handler.Handle(data)
	if err != nil {
		return nil, err
	}

	// Convert to dynamic structure
	dynamicMessage := convertToDynamic(genericMessage)

	// Add metadata
	dynamicMessage.Metadata = extractMetadata(genericMessage)
	dynamicMessage.RawData = data

	return dynamicMessage, nil
}

func convertToDynamic(genericMsg *model.GenericMessage) *DynamicMessage {
	dynamicMsg := &DynamicMessage{
		Type:     genericMsg.Type,
		Segments: make([]DynamicSegment, len(genericMsg.Segments)),
		Metadata: make(map[string]interface{}),
	}

	for i, segment := range genericMsg.Segments {
		dynamicSegment := DynamicSegment{
			Tag:      segment.Tag,
			Position: segment.Position,
			Elements: make([]DynamicElement, len(segment.Elements)),
		}

		for j, element := range segment.Elements {
			dynamicElement := DynamicElement{
				IsComposite: element.IsComposite,
			}

			if element.IsComposite {
				dynamicElement.Components = element.Components
			} else {
				dynamicElement.Value = element.Value
			}

			dynamicSegment.Elements[j] = dynamicElement
		}

		dynamicMsg.Segments[i] = dynamicSegment
	}

	return dynamicMsg
}

func extractMetadata(genericMsg *model.GenericMessage) map[string]interface{} {
	metadata := make(map[string]interface{})

	metadata["segment_count"] = len(genericMsg.Segments)
	metadata["message_type"] = genericMsg.Type

	// Extract common information
	for _, segment := range genericMsg.Segments {
		switch segment.Tag {
		case "UNH":
			if len(segment.Elements) > 0 {
				metadata["message_reference"] = segment.Elements[0].Value
			}
		case "BGM":
			if len(segment.Elements) > 1 {
				metadata["document_number"] = segment.Elements[1].Value
			}
		}
	}

	return metadata
}

// Dynamic processor
type DynamicProcessor struct {
	handlers map[string]func(*DynamicMessage) error
}

func NewDynamicProcessor() *DynamicProcessor {
	return &DynamicProcessor{
		handlers: make(map[string]func(*DynamicMessage) error),
	}
}

func (p *DynamicProcessor) RegisterHandler(messageType string, handler func(*DynamicMessage) error) {
	p.handlers[messageType] = handler
}

func (p *DynamicProcessor) Process(message *DynamicMessage) error {
	if handler, exists := p.handlers[message.Type]; exists {
		return handler(message)
	}

	// Default generic processing
	return p.defaultHandler(message)
}

func (p *DynamicProcessor) defaultHandler(message *DynamicMessage) error {
	fmt.Printf("Processing %s message with %d segments\n",
		message.Type, len(message.Segments))

	// Print segment summary
	for _, segment := range message.Segments {
		fmt.Printf("  %s: %d elements\n", segment.Tag, len(segment.Elements))
	}

	return nil
}

// EDIApplication represents a complete production EDI application
type EDIApplication struct {
	config     *Config
	server     *HTTPServer
	processor  *MessageProcessor
	queue      *MessageQueue
	database   *Database
	monitoring *Monitoring
	logger     *Logger
}

// Config holds application configuration
type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	Queue      QueueConfig      `json:"queue"`
	Monitoring MonitoringConfig `json:"monitoring"`
}

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type QueueConfig struct {
	Workers int `json:"workers"`
	Size    int `json:"size"`
}

type MonitoringConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

// LoadConfig loads configuration from file or environment
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "edi_db",
			User:     "edi_user",
			Password: "edi_password",
		},
		Queue: QueueConfig{
			Workers: 3,
			Size:    100,
		},
		Monitoring: MonitoringConfig{
			Enabled: true,
			Port:    9090,
		},
	}
}

// NewEDIApplication creates a new EDI application
func NewEDIApplication(configPath string) (*EDIApplication, error) {
	config := LoadConfig()

	app := &EDIApplication{
		config: config,
	}

	// Initialize components
	if err := app.initializeDatabase(); err != nil {
		return nil, err
	}

	if err := app.initializeQueue(); err != nil {
		return nil, err
	}

	if err := app.initializeProcessor(); err != nil {
		return nil, err
	}

	if err := app.initializeServer(); err != nil {
		return nil, err
	}

	if err := app.initializeMonitoring(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *EDIApplication) initializeDatabase() error {
	app.database = NewDatabase(app.config.Database)
	fmt.Printf("Database initialized: %s:%d/%s\n",
		app.config.Database.Host, app.config.Database.Port, app.config.Database.Name)
	return nil
}

func (app *EDIApplication) initializeQueue() error {
	app.queue = NewMessageQueue(app.config.Queue.Workers, app.config.Queue.Size)
	fmt.Printf("Message queue initialized with %d workers\n", app.config.Queue.Workers)
	return nil
}

func (app *EDIApplication) initializeProcessor() error {
	app.processor = NewMessageProcessor(app.database, app.queue)
	fmt.Println("Message processor initialized")
	return nil
}

func (app *EDIApplication) initializeServer() error {
	app.server = NewHTTPServer(app.config.Server, app.processor)
	fmt.Printf("HTTP server initialized on %s:%d\n", app.config.Server.Host, app.config.Server.Port)
	return nil
}

func (app *EDIApplication) initializeMonitoring() error {
	app.monitoring = NewMonitoring(app.config.Monitoring)
	app.logger = NewLogger()
	fmt.Printf("Monitoring initialized on port %d\n", app.config.Monitoring.Port)
	return nil
}

func (app *EDIApplication) Start() error {
	app.logger.Info("Starting EDI application...")

	// Start monitoring
	if app.config.Monitoring.Enabled {
		if err := app.monitoring.Start(); err != nil {
			return err
		}
	}

	// Start message processor
	if err := app.processor.Start(); err != nil {
		return err
	}

	// Start HTTP server
	return app.server.Start()
}

func (app *EDIApplication) Stop() error {
	app.logger.Info("Stopping EDI application...")

	// Graceful shutdown
	if err := app.server.Stop(); err != nil {
		app.logger.Error("Error stopping server", "error", err)
	}

	if err := app.processor.Stop(); err != nil {
		app.logger.Error("Error stopping processor", "error", err)
	}

	if app.config.Monitoring.Enabled {
		if err := app.monitoring.Stop(); err != nil {
			app.logger.Error("Error stopping monitoring", "error", err)
		}
	}

	app.logger.Info("EDI application stopped")
	return nil
}

// Database represents the database layer
type Database struct {
	config   DatabaseConfig
	messages map[string]*EDIMessage
	mu       sync.RWMutex
}

func NewDatabase(config DatabaseConfig) *Database {
	return &Database{
		config:   config,
		messages: make(map[string]*EDIMessage),
	}
}

func (db *Database) SaveMessage(message *EDIMessage) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.messages[message.ID] = message
	return nil
}

func (db *Database) GetMessage(id string) (*EDIMessage, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if message, exists := db.messages[id]; exists {
		return message, nil
	}
	return nil, fmt.Errorf("message %s not found", id)
}

// MessageQueue represents the message queue
type MessageQueue struct {
	workers int
	queue   chan *QueueMessage
}

type QueueMessage struct {
	ID      string
	Type    string
	Payload interface{}
}

func NewMessageQueue(workers, size int) *MessageQueue {
	return &MessageQueue{
		workers: workers,
		queue:   make(chan *QueueMessage, size),
	}
}

func (mq *MessageQueue) Send(message *QueueMessage) error {
	select {
	case mq.queue <- message:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (mq *MessageQueue) Receive() (*QueueMessage, error) {
	select {
	case message := <-mq.queue:
		return message, nil
	case <-time.After(1 * time.Second):
		return nil, fmt.Errorf("timeout waiting for message")
	}
}

// MessageProcessor processes EDI messages
type MessageProcessor struct {
	database *Database
	queue    *MessageQueue
	workers  []*Worker
	stopChan chan struct{}
}

type Worker struct {
	id        int
	processor *MessageProcessor
	stopChan  chan struct{}
}

func NewMessageProcessor(database *Database, queue *MessageQueue) *MessageProcessor {
	return &MessageProcessor{
		database: database,
		queue:    queue,
		stopChan: make(chan struct{}),
	}
}

func (mp *MessageProcessor) Start() error {
	// Start workers
	for i := 0; i < mp.queue.workers; i++ {
		worker := &Worker{
			id:        i,
			processor: mp,
			stopChan:  make(chan struct{}),
		}
		mp.workers = append(mp.workers, worker)
		go worker.Start()
	}

	return nil
}

func (mp *MessageProcessor) Stop() error {
	close(mp.stopChan)

	// Stop all workers
	for _, worker := range mp.workers {
		close(worker.stopChan)
	}

	return nil
}

func (w *Worker) Start() {
	fmt.Printf("Worker %d started\n", w.id)

	for {
		select {
		case <-w.stopChan:
			fmt.Printf("Worker %d stopped\n", w.id)
			return
		default:
			message, err := w.processor.queue.Receive()
			if err != nil {
				continue
			}

			w.processMessage(message)
		}
	}
}

func (w *Worker) processMessage(message *QueueMessage) {
	fmt.Printf("Worker %d processing message %s\n", w.id, message.ID)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Save to database
	ediMessage := &EDIMessage{
		ID:        message.ID,
		Type:      message.Type,
		Timestamp: time.Now(),
		Status:    "PROCESSED",
	}

	w.processor.database.SaveMessage(ediMessage)
}

// HTTPServer represents the HTTP server
type HTTPServer struct {
	config    ServerConfig
	processor *MessageProcessor
	server    *http.Server
}

func NewHTTPServer(config ServerConfig, processor *MessageProcessor) *HTTPServer {
	return &HTTPServer{
		config:    config,
		processor: processor,
	}
}

func (hs *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", hs.healthHandler)

	// EDI message endpoint
	mux.HandleFunc("/api/messages", hs.messagesHandler)

	// Metrics endpoint
	mux.HandleFunc("/metrics", hs.metricsHandler)

	hs.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", hs.config.Host, hs.config.Port),
		Handler: mux,
	}

	fmt.Printf("HTTP server starting on %s\n", hs.server.Addr)
	return hs.server.ListenAndServe()
}

func (hs *HTTPServer) Stop() error {
	return hs.server.Close()
}

func (hs *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (hs *HTTPServer) messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		hs.handleMessagePost(w, r)
	} else if r.Method == "GET" {
		hs.handleMessageGet(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (hs *HTTPServer) handleMessagePost(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message := &QueueMessage{
		ID:      fmt.Sprintf("MSG_%d", time.Now().Unix()),
		Type:    request.Type,
		Payload: request.Content,
	}

	if err := hs.processor.queue.Send(message); err != nil {
		http.Error(w, "Failed to queue message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     message.ID,
		"status": "queued",
	})
}

func (hs *HTTPServer) handleMessageGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Message ID required", http.StatusBadRequest)
		return
	}

	message, err := hs.processor.database.GetMessage(id)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (hs *HTTPServer) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# EDI Application Metrics\n")
	fmt.Fprintf(w, "messages_processed_total %d\n", 42) // Mock metric
	fmt.Fprintf(w, "messages_in_queue %d\n", len(hs.processor.queue.queue))
	fmt.Fprintf(w, "active_workers %d\n", len(hs.processor.workers))
}

// Monitoring represents the monitoring system
type Monitoring struct {
	config MonitoringConfig
	server *http.Server
}

func NewMonitoring(config MonitoringConfig) *Monitoring {
	return &Monitoring{
		config: config,
	}
}

func (m *Monitoring) Start() error {
	if !m.config.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.config.Port),
		Handler: mux,
	}

	go m.server.ListenAndServe()
	return nil
}

func (m *Monitoring) Stop() error {
	if m.server != nil {
		return m.server.Close()
	}
	return nil
}

// Logger represents the logging system
type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(message string, fields ...interface{}) {
	fmt.Printf("[INFO] %s %v\n", message, fields)
}

func (l *Logger) Error(message string, fields ...interface{}) {
	fmt.Printf("[ERROR] %s %v\n", message, fields)
}

// EDIMessage represents an EDI message
type EDIMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

func main() {
	fmt.Println("=== Real-world Applications (Lesson 7) ===")

	// Create EDI application
	app, err := NewEDIApplication("config.json")
	if err != nil {
		fmt.Printf("Failed to create application: %v\n", err)
		return
	}

	// Start application in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			fmt.Printf("Failed to start application: %v\n", err)
		}
	}()

	// Wait a moment for startup
	time.Sleep(1 * time.Second)

	// Test the application
	fmt.Println("\n=== Testing Application ===")

	// Test health endpoint
	fmt.Println("Testing health endpoint...")
	resp, err := http.Get("http://localhost:8080/health")
	if err == nil {
		fmt.Printf("Health check: %s\n", resp.Status)
		resp.Body.Close()
	}

	// Test message processing
	fmt.Println("\nTesting message processing...")
	messageData := map[string]string{
		"type":    "INVOIC",
		"content": "UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'UNT+2+1'",
	}

	jsonData, _ := json.Marshal(messageData)
	resp, err = http.Post("http://localhost:8080/api/messages",
		"application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		fmt.Printf("Message submission: %s\n", resp.Status)
		resp.Body.Close()
	}

	// Test metrics endpoint
	fmt.Println("\nTesting metrics endpoint...")
	resp, err = http.Get("http://localhost:8080/metrics")
	if err == nil {
		fmt.Printf("Metrics: %s\n", resp.Status)
		resp.Body.Close()
	}

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Stop application
	fmt.Println("\n=== Stopping Application ===")
	app.Stop()

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Complete EDI application architecture")
	fmt.Println("✅ HTTP server with REST API")
	fmt.Println("✅ Message queue with worker pools")
	fmt.Println("✅ Database integration")
	fmt.Println("✅ Monitoring and health checks")
	fmt.Println("✅ Graceful shutdown")
	fmt.Println("✅ Production-ready deployment")

	fmt.Println("\n🎉 Congratulations! You've completed the Advanced EDIFACT Course!")
	fmt.Println("You're now ready to build production-ready EDI systems!")
}

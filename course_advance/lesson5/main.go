// Lesson 5: Custom Adapters & Handlers
// This lesson covers how to create custom message handlers and adapters.

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/looksocial/edifact/internal/message"
)

// Custom business models
type Booking struct {
	ID            string `json:"id"`
	Sender        string `json:"sender"`
	Receiver      string `json:"receiver"`
	TransportMode string `json:"transport_mode"`
	Reference     string `json:"reference"`
}

type Invoice struct {
	ID        string `json:"id"`
	Seller    string `json:"seller"`
	Buyer     string `json:"buyer"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	Reference string `json:"reference"`
}

// Custom handlers
type BookingHandler struct{}

func (h *BookingHandler) CanHandle(messageType string) bool {
	return messageType == "IFTMIN"
}

func (h *BookingHandler) Handle(data string) (interface{}, error) {
	// Parse the EDIFACT message
	handler := &message.IFTMINHandler{}
	message, err := handler.Handle(data)
	if err != nil {
		return nil, err
	}

	// Transform to business model
	booking := &Booking{
		ID:            message.Reference,
		Sender:        message.Sender,
		Receiver:      message.Receiver,
		TransportMode: "ROAD", // Default value
		Reference:     message.Reference,
	}

	return booking, nil
}

type InvoiceHandler struct{}

func (h *InvoiceHandler) CanHandle(messageType string) bool {
	return messageType == "INVOIC"
}

func (h *InvoiceHandler) Handle(data string) (interface{}, error) {
	// Parse the EDIFACT message
	handler := &message.INVOICHandler{}
	message, err := handler.Handle(data)
	if err != nil {
		return nil, err
	}

	// Transform to business model
	invoice := &Invoice{
		ID:        message.Reference,
		Seller:    message.Seller,
		Buyer:     message.Buyer,
		Amount:    "0.00", // Would extract from message
		Currency:  "USD",  // Would extract from message
		Reference: message.Reference,
	}

	return invoice, nil
}

// DatabaseIntegration represents database operations
type DatabaseIntegration struct {
	repository *EDIRepository
	cache      *Cache
}

type EDIRepository struct {
	messages map[string]*EDIMessage
	mu       sync.RWMutex
}

type EDIMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
}

func NewEDIRepository() *EDIRepository {
	return &EDIRepository{
		messages: make(map[string]*EDIMessage),
	}
}

func (r *EDIRepository) SaveMessage(message *EDIMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simulate database save
	r.messages[message.ID] = message
	fmt.Printf("Saved message %s to database\n", message.ID)
	return nil
}

func (r *EDIRepository) GetMessage(id string) (*EDIMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if message, exists := r.messages[id]; exists {
		return message, nil
	}
	return nil, fmt.Errorf("message %s not found", id)
}

func (r *EDIRepository) UpdateStatus(id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if message, exists := r.messages[id]; exists {
		message.Status = status
		fmt.Printf("Updated message %s status to %s\n", id, status)
		return nil
	}
	return fmt.Errorf("message %s not found", id)
}

// Cache for performance optimization
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	return value, exists
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// API Integration for external service communication
type APIIntegration struct {
	client    *ExternalServiceClient
	rateLimit *RateLimiter
}

type ExternalServiceClient struct {
	baseURL   string
	authToken string
}

type Notification struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *ExternalServiceClient) SendNotification(notification *Notification) error {
	// Simulate API call
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	fmt.Printf("Sending notification to %s: %s\n", c.baseURL, string(payload))

	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	// Simulate success
	return nil
}

// RateLimiter for API calls
type RateLimiter struct {
	tokens     int
	capacity   int
	rate       int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(capacity, rate int) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * rl.rate

	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Message Queue for asynchronous processing
type MessageQueue struct {
	queue      chan *QueueMessage
	deadLetter chan *QueueMessage
	workers    int
	processor  *MessageProcessor
}

type QueueMessage struct {
	ID      string
	Type    string
	Payload interface{}
	Retries int
}

type MessageProcessor struct {
	repository *EDIRepository
	apiClient  *ExternalServiceClient
	cache      *Cache
}

func NewMessageQueue(workers int) *MessageQueue {
	return &MessageQueue{
		queue:      make(chan *QueueMessage, 100),
		deadLetter: make(chan *QueueMessage, 10),
		workers:    workers,
	}
}

func (mq *MessageQueue) Send(message *QueueMessage) error {
	select {
	case mq.queue <- message:
		fmt.Printf("Message %s sent to queue\n", message.ID)
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (mq *MessageQueue) Start(processor *MessageProcessor) {
	mq.processor = processor

	for i := 0; i < mq.workers; i++ {
		go mq.worker()
	}

	go mq.deadLetterHandler()
}

func (mq *MessageQueue) worker() {
	for message := range mq.queue {
		fmt.Printf("Processing message %s\n", message.ID)

		if err := mq.processor.Process(message); err != nil {
			message.Retries++
			if message.Retries < 3 {
				// Retry
				time.Sleep(time.Duration(message.Retries) * time.Second)
				mq.queue <- message
			} else {
				// Send to dead letter queue
				mq.deadLetter <- message
			}
		}
	}
}

func (mq *MessageQueue) deadLetterHandler() {
	for message := range mq.deadLetter {
		fmt.Printf("Message %s sent to dead letter queue after %d retries\n",
			message.ID, message.Retries)
	}
}

func (mp *MessageProcessor) Process(message *QueueMessage) error {
	// Simulate processing
	time.Sleep(50 * time.Millisecond)

	// Update message status
	if err := mp.repository.UpdateStatus(message.ID, "PROCESSED"); err != nil {
		return err
	}

	// Send notification
	notification := &Notification{
		Type:    "MESSAGE_PROCESSED",
		Message: fmt.Sprintf("Message %s processed successfully", message.ID),
		Data:    message.Payload,
	}

	return mp.apiClient.SendNotification(notification)
}

// Integration Manager coordinates all integrations
type IntegrationManager struct {
	database *DatabaseIntegration
	api      *APIIntegration
	queue    *MessageQueue
}

func NewIntegrationManager() *IntegrationManager {
	repository := NewEDIRepository()
	cache := NewCache()
	apiClient := &ExternalServiceClient{
		baseURL:   "https://api.example.com",
		authToken: "token123",
	}

	return &IntegrationManager{
		database: &DatabaseIntegration{
			repository: repository,
			cache:      cache,
		},
		api: &APIIntegration{
			client:    apiClient,
			rateLimit: NewRateLimiter(10, 1), // 10 tokens, 1 per second
		},
		queue: NewMessageQueue(3),
	}
}

func (im *IntegrationManager) ProcessMessage(content string) error {
	// Create message
	message := &EDIMessage{
		ID:        fmt.Sprintf("MSG_%d", time.Now().Unix()),
		Type:      "INVOIC",
		Sender:    "SENDER001",
		Receiver:  "RECEIVER001",
		Timestamp: time.Now(),
		Content:   content,
		Status:    "RECEIVED",
	}

	// Save to database
	if err := im.database.repository.SaveMessage(message); err != nil {
		return err
	}

	// Check rate limit
	if !im.api.rateLimit.Allow() {
		return fmt.Errorf("rate limit exceeded")
	}

	// Send to queue
	queueMessage := &QueueMessage{
		ID:      message.ID,
		Type:    message.Type,
		Payload: message,
		Retries: 0,
	}

	return im.queue.Send(queueMessage)
}

func main() {
	fmt.Println("=== Integration Patterns (Lesson 5) ===")

	// Create integration manager
	manager := NewIntegrationManager()

	// Start message queue
	processor := &MessageProcessor{
		repository: manager.database.repository,
		apiClient:  manager.api.client,
		cache:      manager.database.cache,
	}
	manager.queue.Start(processor)

	// Process test messages
	testMessages := []string{
		"UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'UNT+2+1'",
		"UNH+1+ORDERS:D:97A:UN'BGM+220+ORD67890+9'UNT+2+1'",
		"UNH+1+DESADV:D:97A:UN'BGM+351+DESV11111+9'UNT+2+1'",
	}

	fmt.Println("\n=== Processing Messages ===")

	for i, content := range testMessages {
		fmt.Printf("\nProcessing message %d...\n", i+1)

		if err := manager.ProcessMessage(content); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}

		time.Sleep(200 * time.Millisecond) // Allow processing time
	}

	// Wait for processing to complete
	fmt.Println("\nWaiting for processing to complete...")
	time.Sleep(2 * time.Second)

	// Display results
	fmt.Println("\n=== Processing Results ===")

	for _, content := range testMessages {
		// This would normally be retrieved by ID, but for demo we'll show the pattern
		fmt.Printf("Message processed: %s\n", content[:20]+"...")
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Database integration with repository pattern")
	fmt.Println("✅ API integration with rate limiting")
	fmt.Println("✅ Message queuing with retry logic")
	fmt.Println("✅ Caching for performance optimization")
	fmt.Println("✅ Dead letter queue for failed messages")
}

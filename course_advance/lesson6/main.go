// Lesson 6: Validation & Error Handling
// This lesson covers validation and error handling in EDIFACT processing.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/looksocial/edifact/internal/dispatcher"
	"github.com/looksocial/edifact/pkg/edifact"
)

// Custom error types
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error [%s]: %s - %s", e.Code, e.Field, e.Message)
}

type ProcessingError struct {
	Message string
	Cause   error
}

func (e ProcessingError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("processing error: %s (caused by: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("processing error: %s", e.Message)
}

// Validation functions
func validateEDIFACTFormat(data string) error {
	if strings.TrimSpace(data) == "" {
		return ValidationError{
			Field:   "data",
			Message: "EDIFACT data cannot be empty",
			Code:    "EMPTY_DATA",
		}
	}

	if !strings.Contains(data, "'") {
		return ValidationError{
			Field:   "format",
			Message: "Missing segment terminator",
			Code:    "MISSING_TERMINATOR",
		}
	}

	if !strings.Contains(data, "+") {
		return ValidationError{
			Field:   "format",
			Message: "Missing element separator",
			Code:    "MISSING_SEPARATOR",
		}
	}

	return nil
}

func validateRequiredSegments(data string) error {
	requiredSegments := []string{"UNH", "UNT"}

	for _, segment := range requiredSegments {
		if !strings.Contains(data, segment) {
			return ValidationError{
				Field:   "segments",
				Message: fmt.Sprintf("Missing required segment: %s", segment),
				Code:    "MISSING_SEGMENT",
			}
		}
	}

	return nil
}

func validateMessageType(data string) error {
	detector := edifact.NewDetector()
	messageType, err := detector.DetectMessageType(data)
	if err != nil {
		return ValidationError{
			Field:   "message_type",
			Message: "Cannot detect message type",
			Code:    "INVALID_MESSAGE_TYPE",
		}
	}

	// Check if message type is supported
	supportedTypes := []string{"IFTMIN", "INVOIC", "ORDERS"}
	for _, t := range supportedTypes {
		if messageType == t {
			return nil
		}
	}

	return ValidationError{
		Field:   "message_type",
		Message: fmt.Sprintf("Unsupported message type: %s", messageType),
		Code:    "UNSUPPORTED_TYPE",
	}
}

// Comprehensive validation function
func validateMessage(data string) []error {
	var errors []error

	// Format validation
	if err := validateEDIFACTFormat(data); err != nil {
		errors = append(errors, err)
	}

	// Segment validation
	if err := validateRequiredSegments(data); err != nil {
		errors = append(errors, err)
	}

	// Message type validation
	if err := validateMessageType(data); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// Error handling wrapper
func processWithErrorHandling(data string) (interface{}, error) {
	// Step 1: Validate input
	validationErrors := validateMessage(data)
	if len(validationErrors) > 0 {
		return nil, ProcessingError{
			Message: "Validation failed",
			Cause:   validationErrors[0],
		}
	}

	// Step 2: Process message
	router := dispatcher.NewRouter()
	result, err := router.ProcessMessage(data)
	if err != nil {
		return nil, ProcessingError{
			Message: "Message processing failed",
			Cause:   err,
		}
	}

	return result, nil
}

// ValidationSeverity represents validation error severity
type ValidationSeverity int

const (
	SeverityError ValidationSeverity = iota
	SeverityWarning
	SeverityInfo
)

func (vs ValidationSeverity) String() string {
	switch vs {
	case SeverityError:
		return "ERROR"
	case SeverityWarning:
		return "WARNING"
	case SeverityInfo:
		return "INFO"
	default:
		return "UNKNOWN"
	}
}

// ValidationRule interface for business rules
type ValidationRule interface {
	Validate(message *EDIMessage) error
	GetRuleID() string
	GetSeverity() ValidationSeverity
}

// BusinessRuleValidator manages business rule validation
type BusinessRuleValidator struct {
	rules map[string]ValidationRule
	cache *ValidationCache
}

func NewBusinessRuleValidator() *BusinessRuleValidator {
	return &BusinessRuleValidator{
		rules: make(map[string]ValidationRule),
		cache: NewValidationCache(),
	}
}

func (v *BusinessRuleValidator) AddRule(rule ValidationRule) {
	v.rules[rule.GetRuleID()] = rule
}

func (v *BusinessRuleValidator) Validate(message *EDIMessage) []ValidationError {
	var errors []ValidationError

	for _, rule := range v.rules {
		if err := rule.Validate(message); err != nil {
			errors = append(errors, ValidationError{
				RuleID:   rule.GetRuleID(),
				Message:  err.Error(),
				Severity: rule.GetSeverity(),
			})
		}
	}

	return errors
}

// InvoiceValidationRule validates invoice-specific business rules
type InvoiceValidationRule struct {
	ruleID string
}

func NewInvoiceValidationRule() *InvoiceValidationRule {
	return &InvoiceValidationRule{
		ruleID: "INVOICE_VALIDATION",
	}
}

func (r *InvoiceValidationRule) GetRuleID() string {
	return r.ruleID
}

func (r *InvoiceValidationRule) GetSeverity() ValidationSeverity {
	return SeverityError
}

func (r *InvoiceValidationRule) Validate(message *EDIMessage) error {
	if message.Type != "INVOIC" {
		return nil // Not applicable
	}

	// Extract invoice data
	invoice := r.extractInvoice(message)

	// Validate business rules
	if err := r.validateAmount(invoice); err != nil {
		return fmt.Errorf("amount validation failed: %w", err)
	}

	if err := r.validateDates(invoice); err != nil {
		return fmt.Errorf("date validation failed: %w", err)
	}

	if err := r.validateLineItems(invoice); err != nil {
		return fmt.Errorf("line item validation failed: %w", err)
	}

	return nil
}

func (r *InvoiceValidationRule) extractInvoice(message *EDIMessage) *Invoice {
	// Mock extraction - in real implementation, parse EDIFACT content
	return &Invoice{
		Number:      "INV12345",
		TotalAmount: 1000.50,
		Currency:    "USD",
		IssueDate:   time.Now(),
		DueDate:     time.Now().AddDate(0, 0, 30),
		LineItems: []LineItem{
			{ID: "ITEM001", Quantity: 10, UnitPrice: 50.00, Total: 500.00},
			{ID: "ITEM002", Quantity: 5, UnitPrice: 100.10, Total: 500.50},
		},
	}
}

func (r *InvoiceValidationRule) validateAmount(invoice *Invoice) error {
	if invoice.TotalAmount <= 0 {
		return fmt.Errorf("invoice amount must be positive")
	}

	calculatedTotal := r.calculateLineItemTotal(invoice.LineItems)
	if abs(invoice.TotalAmount-calculatedTotal) > 0.01 {
		return fmt.Errorf("invoice total %.2f doesn't match line item sum %.2f",
			invoice.TotalAmount, calculatedTotal)
	}

	return nil
}

func (r *InvoiceValidationRule) validateDates(invoice *Invoice) error {
	if invoice.DueDate.Before(invoice.IssueDate) {
		return fmt.Errorf("due date cannot be before issue date")
	}

	if invoice.DueDate.Before(time.Now()) {
		return fmt.Errorf("due date cannot be in the past")
	}

	return nil
}

func (r *InvoiceValidationRule) validateLineItems(invoice *Invoice) error {
	if len(invoice.LineItems) == 0 {
		return fmt.Errorf("invoice must have at least one line item")
	}

	for i, item := range invoice.LineItems {
		if item.Quantity <= 0 {
			return fmt.Errorf("line item %d quantity must be positive", i+1)
		}

		if item.UnitPrice <= 0 {
			return fmt.Errorf("line item %d unit price must be positive", i+1)
		}

		calculatedTotal := float64(item.Quantity) * item.UnitPrice
		if abs(item.Total-calculatedTotal) > 0.01 {
			return fmt.Errorf("line item %d total %.2f doesn't match quantity * unit price %.2f",
				i+1, item.Total, calculatedTotal)
		}
	}

	return nil
}

func (r *InvoiceValidationRule) calculateLineItemTotal(items []LineItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Total
	}
	return total
}

// CrossFieldValidator validates dependencies between fields
type CrossFieldValidator struct {
	dependencies map[string][]string
	validators   map[string]CrossFieldRule
}

type CrossFieldRule interface {
	Validate(fields map[string]interface{}) error
	GetDependencies() []string
}

func NewCrossFieldValidator() *CrossFieldValidator {
	return &CrossFieldValidator{
		dependencies: make(map[string][]string),
		validators:   make(map[string]CrossFieldRule),
	}
}

func (v *CrossFieldValidator) AddRule(rule CrossFieldRule) {
	deps := rule.GetDependencies()
	for _, dep := range deps {
		v.dependencies[dep] = append(v.dependencies[dep], dep)
	}
	v.validators[fmt.Sprintf("rule_%d", len(v.validators))] = rule
}

func (v *CrossFieldValidator) Validate(fields map[string]interface{}) []ValidationError {
	var errors []ValidationError

	for ruleID, rule := range v.validators {
		if err := rule.Validate(fields); err != nil {
			errors = append(errors, ValidationError{
				RuleID:   ruleID,
				Message:  err.Error(),
				Severity: SeverityError,
			})
		}
	}

	return errors
}

// DateRangeValidator validates date range relationships
type DateRangeValidator struct {
	startDateField string
	endDateField   string
}

func NewDateRangeValidator(startField, endField string) *DateRangeValidator {
	return &DateRangeValidator{
		startDateField: startField,
		endDateField:   endField,
	}
}

func (v *DateRangeValidator) Validate(fields map[string]interface{}) error {
	startDate, startOK := fields[v.startDateField].(time.Time)
	endDate, endOK := fields[v.endDateField].(time.Time)

	if !startOK || !endOK {
		return fmt.Errorf("date fields not found or invalid")
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	return nil
}

func (v *DateRangeValidator) GetDependencies() []string {
	return []string{v.startDateField, v.endDateField}
}

// Custom Validators
type CustomValidator struct {
	validators []Validator
}

type Validator interface {
	Validate(data interface{}) error
	GetName() string
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validators: make([]Validator, 0),
	}
}

func (cv *CustomValidator) AddValidator(validator Validator) {
	cv.validators = append(cv.validators, validator)
}

func (cv *CustomValidator) Validate(data interface{}) []ValidationError {
	var errors []ValidationError

	for _, validator := range cv.validators {
		if err := validator.Validate(data); err != nil {
			errors = append(errors, ValidationError{
				RuleID:   validator.GetName(),
				Message:  err.Error(),
				Severity: SeverityError,
			})
		}
	}

	return errors
}

// CurrencyValidator validates currency codes
type CurrencyValidator struct {
	allowedCurrencies map[string]bool
}

func NewCurrencyValidator() *CurrencyValidator {
	return &CurrencyValidator{
		allowedCurrencies: map[string]bool{
			"USD": true,
			"EUR": true,
			"GBP": true,
			"JPY": true,
		},
	}
}

func (v *CurrencyValidator) Validate(data interface{}) error {
	currency, ok := data.(string)
	if !ok {
		return fmt.Errorf("currency must be a string")
	}

	if !v.allowedCurrencies[currency] {
		return fmt.Errorf("currency %s is not supported", currency)
	}

	return nil
}

func (v *CurrencyValidator) GetName() string {
	return "CURRENCY_VALIDATOR"
}

// TaxCodeValidator validates tax codes
type TaxCodeValidator struct {
	taxCodes map[string]TaxCodeInfo
}

type TaxCodeInfo struct {
	Code   string
	Active bool
	Rate   float64
}

func NewTaxCodeValidator() *TaxCodeValidator {
	return &TaxCodeValidator{
		taxCodes: map[string]TaxCodeInfo{
			"VAT": {Code: "VAT", Active: true, Rate: 0.20},
			"GST": {Code: "GST", Active: true, Rate: 0.10},
			"OLD": {Code: "OLD", Active: false, Rate: 0.15},
		},
	}
}

func (v *TaxCodeValidator) Validate(data interface{}) error {
	taxCode, ok := data.(string)
	if !ok {
		return fmt.Errorf("tax code must be a string")
	}

	if info, exists := v.taxCodes[taxCode]; !exists {
		return fmt.Errorf("invalid tax code: %s", taxCode)
	} else if !info.Active {
		return fmt.Errorf("tax code %s is inactive", taxCode)
	}

	return nil
}

func (v *TaxCodeValidator) GetName() string {
	return "TAX_CODE_VALIDATOR"
}

// Supporting types
type EDIMessage struct {
	Type    string
	Content string
}

type Invoice struct {
	Number      string
	TotalAmount float64
	Currency    string
	IssueDate   time.Time
	DueDate     time.Time
	LineItems   []LineItem
}

type LineItem struct {
	ID        string
	Quantity  int
	UnitPrice float64
	Total     float64
}

type ValidationCache struct {
	// Mock cache implementation
}

func NewValidationCache() *ValidationCache {
	return &ValidationCache{}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println("=== Advanced Validation (Lesson 6) ===")

	// Create validators
	businessValidator := NewBusinessRuleValidator()
	crossFieldValidator := NewCrossFieldValidator()
	customValidator := NewCustomValidator()

	// Add business rules
	invoiceRule := NewInvoiceValidationRule()
	businessValidator.AddRule(invoiceRule)

	// Add cross-field rules
	dateRangeRule := NewDateRangeValidator("issue_date", "due_date")
	crossFieldValidator.AddRule(dateRangeRule)

	// Add custom validators
	currencyValidator := NewCurrencyValidator()
	taxCodeValidator := NewTaxCodeValidator()
	customValidator.AddValidator(currencyValidator)
	customValidator.AddValidator(taxCodeValidator)

	// Test validation
	fmt.Println("\n=== Business Rule Validation ===")

	testMessage := &EDIMessage{
		Type:    "INVOIC",
		Content: "UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'",
	}

	errors := businessValidator.Validate(testMessage)
	if len(errors) > 0 {
		fmt.Println("Business rule validation errors:")
		for _, err := range errors {
			fmt.Printf("  [%s] %s: %s\n", err.Severity.String(), err.RuleID, err.Message)
		}
	} else {
		fmt.Println("Business rule validation passed")
	}

	// Test cross-field validation
	fmt.Println("\n=== Cross-Field Validation ===")

	fields := map[string]interface{}{
		"issue_date": time.Now(),
		"due_date":   time.Now().AddDate(0, 0, 30),
	}

	crossErrors := crossFieldValidator.Validate(fields)
	if len(crossErrors) > 0 {
		fmt.Println("Cross-field validation errors:")
		for _, err := range crossErrors {
			fmt.Printf("  [%s] %s: %s\n", err.Severity.String(), err.RuleID, err.Message)
		}
	} else {
		fmt.Println("Cross-field validation passed")
	}

	// Test custom validation
	fmt.Println("\n=== Custom Validation ===")

	// Test currency validation
	currencyErrors := customValidator.Validate("USD")
	if len(currencyErrors) > 0 {
		fmt.Println("Currency validation errors:")
		for _, err := range currencyErrors {
			fmt.Printf("  [%s] %s: %s\n", err.Severity.String(), err.RuleID, err.Message)
		}
	} else {
		fmt.Println("Currency validation passed")
	}

	// Test tax code validation
	taxErrors := customValidator.Validate("VAT")
	if len(taxErrors) > 0 {
		fmt.Println("Tax code validation errors:")
		for _, err := range taxErrors {
			fmt.Printf("  [%s] %s: %s\n", err.Severity.String(), err.RuleID, err.Message)
		}
	} else {
		fmt.Println("Tax code validation passed")
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Business rule validation with domain logic")
	fmt.Println("✅ Cross-field validation for data relationships")
	fmt.Println("✅ Custom validators for specific requirements")
	fmt.Println("✅ Validation error categorization")
	fmt.Println("✅ Comprehensive validation pipeline")
}

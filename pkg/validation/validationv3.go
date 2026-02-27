package validation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// DetailedError merepresentasikan error validasi dengan informasi lengkap
type DetailedError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
	Index   *int        `json:"index,omitempty"`
}

// ResponseError merepresentasikan error response standar
type ResponseError struct {
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	StatusCode int             `json:"status_code"`
	Errors     []DetailedError `json:"errors"`
	Err        error           `json:"-"`
}

// ValidatorHelper mengelola validasi struct dengan dukungan GORM
type ValidatorHelper struct {
	validate   *validator.Validate
	db         *gorm.DB
	mu         sync.Mutex
	exclusions map[string]string
}

// NewValidatorHelper membuat instance ValidatorHelper baru
func NewValidatorHelper(db *gorm.DB) *ValidatorHelper {
	vh := &ValidatorHelper{
		validate: validator.New(),
		db:       db,
	}

	vh.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	vh.validate.RegisterValidation("unique", vh.uniqueValidator)
	vh.validate.RegisterValidation("exists", vh.existsValidator)

	return vh
}

// uniqueValidator memastikan nilai belum ada di database
// Tag format: unique=table:column
func (vh *ValidatorHelper) uniqueValidator(fl validator.FieldLevel) bool {
	table, column, ok := parseValidatorParam(fl.Param())
	if !ok {
		return false
	}

	value := fl.Field().String()
	query := vh.db.Table(table).Where(fmt.Sprintf("%s = ?", column), value)

	// Apply exclusion jika ada (thread-safe karena diakses dalam lock)
	if excludedID := vh.getExclusion("id"); excludedID != "" {
		query = query.Where("id != ?", excludedID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false
	}

	return count == 0
}

// existsValidator memastikan nilai ada di database
// Tag format: exists=table:column
func (vh *ValidatorHelper) existsValidator(fl validator.FieldLevel) bool {
	table, column, ok := parseValidatorParam(fl.Param())
	if !ok {
		return false
	}

	value := fl.Field().String()
	var count int64
	err := vh.db.Table(table).
		Where(fmt.Sprintf("%s = ?", column), value).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// getExclusion mengambil exclusion value secara aman (harus dipanggil saat lock aktif)
func (vh *ValidatorHelper) getExclusion(key string) string {
	if vh.exclusions == nil {
		return ""
	}
	return vh.exclusions[key]
}

// ValidateStruct memvalidasi struct untuk operasi create
func (vh *ValidatorHelper) ValidateStruct(s interface{}) []DetailedError {
	return vh.ValidateRequestFormat(s, nil)
}

// ValidateForUpdate memvalidasi struct untuk operasi update dengan excluded ID
// Secara otomatis mengecualikan record dengan ID tersebut dari pengecekan unique
func (vh *ValidatorHelper) ValidateForUpdate(s interface{}, excludedID interface{}) []DetailedError {
	if !hasUniqueTag(s) {
		return vh.ValidateRequestFormat(s, nil)
	}

	exclusions := map[string]string{
		"id": fmt.Sprintf("%v", excludedID),
	}

	return vh.ValidateRequestFormat(s, exclusions)
}

// ValidateRequestFormat memvalidasi struct dengan exclusions opsional
// Thread-safe: menggunakan mutex untuk melindungi akses ke exclusions
func (vh *ValidatorHelper) ValidateRequestFormat(req interface{}, exclusions map[string]string) []DetailedError {
	vh.mu.Lock()
	vh.exclusions = exclusions
	err := vh.validate.Struct(req)
	vh.exclusions = nil
	vh.mu.Unlock()

	if err == nil {
		return nil
	}

	var errs []DetailedError
	for _, ve := range err.(validator.ValidationErrors) {
		errs = append(errs, processValidationError(ve))
	}

	return errs
}

// hasUniqueTag mengecek apakah struct (atau nested struct) memiliki tag `unique`
func hasUniqueTag(s interface{}) bool {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return structHasUniqueTag(t)
}

func structHasUniqueTag(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.Contains(field.Tag.Get("validate"), "unique=") {
			return true
		}
		// Rekursif untuk nested struct
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct && structHasUniqueTag(ft) {
			return true
		}
	}
	return false
}

// processValidationError mengkonversi validator.FieldError ke DetailedError
func processValidationError(err validator.FieldError) DetailedError {
	return DetailedError{
		Field:   err.Field(),
		Tag:     err.Tag(),
		Message: generateErrorMessage(err.Field(), err.Tag(), err.Param()),
		Value:   err.Value(),
		Index:   extractArrayIndex(err.Namespace()),
	}
}

// extractArrayIndex mengekstrak index array dari namespace (e.g. "Items[2]" -> 2)
func extractArrayIndex(namespace string) *int {
	start := strings.Index(namespace, "[")
	end := strings.Index(namespace, "]")

	if start == -1 || end == -1 || end <= start+1 {
		return nil
	}

	indexStr := namespace[start+1 : end]
	if index, err := strconv.Atoi(indexStr); err == nil {
		return &index
	}

	return nil
}

// parseValidatorParam mem-parse parameter format "table:column"
func parseValidatorParam(param string) (table, column string, ok bool) {
	parts := strings.SplitN(param, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// tagMessageTemplates untuk tag tanpa parameter
var tagMessageTemplates = map[string]string{
	"required": "%s is required",
	"email":    "%s must be a valid email address",
	"unique":   "%s already exists",
	"exists":   "%s does not exist",
	"dive":     "invalid data in %s",
}

// generateErrorMessage membuat pesan error yang human-readable
func generateErrorMessage(fieldName, tag, param string) string {
	switch tag {
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fieldName, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fieldName, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fieldName, strings.ReplaceAll(param, " ", ", "))
	}

	if tmpl, ok := tagMessageTemplates[tag]; ok {
		return fmt.Sprintf(tmpl, fieldName)
	}

	return fmt.Sprintf("%s is invalid", fieldName)
}

// --- Helper Functions ---

// HasValidationErrors mengecek apakah ada validation errors
func HasValidationErrors(errors []DetailedError) bool {
	return len(errors) > 0
}

// FormatValidationErrors memformat errors menjadi map field -> messages
// Mendukung multiple error per field
func FormatValidationErrors(errors []DetailedError) map[string][]string {
	result := make(map[string][]string)
	for _, err := range errors {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
}

// --- Response Helpers ---

func ErrorValidation(errors []DetailedError) *ResponseError {
	return &ResponseError{
		Code:       "BAD_REQUEST",
		Message:    "Validation failed",
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     errors,
	}
}

// Tambahkan method ini di validation package
func (e *ResponseError) Error() string {
	return e.Message
}

func (vh *ValidatorHelper) Validate(s interface{}) error {
	if errs := vh.ValidateStruct(s); HasValidationErrors(errs) {
		return ErrorValidation(errs)
	}
	return nil
}

func (vh *ValidatorHelper) ValidateUpdate(s interface{}, excludedID interface{}) error {
	if errs := vh.ValidateForUpdate(s, excludedID); HasValidationErrors(errs) {
		return ErrorValidation(errs)
	}
	return nil
}

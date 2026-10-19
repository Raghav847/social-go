package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}

/* improved validation err can use this in future
func (app *application) handleValidationError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		app.internalServerError(w, r, err)
		return
	}

	fieldErrors := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()

		switch fieldErr.Tag() {
		case "required":
			fieldErrors[field] = "This field is required"

		case "min":
			fieldErrors[field] = "Must be at least " + fieldErr.Param()

		case "max":
			fieldErrors[field] = "Must be at most " + fieldErr.Param()

		case "email":
			fieldErrors[field] = "Invalid email address"

		default:
			fieldErrors[field] = "Invalid value"
		}
	}

	writeJSONError(w, http.StatusBadRequest, map[string]any{
		"error":  "Validation failed",
		"fields": fieldErrors,
	})
}
*/

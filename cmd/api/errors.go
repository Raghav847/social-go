package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
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

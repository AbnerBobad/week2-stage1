package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// envelope wraps every JSON response body under a named top-level key,
// e.g. {"report": {...}} or {"error": "..."}.
type envelope map[string]any

// writeJSON marshals data as JSON, adds any provided headers, and writes
// it (with the given status code) to the http.ResponseWriter.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// readJSON decodes a single JSON object from the request body into dst,
// rejecting bodies containing more than one JSON value or unknown fields.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

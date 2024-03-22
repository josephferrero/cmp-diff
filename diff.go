package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// DiffReporter is a custom reporter that builds a nested map representing the differences.
type diffReporter struct {
	path       cmp.Path
	diffs      map[string]any
	redactions []redactedFields
}

func NewDiffReporter(fields ...redactedFields) diffReporter {
	return diffReporter{
		redactions: fields,
	}
}

type redactedFields struct {
	StructType interface{}
	Fields     []string
}

func NewRedactedFields(typ any, names ...string) redactedFields {
	return redactedFields{
		typ,
		names,
	}
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		// Ensure the root map is initialized
		if r.diffs == nil {
			r.diffs = make(map[string]any)
		}
		// Build the nested map structure
		splitPath := strings.Split(r.path.String(), ".")
		currentMap := r.diffs
		for i, step := range splitPath {
			// check if we are at the last step in the path
			if i < len(splitPath)-1 {
				if nextMap, ok := currentMap[step]; !ok {
					// the next level map on the current map does not exist so create it
					newMap := make(map[string]any)
					// add the new map k/v to the current map
					currentMap[step] = newMap
					// change the reference for the current map to the nested map
					currentMap = newMap
				} else {
					// we have the key for the next map so we change the reference to the next map
					currentMap = nextMap.(map[string]any)
				}
			} else {
				// We're at the last step in the path, so store the difference directly
				if r.isRedacted() {
					currentMap[step] = map[string]any{"From:": "(redacted)", "To:": "(redacted)"}
				} else {
					currentMap[step] = map[string]any{"From:": vx.Interface(), "To:": vy.Interface()}
				}

			}
		}
	}
}

// Helper method to check if a field is redacted
func (r *diffReporter) isRedacted() bool {
	// get the previous step type
	previousStepType := r.path.Index(-2).Type()
	for _, redaction := range r.redactions {
		// check if a redacted struct type matches the previous step type
		if previousStepType == reflect.TypeOf(redaction.StructType) {
			// get the current step field name
			currentStep := strings.Trim(r.path.Last().String(), ".")
			// check if the current step matches any of the redacted fields for the struct type
			for _, field := range redaction.Fields {
				if currentStep == field {
					return true
				}
			}
		}
	}
	return false
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	jsonBytes, err := json.MarshalIndent(r.diffs, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling diffs to JSON: %v", err)
	}
	return string(jsonBytes)
}

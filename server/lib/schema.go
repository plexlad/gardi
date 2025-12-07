package lib

// commented out sections represent future features
// maps, slices for later when adding features
import (
	"encoding/json"
	"time"
)

// The three essential types here:
// - Schema: Base definition of user values, features, and modules
// - Instance: An instance of a schema with user data
// - Visualization: Configuration on how to visualize data
// (Visualiations are specifically for web based usage here
// you can make your own for other display formats)
// Other important types:
// - Initialization: A configuration that initializes your character.
// - Variable: A value that stores something
// - Property: Takes data (such as a variable) and edits it using a formula

type VariableType string

const (
	TypeNumber  VariableType = "number"
	TypeString  VariableType = "string"
	TypeBoolean VariableType = "boolean"
	TypeEnum    VariableType = "enum"
	TypeArray   VariableType = "array"
)

type FormatType string

const (
	FormatFloor FormatType = "floor"
	FormatCeil  FormatType = "ceil"
	FormatRound FormatType = "round"
	FormatRaw   FormatType = "raw"
)

type VisualizationType string

const (
	VisGrid        VisualizationType = "grid"
	VisRow         VisualizationType = "row"
	VisAccordion   VisualizationType = "accordion"
	VisSingleField VisualizationType = "single_field"
	VisSlotTracker VisualizationType = "slot_tracker"
	VisCard        VisualizationType = "card"
	VisTabs        VisualizationType = "tabs"
	VisDefault     VisualizationType = "default" // top to down standard
	VisCustom      VisualizationType = "custom"
)

// Schema //

// The base template structure.
//
// Users defined variables for storing values, properties for values based off
// of variables, features for defining application logic, and modules for
// grouping.
type Schema struct {
	ID          string              `json:"_id"`
	Version     int                 `json:"version"`
	UserVersion int                 `json:"user_version"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Variables   map[string]Variable `json:"variables"`
	//	Properties     map[string]Property `json:"display_values"`
	//	Features       map[string]Feature  `json:"features"`
	//	Modules        map[string]Module   `json:"modules"`
	//	Initialization Initialization      `json:"initialization"`
	Visualizations []Visualization `json:"visualizations"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type Variable struct {
	Type    VariableType `json:"type"`
	Default any          `json:"default,omitempty"`
	Min     *float64     `json:"min,omitempty"`
	Max     *float64     `json:"max,omitempty"`
	Options []string     `json:"options,omitempty"` // for enum
	Items   *Variable    `json:"items,omitempty"`   // for array type
}

type Property struct {
	Formula string     `json:"formula"`
	Format  FormatType `json:"format,omitempty"`
}

type Feature struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AddsModules []string `json:"adds_modules"`
}

type Module struct {
	Name           string              `json:"name"`
	Description    string              `json:"description,omitempty"`
	AddsVariables  map[string]Variable `json:"adds_variables,omitempty"`
	AddsProperties map[string]Property `json:"adds_display_values,omitempty"`
}

// Initialization

// Allows the user to customize values
type Initialization struct {
	Steps []InitializationStep `json:"steps"`
}

// A step in initialization
type InitializationStep struct {
	Title  string  `json:"title,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}

// Prompt and schema variable to adjust
// The formula is applied to the input before being stored in variable
type Field struct {
	Prompt       string `json:"prompt"`
	VariableName string `json:"variable_name"`
	Formula      string `json:"formula"`
}

// Instance //

// An instance of data based off a schema.
type Instance struct {
	ID          string         `json:"_id"`
	SchemaID    string         `json:"schema_id"`
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Variables   map[string]any `json:"variables"`
	//	ActiveFeatures []string       `json:"active_features"`
	//	ActiveModules  []string       `json:"active_modules"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Visualization

// How data is displayed (variables and properties)
// The source modules are needed for this visualization to work.
type Visualization struct {
	Name                string            `json:"name"`
	Type                VisualizationType `json:"type"`
	ChildVisualizations []Visualization   `json:"child_visualizations,omitempty"`
	Config              json.RawMessage   `json:"config,omitempty"`
}

// Basic schema validation
func (s *Schema) Validate() error {
	// TODO: validation logic
	// - Check that all modules exist from Features
	// - Check that variables from formulas in properties exist
	// - Check for proper Initialization (fields exist, etc.)
	// - Verify enum types
	return nil
}

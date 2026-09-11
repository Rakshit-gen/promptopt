// Package types holds the stable, serializable result shapes that promptopt
// commands emit. Everything here is part of the public --json contract, so
// field names and JSON tags should be treated as an API surface.
package types

// TokenReport describes token accounting for a single operation. Counts are
// estimates unless Source is "provider", in which case they come from the
// Groq usage object.
type TokenReport struct {
	Before int `json:"before"`
	After  int `json:"after"`
	// ReductionPercent is positive when the prompt got shorter, negative when
	// it grew. Computed as (Before-After)/Before*100.
	ReductionPercent float64 `json:"reduction_percent"`
	// Source is "estimate" or "provider".
	Source string `json:"source"`
}

// Severity is the level attached to an analysis finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityError   Severity = "ERROR"
)

// Finding is a single actionable observation about a prompt.
type Finding struct {
	// ID is a stable code such as "P003" so findings can be referenced and
	// suppressed later.
	ID       string   `json:"id"`
	Severity Severity `json:"severity"`
	Title    string   `json:"title"`
	// Detail explains the problem and, where possible, how to fix it.
	Detail string `json:"detail"`
}

// ScoreCard is the set of 0-10 sub-scores shared by analyze and eval.
type ScoreCard struct {
	Overall       float64 `json:"overall"`
	Clarity       float64 `json:"clarity"`
	Specificity   float64 `json:"specificity,omitempty"`
	Completeness  float64 `json:"completeness,omitempty"`
	Consistency   float64 `json:"consistency"`
	Efficiency    float64 `json:"efficiency,omitempty"`
	Robustness    float64 `json:"robustness"`
	OutputControl float64 `json:"output_control,omitempty"`
}

// OptimizeResult is returned by `promptopt optimize`.
type OptimizeResult struct {
	Operation string      `json:"operation"`
	Result    string      `json:"result"`
	Changes   []string    `json:"changes"`
	Tokens    TokenReport `json:"tokens"`
	Analysis  ScoreCard   `json:"analysis"`
	Warnings  []string    `json:"warnings"`
}

// CompressResult is returned by `promptopt compress`.
type CompressResult struct {
	Operation string      `json:"operation"`
	Result    string      `json:"result"`
	Changes   []string    `json:"changes"`
	Tokens    TokenReport `json:"tokens"`
	// SemanticPreservation is a 0-10 estimate of how well behavior was kept.
	SemanticPreservation float64  `json:"semantic_preservation"`
	BehaviorRisks        []string `json:"behavior_risks"`
	Warnings             []string `json:"warnings"`
}

// Assumption is a labeled addition made during expansion.
type Assumption struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Note  string `json:"note"`
}

// ExpandResult is returned by `promptopt expand`.
type ExpandResult struct {
	Operation   string       `json:"operation"`
	Result      string       `json:"result"`
	Depth       string       `json:"depth"`
	Added       []string     `json:"added"`
	Assumptions []Assumption `json:"assumptions"`
	Tokens      TokenReport  `json:"tokens"`
	Warnings    []string     `json:"warnings"`
}

// AnalyzeResult is returned by `promptopt analyze`.
type AnalyzeResult struct {
	Operation string      `json:"operation"`
	Scores    ScoreCard   `json:"scores"`
	Findings  []Finding   `json:"findings"`
	Tokens    TokenReport `json:"tokens"`
	Summary   string      `json:"summary"`
}

// TransformResult is returned by `promptopt transform`.
type TransformResult struct {
	Operation string      `json:"operation"`
	Target    string      `json:"target"`
	Result    string      `json:"result"`
	Notes     []string    `json:"notes"`
	Variables []string    `json:"variables,omitempty"`
	Tokens    TokenReport `json:"tokens"`
	Warnings  []string    `json:"warnings"`
}

// EvalResult is returned by `promptopt eval`.
type EvalResult struct {
	Operation  string     `json:"operation"`
	Scores     ScoreCard  `json:"scores"`
	TestCases  []TestCase `json:"test_cases"`
	Weaknesses []string   `json:"weaknesses"`
	Summary    string     `json:"summary"`
}

// TestCase is a generated probe used during evaluation.
type TestCase struct {
	Name        string `json:"name"`
	Input       string `json:"input"`
	Expectation string `json:"expectation"`
	// Assessment is the model's judgement of how the prompt would handle it.
	Assessment string `json:"assessment"`
	// Pass is a coarse verdict; nil when not applicable.
	Pass *bool `json:"pass,omitempty"`
}

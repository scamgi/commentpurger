package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRemoveHTMLComments tests the HTML comment removal logic.
func TestRemoveHTMLComments(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple HTML comment",
			input:    "<div><!-- This is a comment -->Hello</div>",
			expected: "<html><head></head><body><div>Hello</div></body></html>",
		},
		{
			name:     "Multiple HTML comments",
			input:    "<!-- comment 1 --><p>para</p><!-- comment 2 -->",
			expected: "<html><head></head><body><p>para</p></body></html>",
		},
		{
			name:     "No comments",
			input:    "<div><p>No comments here</p></div>",
			expected: "<html><head></head><body><div><p>No comments here</p></div></body></html>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := removeHTMLComments([]byte(tc.input))
			if err != nil {
				t.Fatalf("removeHTMLComments returned an error: %v", err)
			}
			// The html.Render function might add html/head/body tags, so we check for contains
			if !strings.Contains(string(output), tc.expected) {
				t.Errorf("expected to contain:\n%s\ngot:\n%s", tc.expected, string(output))
			}
		})
	}
}

// TestRemoveCSSComments tests the CSS comment removal logic.
func TestRemoveCSSComments(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple CSS comment",
			input:    "body { /* color: red; */ color: blue; }",
			expected: "body {  color: blue; }",
		},
		{
			name:     "Multi-line CSS comment",
			input:    "p {\n  /*\n    font-size: 16px;\n  */\n  font-size: 18px;\n}",
			expected: "p {\n  \n  font-size: 18px;\n}",
		},
		{
			name:     "No comments",
			input:    "a { text-decoration: none; }",
			expected: "a { text-decoration: none; }",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := removeCSSComments([]byte(tc.input))
			if string(output) != tc.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.expected, string(output))
			}
		})
	}
}

// TestRemoveJSComments tests the JavaScript/TypeScript comment removal logic.
func TestRemoveJSComments(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single line comment",
			input:    "const x = 1; // assign one to x\nconst y = 2;",
			expected: "const x = 1; \nconst y = 2;",
		},
		{
			name:     "Multi-line comment",
			input:    "/* This is a function */\nfunction add(a, b) { return a + b; }",
			expected: "\nfunction add(a, b) { return a + b; }",
		},
		{
			name:     "Mixed comments",
			input:    "// Start\nlet a = 1; /* multi-line */\n// End",
			expected: "\nlet a = 1; \n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := removeJSComments([]byte(tc.input))
			if string(output) != tc.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tc.expected, string(output))
			}
		})
	}
}

// TestRemoveVueComments tests the Vue SFC comment removal logic.
func TestRemoveVueComments(t *testing.T) {
	input := `<template>
    <!-- Template Comment -->
    <div>{{ message }}</div>
</template>

<script>
// Script Comment
export default {
    data() {
        return {
            message: 'Hello Vue!' /* Data Comment */
        }
    }
}
</script>

<style scoped>
/* Style Comment */
.message {
    color: blue;
}
</style>
`

	// Note: The reconstructed file will have a slightly different format, but comments should be gone.
	expectedTemplate := "<div>{{ message }}</div>"
	expectedScript := "export default {\n    data() {\n        return {\n            message: 'Hello Vue!' \n        }\n    }\n}"
	expectedStyle := ".message {\n    color: blue;\n}"

	output := removeVueComments([]byte(input))
	outputStr := string(output)

	if !strings.Contains(outputStr, expectedTemplate) {
		t.Errorf("Vue template processing failed. Expected to contain %q in %q", expectedTemplate, outputStr)
	}
	if strings.Contains(outputStr, "<!-- Template Comment -->") {
		t.Error("Vue template comment not removed.")
	}

	if !strings.Contains(outputStr, expectedScript) {
		t.Errorf("Vue script processing failed. Expected to contain %q in %q", expectedScript, outputStr)
	}
	if strings.Contains(outputStr, "// Script Comment") || strings.Contains(outputStr, "/* Data Comment */") {
		t.Error("Vue script comments not removed.")
	}

	if !strings.Contains(outputStr, expectedStyle) {
		t.Errorf("Vue style processing failed. Expected to contain %q in %q", expectedStyle, outputStr)
	}
	if strings.Contains(outputStr, "/* Style Comment */") {
		t.Error("Vue style comment not removed.")
	}
}

// TestRun is an integration test for the main file processing logic.
func TestRun(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "commentpurger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Clean up the directory after the test finishes
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{
		"test.js":     "// js comment\nconsole.log('hello');",
		"test.html":   "<!-- html comment --><p>hello</p>",
		"test.css":    "/* css comment */ body {}",
		"ignored.txt": "this file should be ignored",
	}

	expectedContent := map[string]string{
		"test.js":     "\nconsole.log('hello');",
		"test.html":   "<html><head></head><body><p>hello</p></body></html>",
		"test.css":    " body {}",
		"ignored.txt": "this file should be ignored",
	}

	// Write the test files to the temp directory
	for name, content := range testFiles {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", name, err)
		}
	}

	// Run the main logic on the temporary directory
	args := []string{tmpDir}
	run(rootCmd, args)

	// Check the content of the files after running the program
	for name, expected := range expectedContent {
		content, err := os.ReadFile(filepath.Join(tmpDir, name))
		if err != nil {
			t.Fatalf("Failed to read processed file %s: %v", name, err)
		}

		if string(content) != expected {
			// Special handling for html which gets wrapped in html/body tags
			if filepath.Ext(name) == ".html" && strings.Contains(string(content), expected) {
				continue
			}
			t.Errorf("File %s content mismatch.\nExpected:\n%s\nGot:\n%s", name, expected, string(content))
		}
	}
}

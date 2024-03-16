package main

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	// Define a sample dependency tree
	depD := &Dependency{Name: "D", Version: "1.0"}
	depC := &Dependency{Name: "C", Version: "1.0"}
	depB := &Dependency{Name: "B", Version: "1.0", Dependencies: []*Dependency{depD}}
	depA := &Dependency{Name: "A", Version: "1.0", Dependencies: []*Dependency{depB, depC}}
	roots := []*Dependency{depA}

	// Call toJSON function
	jsonData, err := toJSON(roots)
	if err != nil {
		t.Errorf("toJSON returned an error: %v", err)
	}

	// Define the expected JSON output
	expectedJSON := `[
    {
        "name": "A",
        "version": "1.0",
        "dependencies": [
            {
                "name": "B",
                "version": "1.0",
                "dependencies": [
                    {
                        "name": "D",
                        "version": "1.0",
                        "dependencies": null
                    }
                ]
            },
            {
                "name": "C",
                "version": "1.0",
                "dependencies": null
            }
        ]
    }
]`

	// Compare the actual JSON output with the expected JSON output
	if string(jsonData) != expectedJSON {
		t.Errorf("Unexpected JSON output. Expected: %s, Got: %s", expectedJSON, string(jsonData))
	}
}

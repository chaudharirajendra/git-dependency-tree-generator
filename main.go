package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
)

type GitService struct {
	RepoURL   string
	BranchTag string
}

type Dependency struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Dependencies []*Dependency `json:"dependencies"`
}

type DependencyJSON struct {
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Dependencies []*DependencyJSON    `json:"dependencies"`
	Visited      map[*Dependency]bool `json:"-"`
}

func main() {

	// Input GitHub repository URL and branch/tag
	var repoURL, branchTag string
	fmt.Print("Enter GitHub repository URL: ")
	fmt.Scanln(&repoURL)
	fmt.Print("Enter branch/tag: ")
	fmt.Scanln(&branchTag)

	gs := NewGitService(repoURL, branchTag)

	jsonData, err := gs.generateDependancyTree()
	if err != nil {
		fmt.Println("Error generating dependancy tree:", err)
		return
	}

	// Print the JSON output
	fmt.Println(string(jsonData))

}

func NewGitService(repoURL, branchTag string) GitService {
	return GitService{RepoURL: repoURL, BranchTag: branchTag}
}

func (gs GitService) generateDependancyTree() ([]byte, error) {

	err := cloneRepository(gs.RepoURL)
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		return nil, err
	}

	// Change directory to the cloned repository
	err = os.Chdir(getRepoName(gs.RepoURL))
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return nil, err
	}

	// Checkout the specified branch/tag
	err = checkoutBranch(gs.BranchTag)
	if err != nil {
		fmt.Println("Error checking out branch/tag:", err)
		return nil, err
	}

	// Command to fetch dependencies
	cmd := exec.Command("go", "mod", "graph")
	cmd.Env = append(cmd.Env, fmt.Sprintf("GO111MODULE=on"))

	// Execute the command
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// Convert output to string
	outputStr := string(output)

	// Split dependencies by newline
	dependencies := strings.Split(outputStr, "\n")

	// Create a map to store dependencies by their names
	depMap := make(map[string]*Dependency)

	// Process each dependency and build the dependency tree
	for _, dep := range dependencies {
		parts := strings.Split(dep, " ")
		if len(parts) < 2 {
			continue
		}
		parent := parts[0]
		child := parts[1]

		parentDep, ok := depMap[parent]
		if !ok {
			parentDep = &Dependency{Name: parent}
			depMap[parent] = parentDep
		}

		childDep, ok := depMap[child]
		if !ok {
			childDep = &Dependency{Name: child}
			depMap[child] = childDep
		}

		parentDep.Dependencies = append(parentDep.Dependencies, childDep)
	}

	// Find the root dependencies (ones without parents)
	var roots []*Dependency
	for _, dep := range depMap {
		if !hasParent(depMap, dep) {
			roots = append(roots, dep)
		}
	}

	// Convert the dependency tree to JSON
	jsonData, err := toJSON(roots)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return jsonData, nil
}

// toJSON converts the dependency tree to JSON, handling cyclic dependencies
func toJSON(roots []*Dependency) ([]byte, error) {
	visited := make(map[*Dependency]bool)
	return json.MarshalIndent(convertToJSON(roots, visited), "", "    ")
}

// convertToJSON converts the dependency tree to JSON recursively
func convertToJSON(roots []*Dependency, visited map[*Dependency]bool) []*DependencyJSON {
	var result []*DependencyJSON
	for _, dep := range roots {
		if visited[dep] {
			continue
		}
		visited[dep] = true
		depJSON := &DependencyJSON{Name: dep.Name, Version: dep.Version, Visited: visited}
		depJSON.Dependencies = convertToJSON(dep.Dependencies, visited)
		result = append(result, depJSON)
	}
	return result
}

func cloneRepository(repoURL string) error {

	url := repoURL
	directory := getRepoName(repoURL)

	// Clone the given repository to the given directory
	Info("git clone %s %s --recursive", url, directory)

	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	CheckIfError(err)

	return err
}

// Function to checkout the specified branch/tag
func checkoutBranch(branchTag string) error {
	cmd := exec.Command("git", "checkout", branchTag)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to checkout branch/tag: %s", err)
	}
	return nil
}

// Function to extract repository name from URL
func getRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}

// Function to check if a dependency has a parent
func hasParent(depMap map[string]*Dependency, dep *Dependency) bool {
	for _, d := range depMap {
		for _, child := range d.Dependencies {
			if child == dep {
				return true
			}
		}
	}
	return false
}

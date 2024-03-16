# Git Dependency Tree Generator

This is a simple Go program that generates a dependency tree for a Git repository based on its `go.mod` file. It clones the repository, checks out a specified branch/tag, and then analyzes the `go.mod` file to build the dependency tree. The dependency tree is then converted to JSON format and printed to the console.

## Usage

To use this program, follow these steps:

1. Ensure you have Go installed on your system.

2. Clone this repository:

   ```bash
   git clone https://github.com/example/git-dependency-tree-generator.git

3. Navigate to clone repository 
    ```bash
    cd git-dependency-tree-generator

4. Run the program with the following command and enter REPO_URL and BRANCH_TAG
    ```bash
    go run main.go 

5. The program will generate the dependency tree for the specified repository and branch/tag and print it to the console in JSON format.

## Dependencies

This program uses the following external dependencies:

- github.com/go-git/go-git/v5: Go package for accessing and manipulating Git repositories.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.

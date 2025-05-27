# Git Package

A Go package for interacting with GitHub's REST API. This package provides a simple interface to perform common Git operations like managing branches, files, and pull requests.

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/git
```

## Features

- Branch management (create/get)
- File operations (read/create/update/batch update)
- Pull request management (create/add reviewers)
- Token-based authentication
- Configurable API endpoints

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/pal-paul/go-libraries/pkg/git"
)

func main() {
    // Create a new Git client
    client := git.New(
        git.WithOwner("your-username"),
        git.WithRepo("your-repo"),
        git.WithToken("your-github-token"),
        git.WithContext(context.Background()),
    )

    // Get branch information
    branch, err := client.GetBranch("main")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Branch SHA: %s\n", branch.Object.Sha)

    // Create a new file
    content := []byte("Hello, World!")
    file, err := client.CreateUpdateAFile("main", "hello.txt", content, "Initial commit", "")
    if err != nil {
        panic(err)
    }
    fmt.Printf("File created: %s\n", file.Content.Name)
}
```

## API Reference

### Configuration

The package can be configured using various options:

```go
type Option func(*Config)

// Available options:
WithOwner(owner string)      // Set repository owner
WithRepo(repo string)        // Set repository name
WithToken(token string)      // Set GitHub access token
WithContext(ctx context.Context) // Set context for API requests
WithBaseURL(url string)      // Set custom API base URL
```

### Branch Operations

#### GetBranch

```go
GetBranch(branch string) (*BranchInfo, error)
```

Retrieves information about a specific branch in the repository.

- **Parameters**:
  - `branch`: The name of the branch to retrieve information for.
- **Returns**:
  - `*BranchInfo`: Contains branch details including its ref and SHA.
  - `error`: Any error that occurred during the operation.

#### CreateBranch

```go
CreateBranch(branch string, sha string) (*BranchInfo, error)
```

Creates a new branch in the repository.

- **Parameters**:
  - `branch`: The name of the new branch.
  - `sha`: The SHA of the commit the branch will point to.
- **Returns**:
  - `*BranchInfo`: Information about the created branch.
  - `error`: Any error that occurred during the operation.

### File Operations

#### GetAFile

```go
GetAFile(branch string, filePath string) (*FileInfo, error)
```

Retrieves information about a specific file in the repository.

- **Parameters**:
  - `branch`: The name of the branch containing the file.
  - `filePath`: The path to the file within the repository.
- **Returns**:
  - `*FileInfo`: File information including content and metadata.
  - `error`: Any error that occurred during the operation.

#### CreateUpdateAFile

```go
CreateUpdateAFile(branch string, filePath string, content []byte, message string, sha string) (*FileResponse, error)
```

Creates or updates a file in the repository.

- **Parameters**:
  - `branch`: The branch where the file will be created/updated.
  - `filePath`: The path to the file.
  - `content`: The file content.
  - `message`: The commit message.
  - `sha`: The file's current SHA (required for updates).
- **Returns**:
  - `*FileResponse`: Information about the created/updated file.
  - `error`: Any error that occurred during the operation.

#### CreateUpdateMultipleFiles

```go
CreateUpdateMultipleFiles(batch BatchFileUpdate) error
```

Updates or creates multiple files in a single commit.

- **Parameters**:
  - `batch`: A `BatchFileUpdate` struct containing:
    - `Branch`: Target branch name
    - `Message`: Commit message
    - `Files`: Array of `FileOperation` structs
- **Returns**:
  - `error`: Any error that occurred during the operation.

### Pull Request Operations

#### CreatePullRequest

```go
CreatePullRequest(baseBranch string, branch string, title string, description string) (int, error)
```

Creates a new pull request.

- **Parameters**:
  - `baseBranch`: The branch to merge into.
  - `branch`: The branch containing changes.
  - `title`: Pull request title.
  - `description`: Pull request description.
- **Returns**:
  - `int`: Pull request number.
  - `error`: Any error that occurred during the operation.

#### AddReviewers

```go
AddReviewers(number int, prReviewers Reviewers) error
```

Adds reviewers to a pull request.

- **Parameters**:
  - `number`: Pull request number.
  - `prReviewers`: A `Reviewers` struct containing:
    - `Users`: List of GitHub usernames
    - `Teams`: List of GitHub team names
- **Returns**:
  - `error`: Any error that occurred during the operation.

## Error Handling

The package returns meaningful errors for various scenarios:

- Invalid authentication
- Missing required parameters
- Network errors
- API rate limiting
- Invalid file operations
- Repository access issues

Example error handling:

```go
branch, err := client.GetBranch("nonexistent")
if err != nil {
    switch {
    case err.Error() == "branch not found":
        // Handle missing branch
    case strings.Contains(err.Error(), "401"):
        // Handle authentication error
    default:
        // Handle other errors
    }
}
```

## Testing

The package includes comprehensive tests. To run them:

```bash
go test -v ./...
```

For mock generation (requires `mockgen`):

```bash
go generate ./...
```

## Best Practices

1. **Token Security**: Never hardcode GitHub tokens. Use environment variables or secure configuration management.
2. **Error Handling**: Always check for errors and handle them appropriately.
3. **Resource Cleanup**: Use context with timeout for long-running operations.
4. **Rate Limiting**: Be mindful of GitHub API rate limits in production applications.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This package is released under the MIT License.

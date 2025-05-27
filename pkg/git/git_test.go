package git_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pal-paul/go-libraries/pkg/git"
	"github.com/stretchr/testify/assert"
)

func setupMockServer(t *testing.T, expectedPath string, method string, status int, response []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify required headers
		assert.Equal(t, "token test-token", r.Header.Get("Authorization"), "Authorization header mismatch")
		assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"), "Accept header mismatch")
		assert.Equal(t, expectedPath, r.URL.Path, "Path mismatch. Expected: %s, Got: %s", expectedPath, r.URL.Path)
		assert.Equal(t, method, r.Method, "Method mismatch. Expected: %s, Got: %s", method, r.Method)

		// Set required headers for the response
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(status)
		if response != nil {
			w.Write(response)
		}
	}))
}

func TestGitNew(t *testing.T) {
	tests := []struct {
		name      string
		opts      []git.Option
		wantError bool
	}{
		{
			name: "success with valid options",
			opts: []git.Option{
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := git.New(tt.opts...)
			assert.NotNil(t, client)
		})
	}
}

func TestGitGetBranch(t *testing.T) {
	tests := []struct {
		name      string
		branch    string
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:   "success - branch exists",
			branch: "main",
			response: []byte(`{
				"ref": "refs/heads/main",
				"node_id": "test-node",
				"url": "https://api.github.com/repos/test-owner/test-repo/git/refs/heads/main",
				"object": {
					"sha": "test-sha",
					"type": "commit",
					"url": "https://api.github.com/repos/test-owner/test-repo/git/commits/test-sha"
				}
			}`),
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:      "branch not found",
			branch:    "non-existent",
			status:    http.StatusNotFound,
			wantError: false,
		},
		{
			name:      "server error",
			branch:    "main",
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				fmt.Sprintf("/repos/%s/%s/git/refs/heads/%s", "test-owner", "test-repo", tt.branch),
				http.MethodGet,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			branchInfo, err := client.GetBranch(tt.branch)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.status == http.StatusOK {
					assert.NotNil(t, branchInfo)
					assert.Equal(t, "refs/heads/main", branchInfo.Ref)
					assert.Equal(t, "test-sha", branchInfo.Object.Sha)
				} else {
					assert.Nil(t, branchInfo)
				}
			}
		})
	}
}

func TestGitCreateBranch(t *testing.T) {
	tests := []struct {
		name      string
		branch    string
		sha       string
		status    int
		wantError bool
	}{
		{
			name:      "success",
			branch:    "feature",
			sha:       "test-sha",
			status:    http.StatusCreated,
			wantError: false,
		},
		{
			name:      "failure - server error",
			branch:    "feature",
			sha:       "test-sha",
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				"/repos/test-owner/test-repo/git/refs",
				http.MethodPost,
				tt.status,
				nil,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			_, err := client.CreateBranch(tt.branch, tt.sha)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitGetAFile(t *testing.T) {
	tests := []struct {
		name      string
		branch    string
		filePath  string
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:     "success",
			branch:   "main",
			filePath: "test.txt",
			response: []byte(`{
				"name": "test.txt",
				"path": "test.txt",
				"sha": "test-file-sha",
				"size": 100,
				"url": "https://api.github.com/repos/test-owner/test-repo/contents/test.txt",
				"html_url": "https://github.com/test-owner/test-repo/blob/main/test.txt",
				"git_url": "https://api.github.com/repos/test-owner/test-repo/git/blobs/test-file-sha",
				"content": "SGVsbG8gV29ybGQh",
				"encoding": "base64"
			}`),
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:      "file not found",
			branch:    "main",
			filePath:  "nonexistent.txt",
			status:    http.StatusNotFound,
			wantError: true,
		},
		{
			name:      "server error",
			branch:    "main",
			filePath:  "test.txt",
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				fmt.Sprintf("/repos/%s/%s/contents/%s", "test-owner", "test-repo", tt.filePath),
				http.MethodGet,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			fileContent, err := client.GetAFile(tt.branch, tt.filePath)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, fileContent)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fileContent)
				assert.Equal(t, "test.txt", fileContent.Name)
				assert.Equal(t, "test-file-sha", fileContent.Sha)
				assert.Equal(t, "base64", fileContent.Encoding)
			}
		})
	}
}

func TestGitCreateUpdateAFile(t *testing.T) {
	tests := []struct {
		name      string
		branch    string
		filePath  string
		content   string
		sha       string
		commitMsg string
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:      "create new file",
			branch:    "main",
			filePath:  "new-file.txt",
			content:   "Hello World!",
			commitMsg: "Create new file",
			response: []byte(`{
				"content": {
					"name": "new-file.txt",
					"path": "new-file.txt",
					"sha": "new-file-sha",
					"url": "https://api.github.com/repos/test-owner/test-repo/contents/new-file.txt"
				}
			}`),
			status:    http.StatusCreated,
			wantError: false,
		},
		{
			name:      "update existing file",
			branch:    "main",
			filePath:  "existing-file.txt",
			content:   "Updated content",
			commitMsg: "Update file",
			response: []byte(`{
				"content": {
					"name": "existing-file.txt",
					"path": "existing-file.txt",
					"sha": "updated-file-sha",
					"url": "https://api.github.com/repos/test-owner/test-repo/contents/existing-file.txt"
				}
			}`),
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:      "server error",
			branch:    "main",
			filePath:  "test.txt",
			content:   "test content",
			commitMsg: "Test commit",
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				fmt.Sprintf("/repos/%s/%s/contents/%s", "test-owner", "test-repo", tt.filePath),
				http.MethodPut,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			content := []byte(tt.content)
			response, err := client.CreateUpdateAFile(tt.branch, tt.filePath, content, tt.commitMsg, tt.sha)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.filePath, response.Content.Name)
			}
		})
	}
}

func TestGitCreatePullRequest(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		branch      string
		baseBranch  string
		response    []byte
		status      int
		wantError   bool
		expectedPR  int
	}{
		{
			name:        "success",
			title:       "Feature implementation",
			description: "Implementing new feature",
			branch:      "feature-branch",
			baseBranch:  "main",
			response:    []byte(`{"number": 101}`),
			status:      http.StatusCreated,
			wantError:   false,
			expectedPR:  101,
		},
		{
			name:        "validation error",
			title:       "",
			description: "Empty title",
			branch:      "feature",
			baseBranch:  "main",
			status:      http.StatusUnprocessableEntity,
			wantError:   true,
			expectedPR:  0,
		},
		{
			name:        "server error",
			title:       "Test PR",
			description: "Test body",
			branch:      "feature",
			baseBranch:  "main",
			status:      http.StatusInternalServerError,
			wantError:   true,
			expectedPR:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				"/repos/test-owner/test-repo/pulls",
				http.MethodPost,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			prNumber, err := client.CreatePullRequest(tt.baseBranch, tt.branch, tt.title, tt.description)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, 0, prNumber)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPR, prNumber)
			}
		})
	}
}

func TestGitAddReviewers(t *testing.T) {
	tests := []struct {
		name      string
		prNumber  int
		reviewers git.Reviewers
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:     "success - user reviewers",
			prNumber: 101,
			reviewers: git.Reviewers{
				Users: []string{"reviewer1", "reviewer2"},
			},
			response: []byte(`{
				"users": [
					{"login": "reviewer1"},
					{"login": "reviewer2"}
				]
			}`),
			status:    http.StatusCreated,
			wantError: false,
		},
		{
			name:     "success - team reviewers",
			prNumber: 101,
			reviewers: git.Reviewers{
				Teams: []string{"team1", "team2"},
			},
			response: []byte(`{
				"teams": [
					{"slug": "team1"},
					{"slug": "team2"}
				]
			}`),
			status:    http.StatusCreated,
			wantError: false,
		},
		{
			name:     "invalid reviewers",
			prNumber: 101,
			reviewers: git.Reviewers{
				Users: []string{"nonexistent-user"},
			},
			status:    http.StatusUnprocessableEntity,
			wantError: true,
		},
		{
			name:     "server error",
			prNumber: 101,
			reviewers: git.Reviewers{
				Users: []string{"reviewer1"},
			},
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				fmt.Sprintf("/repos/test-owner/test-repo/pulls/%d/requested_reviewers", tt.prNumber),
				http.MethodPost,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			err := client.AddReviewers(tt.prNumber, tt.reviewers)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitCreateUpdateMultipleFiles(t *testing.T) {
	tests := []struct {
		name      string
		batch     git.BatchFileUpdate
		status    int
		wantError bool
	}{
		{
			name: "success - multiple files",
			batch: git.BatchFileUpdate{
				Branch:  "feature",
				Message: "Update multiple files",
				Files: []git.FileOperation{
					{
						Path:    "file1.txt",
						Content: "content1",
					},
					{
						Path:    "file2.txt",
						Content: "content2",
					},
				},
			},
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name: "success - update with SHA",
			batch: git.BatchFileUpdate{
				Branch:  "main",
				Message: "Update with SHA",
				Files: []git.FileOperation{
					{
						Path:    "existing.txt",
						Content: "updated content",
						Sha:     "existing-sha",
					},
				},
			},
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name: "server error",
			batch: git.BatchFileUpdate{
				Branch:  "main",
				Message: "Failed update",
				Files: []git.FileOperation{
					{
						Path:    "test.txt",
						Content: "test content",
					},
				},
			},
			status:    http.StatusInternalServerError,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				"/repos/test-owner/test-repo/contents",
				http.MethodPut,
				tt.status,
				nil,
			)
			defer server.Close()

			client := git.New(
				git.WithOwner("test-owner"),
				git.WithRepo("test-repo"),
				git.WithToken("test-token"),
				git.WithContext(context.Background()),
				git.WithBaseURL(server.URL),
			)

			err := client.CreateUpdateMultipleFiles(tt.batch)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

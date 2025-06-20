package git

//go:generate mockgen -source=interface.go -destination=mocks/mock-git.go -package=mocks
import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type git struct {
	cfg *Config
}

func New(opts ...Option) *git {
	n := &git{cfg: defaultConfig()}
	for _, opt := range opts {
		opt(n.cfg)
	}
	return &git{
		cfg: n.cfg,
	}
}

// GetBranch retrieves information about a specific branch in the repository.
// Parameters:
//   - branch: The name of the branch to retrieve information for.
//
// Returns:
//   - A pointer to a BranchInfo struct containing branch details, or nil if the branch does not exist.
//   - An error if the request fails or if the response status is not 200 OK.
func (g *git) GetBranch(branch string) (*BranchInfo, error) {
	var branchInfo BranchInfo
	resp, err := g.get(
		"repos",
		fmt.Sprintf("%s/%s/git/refs/heads/%s", g.cfg.Owner, g.cfg.Repo, branch),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get branch %s: %s", branch, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &branchInfo)
	if err != nil {
		fmt.Println(string(body))
		return nil, err
	}
	return &branchInfo, nil
}

// CreateBranch creates a new branch in the repository with the specified name and SHA.
// Parameters:
//   - branch: The name of the new branch to create.
//   - sha: The SHA of the commit that the new branch will point to.
//
// Returns:
//   - A pointer to a BranchInfo struct containing information about the created branch, or nil if successful.
//   - An error if the request fails or if the response status is not 201 Created.
func (g *git) CreateBranch(branch string, sha string) (*BranchInfo, error) {
	reqBody := map[string]string{
		"ref": fmt.Sprintf("refs/heads/%s", branch),
		"sha": sha,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.post("repos", fmt.Sprintf("%s/%s/git/refs", g.cfg.Owner, g.cfg.Repo), nil, reqBodyJson)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to create a branch %s: %s", branch, resp.Status)
	}
	return nil, nil
}

// GetAFile retrieves information about a specific file in the repository at a given branch.
// Parameters:
//   - branch: The name of the branch where the file is located.
//   - filePath: The path to the file within the repository.
//
// Returns:
//   - A pointer to a FileInfo struct containing details about the file, or nil if the file does not exist.
//   - An error if the request fails or if the response status is not 200 OK.
func (g *git) GetAFile(branch string, filePath string) (*FileInfo, error) {
	var fileInfo FileInfo
	qs := url.Values{}
	qs.Add("ref", branch)
	resp, err := g.get("repos", fmt.Sprintf("%s/%s/contents/%s", g.cfg.Owner, g.cfg.Repo, filePath), qs)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file %s: %s", filePath, resp.Status)
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file %s: %s", filePath, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileInfo)
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

// CreateUpdateAFile creates or updates a file in the repository at a specified branch.
// Parameters:
//   - branch: The name of the branch where the file will be created or updated.
//   - filePath: The path to the file within the repository.
//   - content: The content of the file as a byte slice.
//   - message: The commit message for the file creation or update.
//   - sha: The SHA of the file if updating an existing file (optional).
//
// Returns:
//   - A pointer to a FileResponse struct containing details about the created or updated file.
//   - An error if the request fails or if the response status is not 201 Created.
func (g *git) CreateUpdateAFile(
	branch string,
	filePath string,
	content []byte,
	message string,
	sha string,
) (*FileResponse, error) {
	var fileResponse FileResponse
	b64content := b64.StdEncoding.EncodeToString(content)
	reqBody := map[string]string{
		"message": message,
		"content": b64content,
		"branch":  branch,
	}
	if sha != "" {
		reqBody["sha"] = sha
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.put(
		"repos",
		fmt.Sprintf("%s/%s/contents/%s", g.cfg.Owner, g.cfg.Repo, filePath),
		nil,
		reqBodyJson,
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 201 {
		return nil, fmt.Errorf("failed to update file %s: %s", filePath, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return nil, err
	}
	return &fileResponse, nil
}

// CreatePullRequest creates a pull request and returns the pull request number.
// Parameters:
//   - baseBranch: The name of the branch where the pull request will be merged into.
//   - branch: The name of the branch that contains the changes to be merged.
//   - title: The title of the pull request.
//   - description: The description of the pull request.
//
// Returns:
//   - The pull request number if successful, or an error if the request fails or if the response status is not 201 Created.
//   - An error if the request fails or if the response status is not 201 Created.
func (g *git) CreatePullRequest(
	baseBranch string,
	branch string,
	title string,
	description string,
) (int, error) {
	reqBody := map[string]any{
		"title":                 title,
		"body":                  description,
		"head":                  branch,
		"base":                  baseBranch,
		"maintainer_can_modify": true,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}
	resp, err := g.post("repos", fmt.Sprintf("%s/%s/pulls", g.cfg.Owner, g.cfg.Repo), nil, reqBodyJson)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 201 {
		return 0, fmt.Errorf("failed to create pull request: %s", resp.Status)
	}
	var pullResponse PullResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(body, &pullResponse)
	if err != nil {
		return 0, err
	}
	return pullResponse.Number, nil
}

// AddReviewers adds reviewers to a pull request.
// Parameters:
//   - number: The pull request number to which reviewers will be added.
//   - prReviewers: A Reviewers struct containing the list of users and teams to be added as reviewers.
//
// Returns:
//   - An error if the request fails or if the response status is not 201 Created.
//   - nil if the reviewers are successfully added.
func (g *git) AddReviewers(number int, prReviewers Reviewers) error {
	reqBody := make(map[string][]string)
	if len(prReviewers.Users) > 0 {
		reqBody["reviewers"] = prReviewers.Users
	}
	if len(prReviewers.Teams) > 0 {
		reqBody["team_reviewers"] = prReviewers.Teams
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp, err := g.post(
		"repos",
		fmt.Sprintf("%s/%s/pulls/%d/requested_reviewers", g.cfg.Owner, g.cfg.Repo, number),
		nil,
		reqBodyJson,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode == 422 {
		return fmt.Errorf("invalid reviewers: requesters are not collaborators")
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to add reviewers: %s", resp.Status)
	}
	return nil
}

type FileOperation struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Sha     string `json:"sha,omitempty"`
}

type BatchFileUpdate struct {
	Branch  string          `json:"branch"`
	Message string          `json:"message"`
	Files   []FileOperation `json:"files"`
}

// CreateUpdateMultipleFiles updates or creates multiple files in a repository branch using the Git Database API.
// This method creates blobs for each file, creates a new tree with the changes, creates a commit,
// and updates the branch reference to point to the new commit.
//
// Parameters:
//   - batch: A BatchFileUpdate struct containing the branch name, commit message,
//     and a list of files to be created or updated. Each file is represented by a
//     FileOperation struct, which includes the file path and content. The Sha field
//     is ignored for this method as it uses the Git Database API workflow.
//
// Returns:
// - An error if the operation fails, or nil if the files are successfully updated.
func (g *git) CreateUpdateMultipleFiles(batch BatchFileUpdate) error {
	// Step 1: Get the current branch reference to get the current commit SHA
	branchInfo, err := g.GetBranch(batch.Branch)
	if err != nil {
		return fmt.Errorf("failed to get branch %s: %w", batch.Branch, err)
	}

	currentCommitSha := branchInfo.Object.Sha

	// Step 2: Get the current tree from the current commit
	resp, err := g.get("repos", fmt.Sprintf("%s/%s/git/commits/%s", g.cfg.Owner, g.cfg.Repo, currentCommitSha), nil)
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get current commit: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read commit response: %w", err)
	}

	var commitInfo CommitResponse
	if err := json.Unmarshal(body, &commitInfo); err != nil {
		return fmt.Errorf("failed to parse commit response: %w", err)
	}

	currentTreeSha := commitInfo.Tree.Sha

	// Step 3: Create blobs for each file's content
	var treeEntries []TreeEntry
	for _, file := range batch.Files {
		// Create blob for file content
		blobReq := map[string]string{
			"content":  file.Content,
			"encoding": "utf-8",
		}
		blobReqJson, err := json.Marshal(blobReq)
		if err != nil {
			return fmt.Errorf("failed to marshal blob request for %s: %w", file.Path, err)
		}

		resp, err := g.post("repos", fmt.Sprintf("%s/%s/git/blobs", g.cfg.Owner, g.cfg.Repo), nil, blobReqJson)
		if err != nil {
			return fmt.Errorf("failed to create blob for %s: %w", file.Path, err)
		}
		if resp.StatusCode != 201 {
			return fmt.Errorf("failed to create blob for %s: %s", file.Path, resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read blob response for %s: %w", file.Path, err)
		}

		var blobResp BlobResponse
		if err := json.Unmarshal(body, &blobResp); err != nil {
			return fmt.Errorf("failed to parse blob response for %s: %w", file.Path, err)
		}

		// Add tree entry for this file
		treeEntries = append(treeEntries, TreeEntry{
			Path: file.Path,
			Mode: "100644", // Regular file mode
			Type: "blob",
			Sha:  blobResp.Sha,
		})
	}

	// Step 4: Create a new tree with the file changes
	treeReq := map[string]interface{}{
		"base_tree": currentTreeSha,
		"tree":      treeEntries,
	}
	treeReqJson, err := json.Marshal(treeReq)
	if err != nil {
		return fmt.Errorf("failed to marshal tree request: %w", err)
	}

	resp, err = g.post("repos", fmt.Sprintf("%s/%s/git/trees", g.cfg.Owner, g.cfg.Repo), nil, treeReqJson)
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create tree: %s", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read tree response: %w", err)
	}

	var treeResp TreeResponse
	if err := json.Unmarshal(body, &treeResp); err != nil {
		return fmt.Errorf("failed to parse tree response: %w", err)
	}

	// Step 5: Create a commit pointing to the new tree
	commitReq := map[string]interface{}{
		"message": batch.Message,
		"tree":    treeResp.Sha,
		"parents": []string{currentCommitSha},
	}
	commitReqJson, err := json.Marshal(commitReq)
	if err != nil {
		return fmt.Errorf("failed to marshal commit request: %w", err)
	}

	resp, err = g.post("repos", fmt.Sprintf("%s/%s/git/commits", g.cfg.Owner, g.cfg.Repo), nil, commitReqJson)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create commit: %s", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read commit response: %w", err)
	}

	var newCommitResp CommitResponse
	if err := json.Unmarshal(body, &newCommitResp); err != nil {
		return fmt.Errorf("failed to parse commit response: %w", err)
	}

	// Step 6: Update the branch reference to point to the new commit
	refReq := map[string]interface{}{
		"sha":   newCommitResp.Sha,
		"force": false, // Ensure this is a fast-forward update
	}
	refReqJson, err := json.Marshal(refReq)
	if err != nil {
		return fmt.Errorf("failed to marshal ref request: %w", err)
	}

	resp, err = g.patch("repos", fmt.Sprintf("%s/%s/git/refs/heads/%s", g.cfg.Owner, g.cfg.Repo, batch.Branch), nil, refReqJson)
	if err != nil {
		return fmt.Errorf("failed to update branch reference: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to update branch reference: %s", resp.Status)
	}

	return nil
}

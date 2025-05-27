package git

type IGit interface {
	GetBranch(branch string) (*BranchInfo, error)
	CreateBranch(branch string, sha string) (*BranchInfo, error)
	GetAFile(branch string, filePath string) (*FileInfo, error)
	CreateUpdateAFile(branch string, filePath string, content []byte, message string, sha string) (*FileResponse, error)
	CreatePullRequest(baseBranch string, branch string, title string, description string) (int, error)
	AddReviewers(number int, prReviewers Reviewers) error
	CreateUpdateMultipleFiles(batch BatchFileUpdate) error
}

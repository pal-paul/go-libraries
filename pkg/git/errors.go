package git

import (
	"fmt"
)

type ErrBranchNotFound struct {
	Value string
}

func (e ErrBranchNotFound) Error() string {
	return fmt.Sprintf("branch not found: %s", e.Value)
}

type ErrFailedToCreateBranch struct {
	Value string
}

func (e ErrFailedToCreateBranch) Error() string {
	return fmt.Sprintf("failed to create branch: %s", e.Value)
}

type ErrFailedToGetBranch struct {
	Value string
}

func (e ErrFailedToGetBranch) Error() string {
	return fmt.Sprintf("failed to get branch: %s", e.Value)
}

type ErrFailedToDeleteBranch struct {
	Value string
}

func (e ErrFailedToDeleteBranch) Error() string {
	return fmt.Sprintf("failed to delete branch: %s", e.Value)
}

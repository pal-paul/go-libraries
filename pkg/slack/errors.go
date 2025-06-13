package slack

import (
	"fmt"
	"time"
)

type ErrInvalidToken struct {
	Value string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token provided: %s", e.Value)
}

type ErrInvalidChannel struct {
	Value string
}

func (e *ErrInvalidChannel) Error() string {
	return fmt.Sprintf("invalid channel provided: %s", e.Value)
}

type ErrMessageNotFound struct {
	Value string
}

func (e *ErrMessageNotFound) Error() string {
	return fmt.Sprintf("message not found: %s", e.Value)
}

type ErrFileUploadFailed struct {
	Value string
}

func (e *ErrFileUploadFailed) Error() string {
	return fmt.Sprintf("file upload failed for: %s", e.Value)
}

type ErrInvalidReaction struct {
	Value string
}

func (e *ErrInvalidReaction) Error() string {
	return fmt.Sprintf("invalid reaction emoji: %s", e.Value)
}

type ErrUnauthorized struct {
	Value string
}

func (e *ErrUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized: invalid token %s", e.Value)
}

// ErrRateLimit represents the rate limit response from slack
type ErrRateLimit struct {
	Value time.Duration
}

func (e *ErrRateLimit) Error() string {
	return fmt.Sprintf("slack rate limit exceeded, retry after %s", e.Value)
}

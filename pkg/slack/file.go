package slack

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

func (s *slack) UploadFileWithContent(
	fileType string,
	fileName string,
	title string,
	content string,
	messageRef MessageRef,
) error {
	values := url.Values{}
	if fileType != "" {
		values.Add("filetype", fileType)
	}
	if fileName != "" {
		values.Add("filename", fileName)
	}
	if title != "" {
		values.Add("title", title)
	}
	if messageRef.Timestamp != "" {
		values.Add("thread_ts", messageRef.Timestamp)
	}
	if messageRef.Channel != "" {
		values.Add("channels", strings.Join([]string{messageRef.Channel}, ","))
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	var response SlackResponse
	if content != "" {
		values.Add("content", content)
		if resp, err := s.postForm("files.upload", headers, values); err != nil {
			return err
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(body, &response); err != nil {
				return err
			}
			if !response.Ok {
				return &ErrFileUploadFailed{
					Value: response.Error,
				}
			}
		}
	}
	return nil
}

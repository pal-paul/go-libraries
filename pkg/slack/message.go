package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func (s *slack) AddFormattedMessage(
	channel string,
	message Message,
) (messageRef MessageRef, err error) {
	message.Channel = channel
	var response SlackResponse

	apiEndpoint := "chat.postMessage"
	header := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	reqBody, err := json.Marshal(message)
	if err != nil {
		return messageRef, err
	}
	resp, err := s.postRequest(apiEndpoint, header, reqBody)
	if err != nil {
		return messageRef, fmt.Errorf("error post to slack: %v", err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return messageRef, err
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return messageRef, err
		}
		if !response.Ok {
			return messageRef, fmt.Errorf("error slack response")
		}
		messageRef.Channel = response.Channel
		messageRef.Timestamp = response.Ts
	}
	return messageRef, nil
}

func (s *slack) AddScheduleMessage(
	channel string,
	message Message,
	postAt int64,
) (messageRef MessageRef, err error) {
	message.Channel = channel
	var response SlackResponse

	apiEndpoint := fmt.Sprintf("%s/chat.scheduleMessage", s.cfg.BaseURL)
	header := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	reqBody, err := json.Marshal(message)
	if err != nil {
		return messageRef, err
	}
	resp, err := s.postRequest(apiEndpoint, header, reqBody)
	if err != nil {
		return messageRef, fmt.Errorf("error post to slack: %v", err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return messageRef, err
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return messageRef, err
		}
		if !response.Ok {
			return messageRef, fmt.Errorf("error slack response")
		}
		messageRef.Channel = response.Channel
		messageRef.Timestamp = response.Ts
	}
	return messageRef, nil
}

// Get getPermalink retrieves the permalink for a message in a channel.
func (m *slack) GetPermalink(channel string, messageRef MessageRef) (string, error) {
	apiEndpoint := fmt.Sprintf("%s/chat.getPermalink", baseUrl)
	header := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	values := url.Values{}
	if channel != "" {
		values.Add("channel", channel)
	}
	if messageRef.Timestamp != "" {
		values.Add("message_ts", messageRef.Timestamp)
	}
	resp, err := m.getRequest(apiEndpoint, header, values)
	if err != nil {
		return "", fmt.Errorf("error post to slack: %v", err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var response struct {
			Ok        bool   `json:"ok"`
			Permanent string `json:"permalink"`
			Error     string `json:"error"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		if !response.Ok {
			return "", fmt.Errorf("error slack response: %s", response.Error)
		}
		return response.Permanent, nil
	}
}

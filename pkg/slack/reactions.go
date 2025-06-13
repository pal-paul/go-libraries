package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// AddReaction adds a reaction emoji to a message
func (api *slack) AddReaction(name string, item MessageRef) (err error) {
	values := url.Values{}
	if name != "" {
		values.Set("name", name)
	}
	if item.Channel != "" {
		values.Set("channel", item.Channel)
	}
	if item.Timestamp != "" {
		values.Set("timestamp", item.Timestamp)
	}
	var response SlackResponse
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	if resp, err := api.postForm(fmt.Sprintf("%s/reactions.add", baseUrl), headers, values); err != nil {
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
			return fmt.Errorf("error slack response")
		}
	}
	return nil
}

// RemoveReactionContext removes a reaction emoji from a message.
func (api *slack) RemoveReaction(name string, item MessageRef) error {
	values := url.Values{}
	if name != "" {
		values.Set("name", name)
	}
	if item.Channel != "" {
		values.Set("channel", item.Channel)
	}
	if item.Timestamp != "" {
		values.Set("timestamp", item.Timestamp)
	}
	var response SlackResponse
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	if resp, err := api.postForm(fmt.Sprintf("%s/reactions.remove", baseUrl), headers, values); err != nil {
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
			return fmt.Errorf("error slack response")
		}
	}
	return nil
}

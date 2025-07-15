package shopee

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetAccountInfo() (map[string]interface{}, error) {
	body, err := c.DoRequest(RequestOptions{
		Method:   "GET",
		Endpoint: "account/basic/get_account_info",
	})
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid or Expired cookies")
	}

	uniqueIDStr := fmt.Sprintf("%.0f", data["shopid"].(float64))

	filtered := map[string]interface{}{
		"unique_id": uniqueIDStr,
		"name":      data["nickname"],
		"username":  data["username"],
		"email":     data["email"],
		"phone":     data["phone"],
	}
	return filtered, nil
}

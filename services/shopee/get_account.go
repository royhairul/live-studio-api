package shopee

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetAccount(userCookie string) (map[string]interface{}, error) {
	// 1. Siapkan endpoint dan client
	url := "https://shopee.co.id/api/v4/account/basic/get_account_info"
	client := &http.Client{}

	// 2. Buat request baru
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", userCookie)

	// 3. Kirim request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 4. Cek status response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(body))
	}

	// 5. Baca response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// 6. Parse JSON ke map
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// 7. Filter hasil jika perlu (contoh: hanya ambil "data")
	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format: missing 'data' field")
	}

	// Misal kamu hanya ingin ambil nama dan email saja
	filtered := map[string]interface{}{
		"name": data["nickname"],
		"username": data["username"],
		"email":    data["email"],
		"phone":    data["phone"], 
	}

	return filtered, nil
}

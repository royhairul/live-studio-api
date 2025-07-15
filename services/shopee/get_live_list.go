package shopee

import "strconv"

func (c *Client) GetLiveList(cookie string, page, pageSize int, timeDim, endDate string) ([]byte, error) {
	params := map[string]string{
		"page":     strconv.Itoa(page),
		"pageSize": strconv.Itoa(pageSize),
		"name":     "",
		"orderBy":  "",
		"sort":     "",
		"timeDim":  timeDim,
		"endDate":  endDate,
	}

	return c.DoRequest(RequestOptions{
		Method:      "GET",
		Endpoint:    "supply/api/lm/sellercenter/liveList/v2",
		QueryParams: params,
		Headers: map[string]string{
			"Cookie": cookie,
		},
	})
}

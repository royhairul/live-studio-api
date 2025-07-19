package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"go.uber.org/fx"
)

type ShopeeFinanceServiceDeps struct {
	fx.In
	ShopeeClient *shopee.ShopeeClient `name:"creatorShopeeClient"`
}

type ShopeeFinanceServiceImpl struct {
	ShopeeClient *shopee.ShopeeClient
}

func NewShopeeFinanceService(deps ShopeeFinanceServiceDeps) ShopeeFinanceService {
	return &ShopeeFinanceServiceImpl{ShopeeClient: deps.ShopeeClient}
}

func (s *ShopeeFinanceServiceImpl) GetShopeeLiveSalesReport(financeReq params.ShopeeLiveFinanceRequest, cookie string) (*params.ShopeeLiveFinanceResponse, error) {
	endpoint := "/supply/api/lm/sellercenter/liveList/v2"
	query := map[string]string{
		"page":     strconv.Itoa(financeReq.Page),
		"pageSize": strconv.Itoa(financeReq.PageSize),
		"name":     financeReq.Name,
		"orderBy":  financeReq.OrderBy,
		"sort":     financeReq.Sort,
		"timeDim":  financeReq.TimeDim,
		"endDate":  financeReq.EndDate,
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeLiveFinanceResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return &result.Data, nil
}

func (s *ShopeeFinanceServiceImpl) GetShopeeLiveSalesReportRange(financeReq params.ShopeeLiveFinanceRequest, cookie string) ([]params.ShopeeLiveReportItem, error) {
	var allItems []params.ShopeeLiveReportItem

	daysLeft, err := timeDimToDays(financeReq.TimeDim)
	if err != nil {
		return nil, fmt.Errorf("invalid timeDim format '%s': %w", financeReq.TimeDim, err)
	}

	endDate := financeReq.EndDate

	// ✅ Optimasi → kalau daysLeft <= 7 → gunakan 1x request "7d" → lalu filter manual
	if daysLeft <= 7 {
		req := financeReq
		req.TimeDim = "7d"
		req.EndDate = endDate

		resp, err := s.GetShopeeLiveSalesReport(req, cookie)
		if err != nil {
			return nil, err
		}

		log.Println(resp.List)
		log.Println(req.TimeDim)
		log.Println(resp.Total)

		return filterReportsByDays(resp.List, endDate, daysLeft)
	}

	for daysLeft > 0 {
		timeDim, duration := getBestTimeDim(daysLeft)
		if duration == 0 {
			return nil, fmt.Errorf("invalid timeDim generated during processing")
		}
		daysLeft -= duration

		req := financeReq
		req.TimeDim = timeDim
		req.EndDate = endDate

		resp, err := s.GetShopeeLiveSalesReport(req, cookie)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, resp.List...)

		// Geser endDate mundur
		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format: %w", err)
		}
		parsedEndDate = parsedEndDate.AddDate(0, 0, -duration)
		endDate = parsedEndDate.Format("2006-01-02")
	}

	return allItems, nil
}

// GetShopeeLiveRealTime implements ShopeeFinanceService.
func (s *ShopeeFinanceServiceImpl) GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItem, error) {
	req := params.ShopeeLiveFinanceRequest{
		Page:     1,
		PageSize: 20,
		TimeDim:  "1d",
		EndDate:  time.Now().Format("2006-01-02"),
	}

	resp, err := s.GetShopeeLiveSalesReport(req, cookie)
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

func getBestTimeDim(daysLeft int) (string, int) {
	switch {
	case daysLeft >= 30:
		return "30d", 30
	case daysLeft >= 15:
		return "15d", 15
	case daysLeft >= 7:
		return "7d", 7
	case daysLeft >= 1:
		return "1d", 1
	default:
		return "", 0
	}
}

func filterReportsByDays(reports []params.ShopeeLiveReportItem, endDate string, days int) ([]params.ShopeeLiveReportItem, error) {
	loc, _ := time.LoadLocation("Asia/Singapore")
	parsedEndDate, err := time.ParseInLocation("2006-01-02", endDate, loc)
	if err != nil {
		return nil, err
	}
	minDate := parsedEndDate.AddDate(0, 0, -(days - 1))

	var filtered []params.ShopeeLiveReportItem
	for _, item := range reports {
		itemTime := time.UnixMilli(item.StartTime).In(loc)
		if !itemTime.Before(minDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func timeDimToDays(timeDim string) (int, error) {
	trimmed := strings.TrimSuffix(timeDim, "d")
	return strconv.Atoi(trimmed)
}

func getDaysFromTimeDim(timeDim string) int {
	switch timeDim {
	case "30d":
		return 30
	case "15d":
		return 15
	case "7d":
		return 7
	case "1d":
		return 1
	default:
		return 0
	}
}

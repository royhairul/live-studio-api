package service

import (
	"fmt"

	"github.com/royhairul/live-studio-api/helpers/timehandler"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/params"

	performaparams "github.com/royhairul/live-studio-api/internal/domains/performa/params"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	accountadsservice "github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

type DashboardServiceImpl struct {
	// TODO: add repository dependency
	hostSvc           hostservice.HostService
	attendanceSvc     attendanceservice.AttendanceService
	accountSessionSvc accountsessionservice.AccountsessionService
	accountSvc        accountservice.AccountService
	transactionSvc    transactionservice.TransactionService
	studioSvc         studioservice.StudioService
	accountAdsSvc     accountadsservice.AccountadsService
}

func NewDashboardService(
	hostSvc hostservice.HostService,
	attendanceSvc attendanceservice.AttendanceService,
	accountSessionSvc accountsessionservice.AccountsessionService,
	accountSvc accountservice.AccountService,
	transactionSvc transactionservice.TransactionService,
	accountAdsSvc accountadsservice.AccountadsService,
	studioSvc studioservice.StudioService,
) DashboardService {
	return &DashboardServiceImpl{
		hostSvc,
		attendanceSvc,
		accountSessionSvc,
		accountSvc,
		transactionSvc,
		studioSvc,
		accountAdsSvc,
	}
}

// DashboardAdmin implements DashboardService.
func (d *DashboardServiceImpl) DashboardAdmin(startDate string, endDate string) (*params.DashboardResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous period dihitung mundur dengan panjang hari yang sama
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get Attendances (current + previous)
	attendances, err := d.attendanceSvc.WithDateRange(*start, *end).FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	prevAttendances, err := d.attendanceSvc.WithDateRange(prevStart, prevEnd).FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get prev attendances: %v", err)
	}

	// Get All Studio
	studios, err := d.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	list := []performaparams.PerformaStudioItemResponse{}

	var currGMV, prevGMV int64
	var currCommissionPaid, prevCommissionPaid int64
	var currCommissionPending, prevCommissionPending int64
	var currAds, prevAds uint
	var currIncome, prevIncome int64

	for _, studio := range studios {
		// Reset per studio
		var studioGMV, studioPrevGMV int64
		var studioCommissionPaid, studioPrevCommissionPaid int64
		var studioCommissionPending, studioPrevCommissionPending int64
		var studioAds, studioPrevAds uint
		var studioIncome, studioPrevIncome int64

		// Get Account in this studio
		accounts, err := d.accountSvc.WithStudioID(fmt.Sprintf("%d", studio.ID)).FindAll()
		if err != nil {
			return nil, err
		}

		for _, account := range accounts {
			// current transactions
			transactions, err := d.transactionSvc.
				WithAccountID(fmt.Sprintf("%d", account.ID)).
				WithDate(*start, *end).
				FindAll()
			if err != nil {
				return nil, err
			}
			for _, tx := range transactions {
				studioCommissionPaid += int64(tx.Commission.Paid)
				studioCommissionPending += int64(tx.Commission.Pending)
			}

			// previous transactions
			prevTransactions, err := d.transactionSvc.
				WithAccountID(fmt.Sprintf("%d", account.ID)).
				WithDate(prevStart, prevEnd).
				FindAll()
			if err != nil {
				return nil, err
			}
			for _, tx := range prevTransactions {
				studioPrevCommissionPaid += int64(tx.Commission.Paid)
				studioPrevCommissionPending += int64(tx.Commission.Pending)
			}

			// ads
			allAds, err := d.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range allAds {
				studioAds += a.Ads
			}

			prevAllAds, err := d.accountAdsSvc.FindByDateAndAccounts(&prevStart, &prevEnd, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range prevAllAds {
				studioPrevAds += a.Ads
			}
		}

		// GMV current
		for _, att := range attendances {
			accountsessions, err := d.accountSessionSvc.WithAttendanceID(fmt.Sprintf("%d", att.ID)).FindAll()
			if err != nil {
				return nil, err
			}
			for _, session := range accountsessions {
				studioGMV += int64(session.GMVPaid)
			}
		}

		// GMV previous
		for _, att := range prevAttendances {
			accountsessions, err := d.accountSessionSvc.WithAttendanceID(fmt.Sprintf("%d", att.ID)).FindAll()
			if err != nil {
				return nil, err
			}
			for _, session := range accountsessions {
				studioPrevGMV += int64(session.GMVPaid)
			}
		}

		studioIncome = (studioCommissionPaid + studioCommissionPending) - int64(studioAds)
		studioPrevIncome = (studioPrevCommissionPaid + studioPrevCommissionPending) - int64(studioPrevAds)

		// Tambahkan ke list
		list = append(list, performaparams.PerformaStudioItemResponse{
			StudioID:   fmt.Sprintf("%d", studio.ID),
			StudioName: studio.Name,
			Commission: studioCommissionPaid + studioCommissionPending,
			GMV:        studioGMV,
			Ads:        int64(studioAds),
			Income:     studioIncome,
		})

		// Akumulasi ke total metrics
		currGMV += studioGMV
		prevGMV += studioPrevGMV
		currCommissionPaid += studioCommissionPaid
		prevCommissionPaid += studioPrevCommissionPaid
		currCommissionPending += studioCommissionPending
		prevCommissionPending += studioPrevCommissionPending
		currAds += studioAds
		prevAds += studioPrevAds
		currIncome += studioIncome
		prevIncome += studioPrevIncome
	}

	accounts, err := d.accountSvc.FindAll()
	if err != nil {
		return nil, err
	}

	hosts, err := d.hostSvc.FindAll()
	if err != nil {
		return nil, err
	}

	results := &params.DashboardResponse{
		CurrentPeriod: params.PeriodInfo{
			Start: timehandler.FormatDate(start),
			End:   timehandler.FormatDate(end),
			Days:  days,
		},
		PreviousPeriod: params.PeriodInfo{
			Start: timehandler.FormatDate(&prevStart),
			End:   timehandler.FormatDate(&prevEnd),
			Days:  days,
		},
		Metrics: params.Metrics{
			Commission: NewMetric((currCommissionPaid + currCommissionPending), (prevCommissionPaid + prevCommissionPending)),
			GMV:        NewMetric(currGMV, prevGMV),
			Ads:        NewMetric(int64(currAds), int64(prevAds)),
			Income:     NewMetric(currIncome, prevIncome),
			Account:    int64(len(accounts)),
			Host:       int64(len(hosts)),
		},
		List: list,
	}

	return results, nil
}

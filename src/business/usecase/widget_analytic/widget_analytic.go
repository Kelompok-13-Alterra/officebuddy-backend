package widget_analytic

import (
	"context"
	"errors"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	officeDom "go-clean/src/business/domain/office"
	ratingDom "go-clean/src/business/domain/rating"
	transactionDom "go-clean/src/business/domain/transaction"
	userDom "go-clean/src/business/domain/user"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetDashboardWidget(ctx context.Context) (entity.DashboardWidgetResult, error)
	GetOfficeWidget(ctx context.Context, param entity.OfficeWidgetParam) (entity.OfficeWidgetResult, error)
	GetRevenueWidget(ctx context.Context) (entity.RevenueWidgetResult, error)
	GetUserWidget(ctx context.Context) (entity.UserWidgetResult, error)
}

type widgetAnalytic struct {
	office              officeDom.Interface
	transaction         transactionDom.Interface
	midtransTransaction midtransTransactionDom.Interface
	rating              ratingDom.Interface
	user                userDom.Interface
}

func Init(od officeDom.Interface, td transactionDom.Interface, rd ratingDom.Interface, mtd midtransTransactionDom.Interface, ud userDom.Interface) Interface {
	w := &widgetAnalytic{
		office:              od,
		transaction:         td,
		midtransTransaction: mtd,
		rating:              rd,
		user:                ud,
	}
	return w
}

func (wa *widgetAnalytic) GetDashboardWidget(ctx context.Context) (entity.DashboardWidgetResult, error) {
	result := entity.DashboardWidgetResult{}

	officeCount, err := wa.office.GetCount(entity.OfficeParam{
		Type: entity.OfficeType,
	})
	if err != nil {
		return result, err
	}
	result.OfficeTotal = officeCount

	coworkingCount, err := wa.office.GetCount(entity.OfficeParam{
		Type: entity.CoWorkingType,
	})
	if err != nil {
		return result, err
	}
	result.CoWorkingTotal = coworkingCount

	transactions, err := wa.transaction.GetTransactionToday()
	if err != nil {
		return result, err
	}

	officeIDsMap := make(map[uint]bool)
	officeIDs := []uint{}
	for _, t := range transactions {
		officeIDsMap[t.OfficeID] = true
	}
	for k := range officeIDsMap {
		officeIDs = append(officeIDs, k)
	}

	offices, err := wa.office.GetListByID(officeIDs)
	if err != nil {
		return result, err
	}

	officesMap := make(map[uint]entity.Office)
	for _, o := range offices {
		officesMap[o.ID] = o
	}

	coworkingCountTrx := 0
	officeCountTrx := 0
	for _, t := range transactions {
		if officesMap[t.OfficeID].Type == entity.CoWorkingType {
			coworkingCountTrx++
		} else if officesMap[t.OfficeID].Type == entity.OfficeType {
			officeCountTrx++
		}
	}

	result.CoWorkingTransactionToday = coworkingCountTrx
	result.OfficeTransactionToday = officeCountTrx

	return result, nil
}

func (wa *widgetAnalytic) GetOfficeWidget(ctx context.Context, param entity.OfficeWidgetParam) (entity.OfficeWidgetResult, error) {
	result := entity.OfficeWidgetResult{}

	if param.Type != "coworking" && param.Type != "office" {
		return result, errors.New("tipe office tidak tersedia")
	}

	countOffice, err := wa.office.GetCount(entity.OfficeParam{
		Type: param.Type,
	})
	if err != nil {
		return result, err
	}
	result.OfficeCount = countOffice

	offices, err := wa.office.GetList(entity.OfficeParam{
		Type: param.Type,
	})
	if err != nil {
		return result, err
	}

	officeIDsMap := make(map[uint]bool)
	officeIDs := []uint{}
	for _, o := range offices {
		officeIDsMap[o.ID] = true
		officeIDs = append(officeIDs, o.ID)
	}

	transactions, err := wa.transaction.GetTransactionToday()
	if err != nil {
		return result, err
	}

	transactionCount := 0
	for _, t := range transactions {
		if _, ok := officeIDsMap[t.OfficeID]; ok {
			transactionCount++
		}
	}
	result.TotalBooking = transactionCount

	ratingCount, err := wa.rating.GetCountInByID(officeIDs)
	if err != nil {
		return result, err
	}
	result.TotalRating = ratingCount

	return result, nil
}

func (wa *widgetAnalytic) GetRevenueWidget(ctx context.Context) (entity.RevenueWidgetResult, error) {
	result := entity.RevenueWidgetResult{}

	todayTransactions, err := wa.midtransTransaction.GetMidtransTransactionToday()
	if err != nil {
		return result, err
	}

	totalTodayAmount := 0
	for _, t := range todayTransactions {
		totalTodayAmount += t.Amount
	}
	result.TotalTodayRevenue = int64(totalTodayAmount)

	allTransactions, err := wa.midtransTransaction.GetList(entity.MidtransTransactionParam{Status: entity.StatusSuccess})
	if err != nil {
		return result, err
	}

	totalAllAmount := 0
	for _, t := range allTransactions {
		totalAllAmount += t.Amount
	}
	result.TotalAllRevenue = int64(totalAllAmount)

	return result, err
}

func (wa *widgetAnalytic) GetUserWidget(ctx context.Context) (entity.UserWidgetResult, error) {
	result := entity.UserWidgetResult{}

	userCount, err := wa.user.GetCount()
	if err != nil {
		return result, err
	}
	result.TotalUser = userCount

	ratingCount, err := wa.rating.GetCount()
	if err != nil {
		return result, err
	}
	result.TotalRating = ratingCount

	return result, nil
}

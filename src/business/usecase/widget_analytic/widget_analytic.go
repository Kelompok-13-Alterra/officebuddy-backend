package widget_analytic

import (
	"context"
	officeDom "go-clean/src/business/domain/office"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetDashboardWidget(ctx context.Context) (entity.DashboardWidgetResult, error)
}

type widgetAnalytic struct {
	office      officeDom.Interface
	transaction transactionDom.Interface
}

func Init(od officeDom.Interface, td transactionDom.Interface) Interface {
	w := &widgetAnalytic{
		office:      od,
		transaction: td,
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

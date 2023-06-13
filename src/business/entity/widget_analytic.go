package entity

type DashboardWidgetResult struct {
	OfficeTotal               int64
	CoWorkingTotal            int64
	OfficeTransactionToday    int
	CoWorkingTransactionToday int
}

type OfficeWidgetResult struct {
	OfficeCount  int64
	TotalBooking int
	TotalRating  int64
}

type RevenueWidgetResult struct {
	TotalTodayRevenue int64
	TotalAllRevenue   int64
}

type OfficeWidgetParam struct {
	Type string `form:"type" binding:"required"`
}

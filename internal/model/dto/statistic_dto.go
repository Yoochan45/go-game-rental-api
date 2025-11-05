package dto

type DashboardStats struct {
	TotalUsers           int64   `json:"total_users"`
	TotalPartners        int64   `json:"total_partners"`
	TotalGames           int64   `json:"total_games"`
	TotalBookings        int64   `json:"total_bookings"`
	TotalRevenue         float64 `json:"total_revenue"`
	PendingApprovals     int64   `json:"pending_approvals"`
	ActiveBookings       int64   `json:"active_bookings"`
	PendingDisputes      int64   `json:"pending_disputes"`
	MonthlyGrowthUsers   float64 `json:"monthly_growth_users"`
	MonthlyGrowthRevenue float64 `json:"monthly_growth_revenue"`
}

type PartnerStats struct {
	TotalGames        int64   `json:"total_games"`
	ActiveGames       int64   `json:"active_games"`
	TotalBookings     int64   `json:"total_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageRating     float64 `json:"average_rating"`
	PendingApprovals  int64   `json:"pending_approvals"`
}

type CustomerStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	TotalSpent        float64 `json:"total_spent"`
	FavoriteCategory  string  `json:"favorite_category,omitempty"`
}

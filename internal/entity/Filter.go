package entity

type Filter struct {
	Query      *string
	Page       *int
	PerPage    *int
	SortBy     *string
	SortOrder  *string
	StartDate  *string
	EndDate    *string
	CategoryID *string
	UserID     *string
}

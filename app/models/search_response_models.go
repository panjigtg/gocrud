package models

type MetaInfo struct {
	Page 	int 	`json:"page"`
	Limit 	int		`json:"limit"`
	Total 	int 	`json:"total"`
	Pages 	int 	`json:"pages"`
	SortBy 	string 	`json:"sortBy"`
	Order 	string 	`json:"order"`
	Search 	string 	`json:"search"`
}

type UserResponse struct {
	Data []FilterUsers	`json:"data"`
	Meta MetaInfo 	`json:"meta"`
}

type AlumniResponse struct {
	Data []FilterAlumni	`json:"data"`
	Meta MetaInfo 	`json:"meta"`
}

type PekerjaanResponse struct {
	Data []FilterPekerjaan	`json:"data"`
	Meta MetaInfo 	`json:"meta"`
}
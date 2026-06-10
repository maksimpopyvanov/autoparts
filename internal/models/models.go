package models

import "time"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Part struct {
	ID          int        `json:"id"`
	CategoryID  *int       `json:"category_id"`
	Name        string     `json:"name"`
	Article     *string    `json:"article"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PartWithCategory struct {
	Part
	CategoryName *string `json:"category_name"`
	Brands       []Brand `json:"brands"`
}

type Stock struct {
	PartID       int     `json:"part_id"`
	PartName     string  `json:"part_name"`
	Article      *string `json:"article"`
	CategoryName *string `json:"category_name"`
	Quantity     int     `json:"quantity"`
}

type Income struct {
	ID       int     `json:"id"`
	PartID   int     `json:"part_id"`
	PartName string  `json:"part_name"`
	Quantity int     `json:"quantity"`
	Date     string  `json:"date"`
	Comment  *string `json:"comment"`
}

type Outcome struct {
	ID       int     `json:"id"`
	PartID   int     `json:"part_id"`
	PartName string  `json:"part_name"`
	Quantity int     `json:"quantity"`
	Date     string  `json:"date"`
	Comment  *string `json:"comment"`
}

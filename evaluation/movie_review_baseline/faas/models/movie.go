package models

type Movie struct {
	Id       string
	Title    string
	Year     int
	Category string
}

type CategoryListItem struct {
	Id    string
	Title string
}

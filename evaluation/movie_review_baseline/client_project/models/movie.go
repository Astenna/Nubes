package models

type Movie struct {
	Id             string
	Title          string
	ProductionYear int
	Category       string
}

type CategoryListItem struct {
	Id    string
	Title string
}

type MoviesOfCategoryTemplateInput struct {
	Name   string
	Movies []CategoryListItem
}

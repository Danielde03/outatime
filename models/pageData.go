package models

// All page data structs will be of type PageData, so the data can be rendered in a template.
type PageData interface {
	// Display the contents of the data
	Display()
}

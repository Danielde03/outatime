package util

import (
	"net/http"
	"outatime/models"
	"text/template"
)

// Render the layout and the chosen template
//
// page is the name of the page template and pass the data to be used in the page in pageData.
//
// logged in to see what nav bar to use. If not, use standard nav bar.
//
// Return error if an error if one appears in parsing or executing templates.
func RenderTemplate(res http.ResponseWriter, page string, loggedIn bool, pageData *models.PageData) error {

	t, err := template.ParseFiles("templates/layout.html", "templates/include/"+page+".html")

	// parse nav bar
	if loggedIn {
		_, err := t.ParseFiles("templates/include/navbars/userNav.html")

		if err != nil {
			return err
		}
	} else {
		_, err := t.ParseFiles("templates/include/navbars/stndNav.html")

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	err = t.Execute(res, pageData)

	if err != nil {
		return err
	}

	return nil
}

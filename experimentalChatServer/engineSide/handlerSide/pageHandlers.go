package handlerSide

import (
	"github.com/CloudyKit/jet/v6"
	"log"
	"net/http"
)

var jetView = jet.NewSet(
	jet.NewOSFileSystemLoader("./HTML"),
	jet.InDevelopmentMode(), // Removes the need to re-start the application. Take it out before the production.
)

func HomePage(responseWriter http.ResponseWriter, pointerToRequest *http.Request) {
	renderPageError := pageRenderer(responseWriter, "homePageWithTemplate.jet", nil)

	if renderPageError != nil {
		log.Println(renderPageError)
	}
}

func pageRenderer(responseWriter http.ResponseWriter, templateToRender string, templateData jet.VarMap) error {
	pageView, pageError := jetView.GetTemplate(templateToRender)

	if pageError != nil {
		log.Println(pageError)
		return pageError
	}

	pageError = pageView.Execute(responseWriter, templateData, nil)
	if pageError != nil {
		log.Println(pageError)
		return pageError
	}

	return nil
}

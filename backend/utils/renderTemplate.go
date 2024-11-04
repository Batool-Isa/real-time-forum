package utils

import (
	"html/template"
	"net/http"
	"log"
)

type Data struct {
	Errors     map[string]string
	DataPassed interface{}
}

// Function to render HTML files
// Function to render HTML files
func RenderTemplate(w http.ResponseWriter, r *http.Request, templName string, data interface{}, errors ...map[string]string) error {
    temp, err := template.ParseFiles("template/" + templName)
    if err != nil {
        log.Println("Error parsing template:", err)
        http.Error(w, "Unable to render template", http.StatusInternalServerError)
        return err
    }

    templataData := Data{
        DataPassed: data,
    }

    // Check for errors
    if len(errors) > 0 {
        templataData.Errors = errors[0]
    } else {
        templataData.Errors = nil
    }

    // Execute template
    err = temp.Execute(w, templataData)
    if err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Error Executing Template", http.StatusInternalServerError)
        return err
    }

    return nil // Return nil if everything is successful
}

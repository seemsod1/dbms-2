package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"laboratory_databases_2/internal/config"
	"laboratory_databases_2/internal/models"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
)

var functions = template.FuncMap{
	"iterate": Iterate,
	"add":     Add,
}

var app *config.AppConfig
var pathToTemplates = "./templates"

// Add adds two numbers together
func Add(a, b any) any {
	switch aType := a.(type) {
	case int:
		switch bType := b.(type) {
		case int:
			return aType + bType
		case uint:
			return int(aType) + int(bType)
		default:
			log.Printf("Unsupported type for 'b': %v", reflect.TypeOf(b))
		}
	case uint:
		switch bType := b.(type) {
		case uint:
			return aType + bType
		case int:
			return int(aType) + bType
		default:
			log.Printf("Unsupported type for 'b': %v", reflect.TypeOf(b))
		}
	default:
		log.Printf("Unsupported type for 'a': %v", reflect.TypeOf(a))
	}

	return nil
}

// Iterate creates a slice of integers starting at 0 and ending at count
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("error rendering the template: ", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

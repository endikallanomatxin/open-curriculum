package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"text/template"

	structtomap "github.com/Klathmon/StructToMap"
)

// parseAcceptLanguage returns the best match for the languages supported by your application.
func ParseAcceptLanguage(acceptLangHeader string, supportedLangs []string) string {
	// Default language
	defaultLang := "en"

	// Parse the Accept-Language header and split by comma
	options := strings.Split(acceptLangHeader, ",")

	// Map to store priority
	langQuality := make(map[string]float64)
	for _, option := range options {
		parts := strings.Split(strings.TrimSpace(option), ";q=")
		lang := strings.Split(parts[0], "-")[0] // Normalize to base language code
		quality := 1.0                          // Default quality is 1.0

		if len(parts) == 2 {
			fmt.Sscanf(parts[1], "%f", &quality)
		}

		// Store the highest quality found for each language
		if existingQuality, exists := langQuality[lang]; !exists || quality > existingQuality {
			langQuality[lang] = quality
		}
	}

	// Find the supported language with the highest quality
	highestQuality := -1.0
	var selectedLang string
	for _, lang := range supportedLangs {
		if quality, exists := langQuality[lang]; exists && quality > highestQuality {
			highestQuality = quality
			selectedLang = lang
		}
	}

	if selectedLang == "" {
		return defaultLang
	}

	return selectedLang
}

func SetLanguageCookie(w http.ResponseWriter, r *http.Request) {
	// Recuperar el código de idioma enviado por el usuario, por ejemplo 'en', 'es', etc.
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		http.Error(w, "Language parameter is missing", http.StatusBadRequest)
		return
	}

	redirectURL := r.URL.Query().Get("redirect")

	// Crear una nueva cookie con el código de idioma
	cookie := http.Cookie{
		Name:   "lang",
		Value:  lang,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 7, // 1 semana
	}

	// Añadir la cookie a la respuesta
	http.SetCookie(w, &cookie)

	// Redirigir al usuario a la página de la que vino
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	supportedLanguages := []string{"en", "es", "eu"} // Update this list based on your available languages
	var lang string

	langCookie, err := r.Cookie("lang")
	if err != nil {
		acceptLang := r.Header.Get("Accept-Language")
		lang = ParseAcceptLanguage(acceptLang, supportedLanguages)
	} else {
		lang = langCookie.Value
	}

	// Normalize language code (example: convert "en-US" to "en-us")
	lang = strings.ToLower(lang)

	var tmpl_lang string
	if lang == "en" {
		tmpl_lang = tmpl
	} else {
		tmpl_lang = tmpl[:len(tmpl)-5] + "-" + lang + ".html"
	}

	t, err := template.ParseFiles("web/templates/base.html", "web/templates/"+tmpl_lang)
	if err != nil {
		// Get the english version of the template
		t, err = template.ParseFiles("web/templates/base.html", "web/templates/"+tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var newdata map[string]interface{}

	if data == nil {
		newdata = make(map[string]interface{})
		// Else if data is a struct, create a new struct that contains the original data plus the URL
	} else if reflect.TypeOf(data).Kind() == reflect.Struct {
		newdata, _ = structtomap.Convert(data)
	} else if _, ok := data.(map[string]interface{}); ok {
		newdata = data.(map[string]interface{})
	} else {
		http.Error(w, "Data type not supported", http.StatusInternalServerError)
		return
	}
	// Add a new field to the data struct
	// URL
	newdata["URL"] = r.URL.Path

	err = t.ExecuteTemplate(w, "base.html", newdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

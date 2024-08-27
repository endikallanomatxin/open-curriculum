package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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

// typeIs checks if the type of an object is exactly the provided type name.
func typeIs(typeName string, i interface{}) bool {
	return reflect.TypeOf(i).String() == typeName
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, block interface{}) {
	supportedLanguages := []string{"en", "es", "eu"} // Update this list based on your available languages
	var lang string

	// Get the language from the cookie or Accept-Language header
	langCookie, err := r.Cookie("lang")
	if err != nil {
		acceptLang := r.Header.Get("Accept-Language")
		lang = ParseAcceptLanguage(acceptLang, supportedLanguages)
	} else {
		lang = langCookie.Value
	}

	// Normalize language code (example: convert "en-US" to "en-us")
	lang = strings.ToLower(lang)

	var tmplLang string
	if lang == "en" {
		tmplLang = tmpl
	} else {
		tmplLang = tmpl[:len(tmpl)-5] + "-" + lang + ".html"
	}

	// Add custom functions to the template
	funcs := template.FuncMap{
		"typeIs": typeIs,
	}

	// Try to parse the template files with custom functions
	t, err := template.New("base.html").Funcs(funcs).ParseFiles("web/templates/base.html", "web/templates/"+tmplLang)
	if err != nil {
		// Fallback to the default template if language-specific one fails
		t, err = template.New("base.html").Funcs(funcs).ParseFiles("web/templates/base.html", "web/templates/"+tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Prepare the data for the template
	var newdata map[string]interface{}
	if data == nil {
		newdata = make(map[string]interface{})
	} else {
		switch reflect.TypeOf(data).Kind() {
		case reflect.Struct:
			newdata, err = structtomap.Convert(data)
			if err != nil {
				log.Printf("Error converting struct to map: %v", err)
				http.Error(w, "Failed to convert struct to map", http.StatusInternalServerError)
				return
			}
		case reflect.Map:
			if m, ok := data.(map[string]interface{}); ok {
				newdata = m
			} else {
				http.Error(w, "Data type not supported", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Unsupported data type", http.StatusInternalServerError)
			return
		}
	}
	newdata["URL"] = r.URL.Path

	// Render the specific block if provided
	if block != nil && block != "" {
		if blockName, ok := block.(string); ok {
			err = t.ExecuteTemplate(w, blockName, newdata)
			if err != nil {
				log.Printf("Error executing template block %s: %v", blockName, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			http.Error(w, "Invalid block type", http.StatusInternalServerError)
			return
		}
	}

	// Render the base template
	err = t.Execute(w, newdata)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SetActiveProposalID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("active_proposal_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid proposal ID", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:   "active_proposal_id",
		Value:  strconv.FormatInt(id, 10),
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 7, // 1 week
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/teach", http.StatusFound)
}

func GetActiveProposalID(r *http.Request) int64 {
	activeProposalIDCookie, err := r.Cookie("active_proposal_id")
	if err != nil {
		return 0
	}

	activeProposalIDStr := activeProposalIDCookie.Value

	activeProposalID, err := strconv.ParseInt(activeProposalIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting active_proposal_id to int:", err)
		return 0
	}
	return int64(activeProposalID)
}

func SetOpenUnit(w http.ResponseWriter, r *http.Request) {

	isProposed := r.URL.Query().Get("is_proposed")
	isProposedCookie := http.Cookie{
		Name:   "open_unit_is_proposed",
		Value:  isProposed,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 7, // 1 week
	}
	http.SetCookie(w, &isProposedCookie)

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		fmt.Println("Invalid ID")
		return
	}
	idCookie := http.Cookie{
		Name:   "open_unit_id",
		Value:  strconv.FormatInt(id, 10),
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 7, // 1 week
	}
	http.SetCookie(w, &idCookie)

	if isProposed == "none" {
		http.Redirect(w, r, "/learn", http.StatusFound)
	} else {
		http.Redirect(w, r, "/teach", http.StatusFound)
	}
}

func GetOpenUnit(r *http.Request) (bool, int64) {
	openUnitIsProposedCookie, err := r.Cookie("open_unit_is_proposed")
	if err != nil {
		return false, 0
	}
	openUnitIsProposedStr := openUnitIsProposedCookie.Value
	var openUnitIsProposed bool
	if openUnitIsProposedStr == "true" {
		openUnitIsProposed = true
	} else if openUnitIsProposedStr == "false" {
		openUnitIsProposed = false
	} else if openUnitIsProposedStr == "none" {
		openUnitIsProposed = false
	} else {
		fmt.Println("Invalid open_unit_is_proposed value:", openUnitIsProposedStr)
		return false, 0
	}

	openUnitIDCookie, err := r.Cookie("open_unit_id")
	if err != nil {
		fmt.Println("Error getting open_unit_id cookie:", err)
		return false, 0
	}
	openUnitIDStr := openUnitIDCookie.Value
	openUnitID, err := strconv.ParseInt(openUnitIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting open_unit_id to int:", err)
		return false, 0
	}
	return openUnitIsProposed, openUnitID
}

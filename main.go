package main

import (
	"html/template"
	"go_project2/data" // импортируем наш пакет data (путь зависит от module name)
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// Обслуживание статических файлов (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршруты
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/electrolyte/", electrolyteHandler)
	http.HandleFunc("/concentration/", concentrationHandler)

	log.Println("Сервер запущен на http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Обработчик главной страницы (список растворов)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем поисковый запрос из параметра ?search=
	search := r.URL.Query().Get("search")
	filtered := data.Electrolytes
	if search != "" {
		var tmp []data.Electrolyte
		for _, e := range data.Electrolytes {
			if strings.Contains(strings.ToLower(e.Name), strings.ToLower(search)) {
				tmp = append(tmp, e)
			}
		}
		filtered = tmp
	}

	// Данные для шаблона
	type PageData struct {
		Electrolytes []data.Electrolyte
		Search       string
		CartCount    int
		CalcID       int
	}
	pageData := PageData{
		Electrolytes: filtered,
		Search:       search,
		CartCount:    len(data.CurrentConcentration.Items),
		CalcID:       data.CurrentConcentration.ID,
	}

	// Парсим и выполняем шаблон
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик детальной страницы раствора
func electrolyteHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL, например /electrolyte/1
	idStr := strings.TrimPrefix(r.URL.Path, "/electrolyte/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var electrolyte *data.Electrolyte
	for _, e := range data.Electrolytes {
		if e.ID == id {
			electrolyte = &e
			break
		}
	}
	if electrolyte == nil {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/electrolyte.html"))
	err = tmpl.Execute(w, electrolyte)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик страницы расчёта (заявки)
func concentrationHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID, например /concentration/1
	idStr := strings.TrimPrefix(r.URL.Path, "/concentration/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id != data.CurrentConcentration.ID {
		http.NotFound(w, r)
		return
	}

	// Для каждого элемента подставляем полные данные раствора
	type ItemWithDetails struct {
		data.ConcentrationItem
		Name         string
		Concentration float64
		Ions         string
		Image        string
	}
	var items []ItemWithDetails
	for _, item := range data.CurrentConcentration.Items {
		for _, e := range data.Electrolytes {
			if e.ID == item.ElectrolyteID {
				items = append(items, ItemWithDetails{
					ConcentrationItem: item,
					Name:            e.Name,
					Concentration:   e.Concentration,
					Ions:            e.Ions,
					Image:           e.Image,
				})
				break
			}
		}
	}

	// Структура для передачи в шаблон
	type PageData struct {
		Concentration data.Concentration
		Items       []ItemWithDetails
	}
	pageData := PageData{
		Concentration: data.CurrentConcentration,
		Items:       items,
	}

	tmpl := template.Must(template.ParseFiles("templates/concentration.html"))
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
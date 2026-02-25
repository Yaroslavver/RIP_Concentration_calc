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
	http.HandleFunc("/calculation/", calculationHandler)

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
		CartCount:    len(data.CurrentCalculation.Items),
		CalcID:       data.CurrentCalculation.ID,
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
func calculationHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID, например /calculation/1
	idStr := strings.TrimPrefix(r.URL.Path, "/calculation/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id != data.CurrentCalculation.ID {
		http.NotFound(w, r)
		return
	}

	// Для каждого элемента подставляем полные данные раствора
	type ItemWithDetails struct {
		data.CalculationItem
		Name         string
		Concentration float64
		Ions         string
		Image        string
	}
	var items []ItemWithDetails
	for _, item := range data.CurrentCalculation.Items {
		for _, e := range data.Electrolytes {
			if e.ID == item.ElectrolyteID {
				items = append(items, ItemWithDetails{
					CalculationItem: item,
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
		Calculation data.Calculation
		Items       []ItemWithDetails
	}
	pageData := PageData{
		Calculation: data.CurrentCalculation,
		Items:       items,
	}

	tmpl := template.Must(template.ParseFiles("templates/calculation.html"))
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
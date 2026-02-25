package data

// Electrolyte - структура раствора (услуги)
type Electrolyte struct {
	ID           int
	Name         string
	Concentration float64
	Ions         string
	PH           float64
	Description  string
	Image        string // ссылка на изображение в Minio
	Video        string // ссылка на видео в Minio
}

// CalculationItem - элемент расчёта (выбранный раствор с объёмом и комментарием)
type CalculationItem struct {
	ElectrolyteID int
	Volume        int
}

// Calculation - текущий расчёт (заявка)
type Calculation struct {
	ID          int
	Date        string
	Description string
	Result      string
	Items       []CalculationItem
}

// Исходные данные (вместо БД)
var Electrolytes = []Electrolyte{
	{
		ID:           1,
		Name:         "Соляная кислота",
		Concentration: 0.1,
		Ions:         "H⁺, Cl⁻",
		PH:           1.0,
		Description:  "Сильная одноосновная кислота, полностью диссоциирует в воде. Используется в лабораториях для регулирования pH и в промышленности.",
		Image:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC9pbWcxLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPVJPSzZBQUg3VlBZM0dST0UwNkQyJTJGMjAyNjAyMjQlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjYwMjI0VDIzNDg1N1omWC1BbXotRXhwaXJlcz00MzE5OSZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSlNUMHMyUVVGSU4xWlFXVE5IVWs5Rk1EWkVNaUlzSW1WNGNDSTZNVGMzTWpBeE9UYzJOU3dpY0dGeVpXNTBJam9pY205dmRDSjkua2d6akdLWTF6ZlA3c080d256RF92M1lvSGdnVGxNRWdoclJ1dkN0NDI5bTh1b1I4d0k5T2VieWt6ekZXQ2VEaFF2SDB6NG83LWExZVBQLU1uc09CV1EmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT04MWY4OGVlYWFhMmQ5MWYyOTNmNjdjNDFiM2MwODMyMzI3NjA1YjQ0NzBjYWU2YTJmNjdkMmQzYzVmNDE0OGI4", // замените на реальные ссылки
		Video:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC92aWRlbzEubXA0P1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9Uk9LNkFBSDdWUFkzR1JPRTA2RDIlMkYyMDI2MDIyNCUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNjAyMjRUMjM1MDA3WiZYLUFtei1FeHBpcmVzPTQzMjAwJlgtQW16LVNlY3VyaXR5LVRva2VuPWV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpoWTJObGMzTkxaWGtpT2lKU1QwczJRVUZJTjFaUVdUTkhVazlGTURaRU1pSXNJbVY0Y0NJNk1UYzNNakF4T1RjMk5Td2ljR0Z5Wlc1MElqb2ljbTl2ZENKOS5rZ3pqR0tZMXpmUDdzTzR3bnpEX3YzWW9IZ2dUbE1FZ2hyUnV2Q3Q0MjltOHVvUjh3STlPZWJ5a3p6RldDZURoUXZIMHo0bzctYTFlUFAtTW5zT0JXUSZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmdmVyc2lvbklkPW51bGwmWC1BbXotU2lnbmF0dXJlPWRjMDk4NGFhZmY3OGE4NTg3N2E5ZjE4MzYyZDI4MjBkNzM5N2NlNjE5ZjE2NmQ1Njc0ZTg2NGVhNGJkNDZkOTM",
	},
	{
		ID:           2,
		Name:         "Гидроксид натрия",
		Concentration: 0.05,
		Ions:         "Na⁺, OH⁻",
		PH:           12.7,
		Description:  "Сильное основание, щёлочь.",
		Image:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC9pbWcyLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPVJPSzZBQUg3VlBZM0dST0UwNkQyJTJGMjAyNjAyMjQlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjYwMjI0VDIzNTEwNFomWC1BbXotRXhwaXJlcz00MzIwMCZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSlNUMHMyUVVGSU4xWlFXVE5IVWs5Rk1EWkVNaUlzSW1WNGNDSTZNVGMzTWpBeE9UYzJOU3dpY0dGeVpXNTBJam9pY205dmRDSjkua2d6akdLWTF6ZlA3c080d256RF92M1lvSGdnVGxNRWdoclJ1dkN0NDI5bTh1b1I4d0k5T2VieWt6ekZXQ2VEaFF2SDB6NG83LWExZVBQLU1uc09CV1EmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT1jYWZhM2ZhMTdhYjAwOTYyODJkYzYyYjFkMWVlYzAzNzYyMjUwMmNhZjQxNGQzYjBhZWM0NTU2NDllYmRkYTdm",
		Video:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC92aWRlbzIubXA0P1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9Uk9LNkFBSDdWUFkzR1JPRTA2RDIlMkYyMDI2MDIyNCUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNjAyMjRUMjM1MTI0WiZYLUFtei1FeHBpcmVzPTQzMTk5JlgtQW16LVNlY3VyaXR5LVRva2VuPWV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpoWTJObGMzTkxaWGtpT2lKU1QwczJRVUZJTjFaUVdUTkhVazlGTURaRU1pSXNJbVY0Y0NJNk1UYzNNakF4T1RjMk5Td2ljR0Z5Wlc1MElqb2ljbTl2ZENKOS5rZ3pqR0tZMXpmUDdzTzR3bnpEX3YzWW9IZ2dUbE1FZ2hyUnV2Q3Q0MjltOHVvUjh3STlPZWJ5a3p6RldDZURoUXZIMHo0bzctYTFlUFAtTW5zT0JXUSZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmdmVyc2lvbklkPW51bGwmWC1BbXotU2lnbmF0dXJlPTliODhhZjhmMjM2NGVlZGZiN2MyODk4ZTdjODY4MDUwNTYxMmExOWZhMTUzYThmNWNlYjc0MDZjZmYwMzg3ZjA",
	},
	{
		ID:           3,
		Name:         "Хлорид натрия",
		Concentration: 0.2,
		Ions:         "H⁺, Cl⁻",
		PH:           12.7,
		Description:  "Сильное основание, щёлочь.",
		Image:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC9pbWczLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPVJPSzZBQUg3VlBZM0dST0UwNkQyJTJGMjAyNjAyMjUlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjYwMjI1VDAxNDE1NFomWC1BbXotRXhwaXJlcz00MzIwMCZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSlNUMHMyUVVGSU4xWlFXVE5IVWs5Rk1EWkVNaUlzSW1WNGNDSTZNVGMzTWpBeE9UYzJOU3dpY0dGeVpXNTBJam9pY205dmRDSjkua2d6akdLWTF6ZlA3c080d256RF92M1lvSGdnVGxNRWdoclJ1dkN0NDI5bTh1b1I4d0k5T2VieWt6ekZXQ2VEaFF2SDB6NG83LWExZVBQLU1uc09CV1EmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT0xMTdlOTg0Y2E2ZmYyMmI3M2UwYWZlYjA1OThlNTBlYTg5NjFjZTViMWU5YzhhYzM0YTQzNTFhOGQ3YjZhNjE0",
		Video:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC92aWRlbzMubXA0P1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9Uk9LNkFBSDdWUFkzR1JPRTA2RDIlMkYyMDI2MDIyNSUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNjAyMjVUMDE0MjIyWiZYLUFtei1FeHBpcmVzPTQzMTk5JlgtQW16LVNlY3VyaXR5LVRva2VuPWV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpoWTJObGMzTkxaWGtpT2lKU1QwczJRVUZJTjFaUVdUTkhVazlGTURaRU1pSXNJbVY0Y0NJNk1UYzNNakF4T1RjMk5Td2ljR0Z5Wlc1MElqb2ljbTl2ZENKOS5rZ3pqR0tZMXpmUDdzTzR3bnpEX3YzWW9IZ2dUbE1FZ2hyUnV2Q3Q0MjltOHVvUjh3STlPZWJ5a3p6RldDZURoUXZIMHo0bzctYTFlUFAtTW5zT0JXUSZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmdmVyc2lvbklkPW51bGwmWC1BbXotU2lnbmF0dXJlPTkyODAxMzdjZDc5YWJkMTI0ODZkOWJlYzZkNjQ2ODU2NmUxM2ZiZmE3ZDU4YzhlZDc2NmQwYWIzMzlkYTFkMzY",
	},
	{
		ID:           4,
		Name:         "Уксусная кислота",
		Concentration: 0.1,
		Ions:         "H⁺, Cl⁻",
		PH:           12.7,
		Description:  "Сильное основание, щёлочь.",
		Image:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC9pbWc0LnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPVJPSzZBQUg3VlBZM0dST0UwNkQyJTJGMjAyNjAyMjUlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjYwMjI1VDAxNDIxMVomWC1BbXotRXhwaXJlcz00MzIwMCZYLUFtei1TZWN1cml0eS1Ub2tlbj1leUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaFkyTmxjM05MWlhraU9pSlNUMHMyUVVGSU4xWlFXVE5IVWs5Rk1EWkVNaUlzSW1WNGNDSTZNVGMzTWpBeE9UYzJOU3dpY0dGeVpXNTBJam9pY205dmRDSjkua2d6akdLWTF6ZlA3c080d256RF92M1lvSGdnVGxNRWdoclJ1dkN0NDI5bTh1b1I4d0k5T2VieWt6ekZXQ2VEaFF2SDB6NG83LWExZVBQLU1uc09CV1EmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JnZlcnNpb25JZD1udWxsJlgtQW16LVNpZ25hdHVyZT1hYmY4OWY3N2RiZmZlNmJjYWExYzE2MjFkZWYwNGE1MmVlNDViZTViMWU2YTk0MTI5YmUxZjRlMjJkYmNmNDRm",
		Video:        "http://localhost:9001/api/v1/download-shared-object/aHR0cDovLzEyNy4wLjAuMTo5MDAwL2J1Y2tldC92aWRlbzQubXA0P1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9Uk9LNkFBSDdWUFkzR1JPRTA2RDIlMkYyMDI2MDIyNSUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNjAyMjVUMDE0MjM2WiZYLUFtei1FeHBpcmVzPTQzMjAwJlgtQW16LVNlY3VyaXR5LVRva2VuPWV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpoWTJObGMzTkxaWGtpT2lKU1QwczJRVUZJTjFaUVdUTkhVazlGTURaRU1pSXNJbVY0Y0NJNk1UYzNNakF4T1RjMk5Td2ljR0Z5Wlc1MElqb2ljbTl2ZENKOS5rZ3pqR0tZMXpmUDdzTzR3bnpEX3YzWW9IZ2dUbE1FZ2hyUnV2Q3Q0MjltOHVvUjh3STlPZWJ5a3p6RldDZURoUXZIMHo0bzctYTFlUFAtTW5zT0JXUSZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmdmVyc2lvbklkPW51bGwmWC1BbXotU2lnbmF0dXJlPTRkZDgyOWI4OGE0NmU3ZGFmOTFkMzIwMGI5ODY0NTI2MmY1ZDVkMzQzZGQ2MjI4NDQ4NjVjZGY1ZmY4YTQ4OTk",
	},
}

// CurrentCalculation - пример текущего расчёта
var CurrentCalculation = Calculation{
	ID:          1,
	Date:        "2026-02-25",
	Description: "Смешивание соляной кислоты и гидроксида натрия",
	Result:      "[H⁺] = 0.045 моль/л, pH = 1.35",
	Items: []CalculationItem{
		{ElectrolyteID: 1, Volume: 50},
		{ElectrolyteID: 2, Volume: 30},
	},
}
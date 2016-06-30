package Models

type Response struct {
	Success           bool     `json:"success"`
	WebinarsProcessed []string `json:"webinars_processed"`
}

package dto

type CustomerDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Cell    string `json:"cell"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	ZipCode string `json:"zipcode"`
	City    string `json:"city"`
}

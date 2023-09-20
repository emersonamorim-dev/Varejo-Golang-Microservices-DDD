package dto

type IntegrationDTO struct {
	ID       string            `json:"id" bson:"_id"`
	Name     string            `json:"name" bson:"name"`
	Endpoint string            `json:"endpoint" bson:"endpoint"`
	APIKey   string            `json:"api_key" bson:"api_key"`
	Data     string            `json:"data" bson:"data"`
	Other    map[string]string `json:"other" bson:"other"`
}

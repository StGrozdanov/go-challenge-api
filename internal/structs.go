package internal

type User struct {
	Id        int    `db:"user_id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
}

type Location struct {
	City    string  `db:"home_city" json:"city"`
	State   string  `db:"home_state" json:"state"`
	Zip     string  `db:"home_zip" json:"zip"`
	Country string  `db:"home_country" json:"country"`
	Lat     float64 `db:"lat" json:"lat"`
	Lng     float64 `db:"lng" json:"lng"`
}

type Price struct {
	Day int `db:"price_per_day" json:"day"`
}

type Rental struct {
	IdRental        int     `db:"id" json:"id"`
	Name            string  `db:"name" json:"name"`
	Description     string  `db:"description" json:"description"`
	Type            string  `db:"type" json:"type"`
	Make            string  `db:"vehicle_make" json:"make"`
	Model           string  `db:"vehicle_model" json:"model"`
	Year            int     `db:"vehicle_year" json:"year"`
	Length          float64 `db:"vehicle_length" json:"length"`
	Sleeps          int     `db:"sleeps" json:"sleeps"`
	PrimaryImageURL string  `db:"primary_image_url" json:"primary_image_url"`
	Price           `json:"price"`
	Location        `json:"location"`
	User            `json:"user"`
}

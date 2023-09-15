package internal

var selectSingleRentalQuery = `
				SELECT rentals.id,
					   rentals.name,
					   rentals.description,
					   rentals.type,
					   rentals.vehicle_make,
					   rentals.vehicle_model,
					   rentals.vehicle_year,
					   rentals.vehicle_length,
					   rentals.sleeps,
					   rentals.primary_image_url,
					   rentals.price_per_day,
					   rentals.home_city,
					   rentals.home_state,
					   rentals.home_zip,
					   rentals.home_country,
					   rentals.lat,
					   rentals.lng,
					   users.id AS user_id,
					   users.first_name,
					   users.last_name
				FROM rentals
						LEFT JOIN users ON users.id = rentals.user_id
				WHERE rentals.id = :id;`

var selectAllRentalsQuery = `
				SELECT rentals.id,
					   rentals.name,
					   rentals.description,
					   rentals.type,
					   rentals.vehicle_make,
					   rentals.vehicle_model,
					   rentals.vehicle_year,
					   rentals.vehicle_length,
					   rentals.sleeps,
					   rentals.primary_image_url,
					   rentals.price_per_day,
					   rentals.home_city,
					   rentals.home_state,
					   rentals.home_zip,
					   rentals.home_country,
					   rentals.lat,
					   rentals.lng,
					   users.id AS user_id,
					   users.first_name,
					   users.last_name
				FROM rentals
						LEFT JOIN users ON users.id = rentals.user_id`

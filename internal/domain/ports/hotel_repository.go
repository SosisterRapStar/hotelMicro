package ports

import "github.com/SosisterRapStar/hotels/internal/domain/hotel"

type HotelRepository interface {
	hotel.Repository
}

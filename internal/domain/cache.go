package domain

import "time"

const (
	HalfCacheExpiration = 12 * time.Hour

	ProductKeyCache = "product"
	CartKeyCache    = "cart"
)

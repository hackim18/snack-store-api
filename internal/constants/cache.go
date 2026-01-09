package constants

import "time"

const (
	ProductCacheKeyPrefix = "products:date:"
	ReportCacheKeyPrefix  = "report:transactions:"
)

const (
	ProductCacheTTL = 5 * time.Minute
	ReportCacheTTL  = 2 * time.Minute
)

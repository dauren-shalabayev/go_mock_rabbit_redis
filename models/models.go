package models

// CacheValue содержит информацию о номере.
type CacheValue struct {
	Imsi     string `json:"imsi"`
	LacCell  string `json:"lac_cell"`
	SectorID int    `json:"sector_id"`
}

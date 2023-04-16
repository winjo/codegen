package dao

import "database/sql"

type (
	SampleDAO struct {
		*baseSampleDAO
	}
)

func NewSampleDAO(db *sql.DB) *SampleDAO {
	return &SampleDAO{
		baseSampleDAO: newBaseSampleDAO(db),
	}
}

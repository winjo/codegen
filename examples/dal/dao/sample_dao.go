package dao

type (
	SampleDAO struct {
		*baseSampleDAO
	}
)

func NewSampleDAO(q Queryer) *SampleDAO {
	return &SampleDAO{
		baseSampleDAO: newBaseSampleDAO(q),
	}
}

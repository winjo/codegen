package dao

import (
	"context"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
	"time"
)

var sampleTestDataList = []*SampleData{
	{
		RInt:    0,
		NInt:    null.Int{},
		RFloat:  0,
		NFloat:  null.Float{},
		RString: "r0",
		NString: null.String{},
		RTime:   time.Unix(time.Now().Unix()+int64(0), 0).UTC(),
		NTime:   null.Time{},
		Union1:  "u0",
		Union2:  "u0",
		Union3:  null.String{},
	},
	{
		RInt:    1,
		NInt:    null.Int{},
		RFloat:  1,
		NFloat:  null.Float{},
		RString: "r1",
		NString: null.String{},
		RTime:   time.Unix(time.Now().Unix()+int64(1), 0).UTC(),
		NTime:   null.Time{},
		Union1:  "u1",
		Union2:  "u1",
		Union3:  null.String{},
	},
}
var sampleTestNoData = &SampleData{
	RInt:    404,
	NInt:    null.IntFrom(404),
	RFloat:  404,
	NFloat:  null.FloatFrom(404),
	RString: "r404",
	NString: null.StringFrom("n404"),
	RTime:   time.Unix(int64(404), 0),
	NTime:   null.TimeFrom(time.Unix(int64(404), 0)),
	Union1:  "u404",
	Union2:  "u404",
	Union3:  null.StringFrom("u404"),
}

// 设置为零值，便于比较
func emptyAutoValue(data *SampleData) *SampleData {
	clone := *data
	clone.Id = 0
	clone.GmtCreate = time.Time{}
	clone.GmtModified = time.Time{}
	return &clone
}

func TestSampleCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	// 新建
	dao := NewSampleDAO(tdb)
	r0, err := dao.Insert(ctx, sampleTestDataList[0])
	assert.NoError(err)
	id0, err := r0.LastInsertId()
	assert.NoError(err)

	r1, err := dao.Insert(ctx, sampleTestDataList[1])
	assert.NoError(err)
	id1, err := r1.LastInsertId()
	assert.NoError(err)

	pageInData, err := dao.PageInById(ctx, []int64{id0, id1}, 1, 2)
	assert.NoError(err)

	assert.Equal(pageInData.Count, 2)
	assert.Equal(len(pageInData.Data), 2)
	assert.Equal(*sampleTestDataList[0], *emptyAutoValue(pageInData.Data[0]))
	assert.Equal(*sampleTestDataList[1], *emptyAutoValue(pageInData.Data[1]))

	pageInDataOrderByGmtCreate, err := dao.PageInById(ctx, []int64{id0, id1}, 1, 2, QueryOption{OrderBy: mo.Some("gmt_create")})
	assert.NoError(err)
	assert.Equal(*sampleTestDataList[0], *emptyAutoValue(pageInDataOrderByGmtCreate.Data[0]))
	assert.Equal(*sampleTestDataList[1], *emptyAutoValue(pageInDataOrderByGmtCreate.Data[1]))

	pageInDataDesc, err := dao.PageInById(ctx, []int64{id0, id1}, 1, 2, QueryOption{Desc: mo.Some(true)})
	assert.NoError(err)
	assert.Equal(*sampleTestDataList[0], *emptyAutoValue(pageInDataDesc.Data[1]))
	assert.Equal(*sampleTestDataList[1], *emptyAutoValue(pageInDataDesc.Data[0]))

	// 查询
	findData, err := dao.Find(ctx)
	assert.NoError(err)
	assert.Equal(len(findData), 2)

	pageData, err := dao.Page(ctx, 1, 100)
	assert.NoError(err)
	assert.Equal(len(pageData.Data), 2)

	count, err := dao.Count(ctx)
	assert.NoError(err)
	assert.Equal(count, 2)

	data := pageInData.Data[0]
	noData := sampleTestNoData

	findByNTime, err := dao.FindByNTime(ctx, data.NTime)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByNTime), 1)

	findByNTime404, err := dao.FindByNTime(ctx, noData.NTime)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByNTime404), 0)

	pageByNTime, err := dao.PageByNTime(ctx, data.NTime, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByNTime.Count, 1)

	pageByNTime404, err := dao.PageByNTime(ctx, noData.NTime, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByNTime404.Count, 0)

	findInByNTime, err := dao.FindInByNTime(ctx, []null.Time{data.NTime})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByNTime), 1)

	findInByNTime404, err := dao.FindInByNTime(ctx, []null.Time{noData.NTime})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByNTime404), 0)

	pageInByNTime, err := dao.PageInByNTime(ctx, []null.Time{data.NTime}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByNTime.Count, 1)

	pageInByNTime404, err := dao.PageInByNTime(ctx, []null.Time{noData.NTime}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByNTime404.Count, 0)

	findByRTime, err := dao.FindByRTime(ctx, data.RTime)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByRTime), 1)

	findByRTime404, err := dao.FindByRTime(ctx, noData.RTime)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByRTime404), 0)

	pageByRTime, err := dao.PageByRTime(ctx, data.RTime, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByRTime.Count, 1)

	pageByRTime404, err := dao.PageByRTime(ctx, noData.RTime, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByRTime404.Count, 0)

	findInByRTime, err := dao.FindInByRTime(ctx, []time.Time{data.RTime})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByRTime), 1)

	findInByRTime404, err := dao.FindInByRTime(ctx, []time.Time{noData.RTime})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByRTime404), 0)

	pageInByRTime, err := dao.PageInByRTime(ctx, []time.Time{data.RTime}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByRTime.Count, 1)

	pageInByRTime404, err := dao.PageInByRTime(ctx, []time.Time{noData.RTime}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByRTime404.Count, 0)

	findByUnion1AndUnion3, err := dao.FindByUnion1AndUnion3(ctx, data.Union1, data.Union3)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByUnion1AndUnion3), 1)

	findByUnion1AndUnion3404, err := dao.FindByUnion1AndUnion3(ctx, noData.Union1, noData.Union3)
	assert.NoError(err)
	assert.GreaterOrEqual(len(findByUnion1AndUnion3404), 0)

	pageByUnion1AndUnion3, err := dao.PageByUnion1AndUnion3(ctx, data.Union1, data.Union3, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByUnion1AndUnion3.Count, 1)

	pageByUnion1AndUnion3404, err := dao.PageByUnion1AndUnion3(ctx, noData.Union1, noData.Union3, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageByUnion1AndUnion3404.Count, 0)

	getById, err := dao.GetById(ctx, data.Id)
	assert.NoError(err)
	assert.NotNil(getById)

	getById404, err := dao.GetById(ctx, noData.Id)
	assert.NoError(err)
	assert.Nil(getById404)

	findInById, err := dao.FindInById(ctx, []int64{data.Id})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInById), 1)

	findInById404, err := dao.FindInById(ctx, []int64{noData.Id})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInById404), 0)

	pageInById, err := dao.PageInById(ctx, []int64{data.Id}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInById.Count, 1)

	pageInById404, err := dao.PageInById(ctx, []int64{noData.Id}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInById404.Count, 0)

	getByNInt, err := dao.GetByNInt(ctx, data.NInt)
	assert.NoError(err)
	assert.NotNil(getByNInt)

	getByNInt404, err := dao.GetByNInt(ctx, noData.NInt)
	assert.NoError(err)
	assert.Nil(getByNInt404)

	findInByNInt, err := dao.FindInByNInt(ctx, []null.Int{data.NInt})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByNInt), 1)

	findInByNInt404, err := dao.FindInByNInt(ctx, []null.Int{noData.NInt})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByNInt404), 0)

	pageInByNInt, err := dao.PageInByNInt(ctx, []null.Int{data.NInt}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByNInt.Count, 1)

	pageInByNInt404, err := dao.PageInByNInt(ctx, []null.Int{noData.NInt}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByNInt404.Count, 0)

	getByRInt, err := dao.GetByRInt(ctx, data.RInt)
	assert.NoError(err)
	assert.NotNil(getByRInt)

	getByRInt404, err := dao.GetByRInt(ctx, noData.RInt)
	assert.NoError(err)
	assert.Nil(getByRInt404)

	findInByRInt, err := dao.FindInByRInt(ctx, []int64{data.RInt})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByRInt), 1)

	findInByRInt404, err := dao.FindInByRInt(ctx, []int64{noData.RInt})
	assert.NoError(err)
	assert.GreaterOrEqual(len(findInByRInt404), 0)

	pageInByRInt, err := dao.PageInByRInt(ctx, []int64{data.RInt}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByRInt.Count, 1)

	pageInByRInt404, err := dao.PageInByRInt(ctx, []int64{noData.RInt}, 1, 100)
	assert.NoError(err)
	assert.GreaterOrEqual(pageInByRInt404.Count, 0)

	getByUnion1AndUnion2, err := dao.GetByUnion1AndUnion2(ctx, data.Union1, data.Union2)
	assert.NoError(err)
	assert.NotNil(getByUnion1AndUnion2)

	getByUnion1AndUnion2404, err := dao.GetByUnion1AndUnion2(ctx, noData.Union1, noData.Union2)
	assert.NoError(err)
	assert.Nil(getByUnion1AndUnion2404)

	// 修改
	err = dao.UpdatePartial(ctx, id0, &SamplePartialData{})
	assert.NoError(err)
	afterUpdatePartial, err := dao.GetById(ctx, id0)
	assert.NoError(err)
	assert.NotNil(afterUpdatePartial)
	assert.Equal(*sampleTestDataList[0], *emptyAutoValue(afterUpdatePartial))

	err = dao.Delete(ctx, id1)
	assert.NoError(err)
	err = dao.Update(ctx, id0, sampleTestDataList[1])
	assert.NoError(err)
	afterUpdate, err := dao.GetById(ctx, id0)
	assert.NoError(err)
	afterUpdate.Id = 0
	afterUpdate.GmtCreate = time.Time{}
	afterUpdate.GmtModified = time.Time{}
	assert.Equal(*sampleTestDataList[1], *emptyAutoValue(afterUpdate))

	// 删除
	err = dao.Delete(ctx, id0)
	assert.NoError(err)
	afterDelete, err := dao.GetById(ctx, id1)
	assert.Nil(afterDelete)
}

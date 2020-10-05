package db

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestTableDB struct {
	Id     uint32        `gorm:"type:integer;primary_key"`
	Name   string        `gorm:"type:varchar(32)"`
	ISlc   pq.Int64Array `gorm:"type:smallint[]"`
	Amount string        `gorm:"type:money"`
}

func TestDB(t *testing.T) {
	c := DefaultConf()
	err := Init(&c)
	require.Empty(t, err)

	defer Close()

	if DB().HasTable(&TestTableDB{}) {
		err = DB().DropTable(&TestTableDB{}).Error
		require.Empty(t, err)
	}

	err = DB().CreateTable(&TestTableDB{}).Error
	require.Empty(t, err)

	dbData := &TestTableDB{
		Id:   666,
		Name: "hello world",
		ISlc: pq.Int64Array{2333, 3332},
	}

	err = DB().Create(dbData).Error
	require.Empty(t, err)
}

package db

import (
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"xorm.io/core"
)

type TestTableDB struct {
	Id     uint32        `xorm:"integer autoincr pk"`
	Name   string        `xorm:"varchar(32)"`
	ISlc   pq.Int64Array `xorm:"i_slc" pq:",array"`
	Amount string        `xorm:"money"`
}

func TestDB(t *testing.T) {
	c := DefaultConf()
	err := Init(&c)
	require.Empty(t, err)

	defer Close()

	var has bool
	has, err = DB().Exist(&TestTableDB{})
	require.Empty(t, err)

	if has {
		_, err = DB().Exec("truncate table test_table_d_b")
		require.Empty(t, err)
	}

	// err = DB().CreateTables(&TestTableDB{})
	// require.Empty(t, err)

	dbData := &TestTableDB{
		// Id:     1,
		Name:   "hello world",
		ISlc:   pq.Int64Array{2333, 3332},
		Amount: "50000",
	}

	// var num int64
	_, err = DB().Exec(`INSERT INTO test_table_d_b(name,i_slc,amount) VALUES ($1, $2,$3)`, dbData.Name, dbData.ISlc, dbData.Amount)
	// num, err = DB().Insert(dbData)
	require.Empty(t, err)
	// require.Equal(t, int64(1), num)

	// var res []map[string][]byte
	var newData TestTableDB
	// jres, err = DB().
	// j	Query("select id,i_slc from test_table_d_b")
	// j// _, err = DB().Table("test_table_d_b").Get(&newData)
	// jnewData.ISlc.Scan(res[0]["i_slc"])
	// jfmt.Println("=====", newData.ISlc, res)

	var rows *core.Rows
	rows, err = CDB().Query("select id,i_slc from test_table_d_b")
	require.Empty(t, err)

	for rows.Next() {
		err = rows.Scan(&newData.Id, &newData.ISlc)
		require.Empty(t, err)
	}

	fmt.Println("=====", newData.ISlc, newData.Id)

	// require.Equal(t, uint32(1), dbData.Id)

	// dbData.Id = 0
	// dbData.Name = "hello slc"
	// dbSlcData := []*TestTableDB{dbData}
	// num, err = DB().Insert(dbSlcData)
	// require.Empty(t, err)
	// require.Empty(t, int64(1), num)
	// require.Equal(t, uint32(2), dbSlcData[0].Id)
}

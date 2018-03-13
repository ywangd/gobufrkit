package table

import (
    "testing"
    assert2 "github.com/seanpont/assert"
)

func TestLoadTableB(t *testing.T) {

    assert := assert2.Assert(t)

    b, err := LoadTableB("../_definitions/tables/0/0/0/25/TableB.csv")
    assert.Nil(err)

    descriptor, err := b.Lookup(ID(1001))
    assert.Nil(err)
    assert.Equal(descriptor.Id(), ID(1001))
    assert.Equal(descriptor.Entry().Name(), "WMO BLOCK NUMBER")
}

func TestLoadTableD(t *testing.T) {
    assert := assert2.Assert(t)

    d, err := LoadTableD("../_definitions/tables/0/0/0/25/TableD.csv")
    assert.Nil(err)

    descriptor, err := d.Lookup(ID(301001))
    assert.Nil(err)
    assert.Equal(descriptor.Id(), ID(301001))
    assert.Equal(descriptor.Entry().Name(), "(WMO block and station numbers)")

    assert.Equal(descriptor.Entry().(*Dentry).Members, []ID{ID(1001), ID(1002)})

}

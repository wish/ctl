package client

import (
	"reflect"
	"sort"
	"strings"
)

type alphabetical struct {
	reflect.Value
}

func (a alphabetical) Len() int { return a.Value.Len() }
func (a alphabetical) Swap(i, j int) {
	temp := a.Index(i).Interface()
	a.Index(i).Set(a.Index(j))
	a.Index(j).Set(reflect.ValueOf(temp))
}
func (a alphabetical) Less(i, j int) bool {
	x := a.Index(i)
	y := a.Index(j)

	c1 := x.FieldByName("Context").String()
	c2 := y.FieldByName("Context").String()
	if strings.Compare(c1, c2) < 0 {
		return true
	} else if strings.Compare(c1, c2) > 0 {
		return false
	}

	ns1 := x.FieldByIndex([]int{1, 1, 2}).String()
	ns2 := y.FieldByIndex([]int{1, 1, 2}).String()
	if strings.Compare(ns1, ns2) < 0 {
		return true
	}
	return false
}

func sortObjs(objs interface{}) {
	v := reflect.ValueOf(objs)
	if v.Len() == 0 {
		return
	}
	sort.Stable(alphabetical{v})
}

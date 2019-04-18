package dbdiff

import (
	"reflect"
	"sort"
)

type KeySlice struct {
	keyCompareAction KeyCompareAction
	keyComparator    func(left, right interface{}) int
}

type KeyCompareAction interface {
	ActionBothExists(itemLeft, itemRight interface{})
	ActionLeftExists(itemLeft interface{})
	ActionRightExists(itemRight interface{})
}

func (comp *KeySlice) Compare(itemsLeft, itemsRight interface{}) {
	var (
		sortableSliceLeft = sortableSlice{
			items:      itemsLeft,
			comparator: comp.keyComparator,
		}
		sortableSliceRight = sortableSlice{
			items:      itemsRight,
			comparator: comp.keyComparator,
		}
		idxLeft  = 0
		idxRight = 0

		sizeLeft  = reflect.ValueOf(itemsLeft).Elem().Len()
		sizeRight = reflect.ValueOf(itemsRight).Elem().Len()
	)

	sort.Sort(sortableSliceLeft)
	sort.Sort(sortableSliceRight)

	for ; idxLeft < sizeLeft && idxRight < sizeRight; {
		var (
			itemLeft  = reflect.ValueOf(itemsLeft).Elem().Index(idxLeft).Interface()
			itemRight = reflect.ValueOf(itemsRight).Elem().Index(idxRight).Interface()
			res       = comp.keyComparator(itemLeft, itemRight)
		)
		if res == 0 {
			comp.keyCompareAction.ActionBothExists(itemLeft, itemRight)
			idxLeft++
			idxRight++
		} else if res <= 0 {
			comp.keyCompareAction.ActionLeftExists(itemLeft)
			idxLeft++
		} else {
			comp.keyCompareAction.ActionRightExists(itemRight)
			idxRight++
		}
	}

	for ; idxLeft < sizeLeft; {
		itemLeft := reflect.ValueOf(itemsLeft).Elem().Index(idxLeft).Interface()
		comp.keyCompareAction.ActionLeftExists(itemLeft)
		idxLeft++
	}

	for ; idxRight < sizeRight; {
		itemRight := reflect.ValueOf(itemsRight).Elem().Index(idxRight).Interface()
		comp.keyCompareAction.ActionRightExists(itemRight)
		idxRight++
	}
}

type sortableSlice struct {
	items      interface{}
	comparator func(left, right interface{}) int
}

func (this sortableSlice) Len() int {
	return reflect.ValueOf(this.items).Elem().Len()
}

func (this sortableSlice) Less(i, j int) bool {
	v := reflect.ValueOf(this.items).Elem()
	return this.comparator(v.Index(i).Interface(), v.Index(j).Interface()) < 0
}

func (this sortableSlice) Swap(i, j int) {
	v := reflect.ValueOf(this.items).Elem()
	x, y := v.Index(i).Interface(), v.Index(j).Interface()
	v.Index(i).Set(reflect.ValueOf(y))
	v.Index(j).Set(reflect.ValueOf(x))
}

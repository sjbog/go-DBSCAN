package dbscan

import (
	"strconv"
	"testing"
)

func Test_NamedPoint(t *testing.T) {
	t.Parallel()
	var (
		size = 10
		data = make([]ClusterablePoint, 0, size)
		prev float64
	)

	for i := 0; i < size; i += 1 {
		data = append(data, &NamedPoint{
			Name:  strconv.Itoa(i),
			Point: []float64{float64(i), float64(size - i - 1)},
		})
	}

	// t.Logf("%v", data)
	ClusterablePointSlice{Data: data, SortDimension: 1}.Sort()
	prev = data[0].GetPoint()[1]
	// t.Logf("%v", data)

	for i := 1; i < size; i += 1 {
		if prev > data[i].GetPoint()[1] {
			t.FailNow()
		}
		prev = data[i].GetPoint()[1]
	}

	ClusterablePointSlice{Data: data, SortDimension: 0}.Sort()
	prev = data[0].GetPoint()[0]
	// t.Logf("%v", data)

	for i := 1; i < size; i += 1 {
		if prev > data[i].GetPoint()[0] {
			t.FailNow()
		}
		prev = data[i].GetPoint()[0]
	}
}

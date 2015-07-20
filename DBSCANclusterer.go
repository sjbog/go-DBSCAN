package dbscan

import (
	"container/list"
	"sync"
)

type Clusterer interface {
	Cluster([]ClusterablePoint) [][]ClusterablePoint
}

type DBSCANClusterer struct {
	eps, eps2                                 float64
	MinPts, numDimensions, SortDimensionIndex int
	AutoSelectDimension                       bool
}

func NewDBSCANClusterer(eps float64, minPts int) *DBSCANClusterer {
	return &DBSCANClusterer{
		eps:    eps,
		eps2:   eps * eps,
		MinPts: minPts,

		AutoSelectDimension: true,
	}
}

func (this *DBSCANClusterer) GetEps() float64 {
	return this.eps
}
func (this *DBSCANClusterer) SetEps(eps float64) {
	this.eps = eps
	this.eps2 = eps * eps
}

/**
step 1: sort data by a dimension
step 2: slide through sorted data (in parallel), and compute all points in range of eps (everything above eps is definitely isn't directly reachable)
step 3: build neighborhood map & proceed DFS
**/
func (this *DBSCANClusterer) Cluster(data []ClusterablePoint) [][]ClusterablePoint {
	if len(data) == 0 {
		return [][]ClusterablePoint{}
	}
	var (
		dataSize   = len(data)
		clusters   = make([][]ClusterablePoint, 0, 64)
		visitedMap = make([]bool, dataSize)
		cluster    = make([]ClusterablePoint, 0, 64)

		neighborhoodMap []*ConcurrentQueue_InsertOnly
	)

	this.numDimensions = len(data[0].GetPoint())

	if this.AutoSelectDimension {
		this.SortDimensionIndex = this.PredictDimensionByMaxVariance(data)
	}

	ClusterablePointSlice{
		Data:          data,
		SortDimension: this.SortDimensionIndex,
	}.Sort()

	neighborhoodMap = this.BuildNeighborhoodMap(data)

	// Early exit - 1 huge cluster
	if neighborhoodMap[0].Size == uint64(dataSize) {
		cluster = make([]ClusterablePoint, 0, dataSize)

		for _, v := range neighborhoodMap[0].Slice() {
			cluster = append(cluster, data[v])
		}

		clusters = append(clusters, cluster)
		return clusters
	}

	var (
		queue = list.New()
		elem  *list.Element
	)

	for pointIndex, tmpIndex := 0, uint(0); pointIndex < dataSize; pointIndex += 1 {
		if visitedMap[pointIndex] {
			continue
		}
		// Expand cluster
		queue.PushBack(uint(pointIndex))

		// DFS
		for queue.Len() > 0 {
			// Pop last elem
			elem = queue.Back()
			queue.Remove(elem)

			tmpIndex = elem.Value.(uint)
			if visitedMap[tmpIndex] {
				continue
			}

			cluster = append(cluster, data[tmpIndex])
			visitedMap[tmpIndex] = true

			for _, v := range neighborhoodMap[tmpIndex].Slice() {
				queue.PushBack(v)
			}
		}

		if len(cluster) >= this.MinPts {
			clusters = append(clusters, cluster)
		}

		cluster = make([]ClusterablePoint, 0, 64)
	}
	return clusters
}

func (this *DBSCANClusterer) CalcDistance(aPoint, bPoint []float64) float64 {
	var sum = 0.0
	for i, size := 0, this.numDimensions; i < size; i += 1 {
		sum += (aPoint[i] - bPoint[i]) * (aPoint[i] - bPoint[i])
	}
	return sum
}

func (this *DBSCANClusterer) BuildNeighborhoodMap(data []ClusterablePoint) []*ConcurrentQueue_InsertOnly {
	var (
		dataSize  = len(data)
		result    = make([]*ConcurrentQueue_InsertOnly, dataSize)
		waitGroup = new(sync.WaitGroup)

		fn = func(start int) {
			defer waitGroup.Done()
			var (
				x, head ClusterablePoint = nil, data[start]

				headV    []float64 = head.GetPoint()
				headDimV float64   = headV[this.SortDimensionIndex] + this.eps
			)
			if result[start] == nil {
				result[start] = NewConcurrentQueue_InsertOnly()
			}
			result[start].Add(uint(start))

			for i := start + 1; i < dataSize && data[i].GetPoint()[this.SortDimensionIndex] <= headDimV; i += 1 {
				x = data[i]

				if this.CalcDistance(headV, x.GetPoint()) <= this.eps2 {
					result[start].Add(uint(i))
					if result[i] == nil {
						result[i] = NewConcurrentQueue_InsertOnly()
					}
					result[i].Add(uint(start))
				}
			}
		}
	)
	waitGroup.Add(dataSize)

	// Early exit - 1 huge cluster
	fn(0)

	if result[0].Size == uint64(dataSize) {
		return result
	}

	for i := 1; i < dataSize; i += 1 {
		go fn(i)
	}
	waitGroup.Wait()

	return result
}

/**
 * Calculate variance for each dimension (in parallel), returns dimension index with max variance
 */
func (this *DBSCANClusterer) PredictDimensionByMaxVariance(data []ClusterablePoint) int {
	var (
		waitGroup = new(sync.WaitGroup)
		result    = make([]float64, this.numDimensions)
	)
	waitGroup.Add(int(this.numDimensions))

	for i, size := 0, this.numDimensions; i < size; i += 1 {
		go func(dim int) {
			result[dim] = Variance(data, dim)
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	var (
		maxV = 0.0
		maxI = 0
	)
	for i, v := range result {
		if maxV <= v {
			maxV = v
			maxI = i
		}
	}
	return maxI
}

func Variance(data []ClusterablePoint, dimension int) float64 {
	var (
		size     = len(data)
		avg      = 0.0
		sum      = 0.0
		delta, v float64
	)
	if size < 2 {
		return 0.0
	}
	for i, point := range data {
		v = point.GetPoint()[dimension]
		delta = v - avg
		avg += delta / float64(i+1)
		sum += delta * (v - avg)
	}
	return sum / float64(size-1)
}

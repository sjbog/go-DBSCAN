package dbscan

import "testing"

func Test_DBSCAN(t *testing.T) {
	t.Parallel()
	var (
		data = []ClusterablePoint{
			&NamedPoint{"0", []float64{2, 4}},
			&NamedPoint{"1", []float64{7, 3}},
			&NamedPoint{"2", []float64{3, 5}},
			&NamedPoint{"3", []float64{5, 3}},
			&NamedPoint{"4", []float64{7, 4}},
			// Noise point
			&NamedPoint{"5", []float64{6, 8}},
			&NamedPoint{"6", []float64{6, 5}},
			&NamedPoint{"7", []float64{8, 4}},
			&NamedPoint{"8", []float64{2, 5}},
			&NamedPoint{"9", []float64{3, 7}},
		}
		clusterer = NewDBSCANClusterer(2.0, 2)
		result    [][]ClusterablePoint

		expectedClusters = map[ClusterablePoint]int{
			// Cluster 0
			data[0]: 4,
			data[2]: 4,
			data[8]: 4,
			data[9]: 4,

			// Cluster 1
			data[1]: 5,
			data[3]: 5,
			data[4]: 5,
			data[6]: 5,
			data[7]: 5,
		}
		compare = func(result [][]ClusterablePoint, expectedClusters map[ClusterablePoint]int) {
			for _, cluster := range result {
				size := len(cluster)
				for _, p := range cluster {
					if expectedClusters[p] != size {
						t.Fatalf("Expected point %v to be in cluster %v, got %v", p, expectedClusters[p], size)
					}
				}
			}
		}
	)

	result = clusterer.Cluster(data)
	compare(result, expectedClusters)

	clusterer.AutoSelectDimension = false
	clusterer.SetEps(7.0)

	expectedClusters = make(map[ClusterablePoint]int)
	for _, v := range data {
		expectedClusters[v] = len(data)
	}

	result = clusterer.Cluster(data)
	compare(result, expectedClusters)
}

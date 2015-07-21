package dbscan

import (
	// "bufio"
	// "io"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	"testing"
)

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

/*var (
	blockSize = 64 * 1024 * 1024
	LineDelim = byte('\n')
)

func Test_DBSCAN2(t *testing.T) {
	var (
		filename = "intermediate-92.txt"

		data   []ClusterablePoint = NamedPointToClusterablePoint(transformToPoints(t, filename))
		result [][]ClusterablePoint

		clusterer = NewDBSCANClusterer(1.465, 10)

		printCluster = func(result [][]ClusterablePoint) {
			for i, cluster := range result {
				t.Logf("Cluster %d: %d items", i, len(cluster))
			}
		}
	)

	t.Logf("Eps = %f", clusterer.GetEps())
	result = clusterer.Cluster(data)
	printCluster(result)
}

func transformToPoints(t *testing.T, filename string) []*NamedPoint {
	var (
		fileHandle, err = os.Open(filepath.Clean(filename))
		reader          = bufio.NewReaderSize(fileHandle, blockSize)
		line            string
		records         = make([]*NamedPoint, 0, 10000)
	)
	defer fileHandle.Close()
	if err != nil {
		t.Fatalf(err.Error())
		return nil
	}

	line, err = reader.ReadString(LineDelim)
	line = strings.TrimSpace(line)

	for _, name := range strings.Split(line, ",")[2:] {
		records = append(records, NewNamedPoint(
			name,
			make([]float64, 0, 10000),
		))
	}

	for err != io.EOF {
		line, err = reader.ReadString(LineDelim)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for i, v := range strings.Split(line, ",")[2:] {
			floatValue, _ := strconv.ParseFloat(v, 64)
			records[i].Point = append(records[i].Point, floatValue)
		}
	}

	return records
}*/

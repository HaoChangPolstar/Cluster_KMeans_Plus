package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Point struct {
	X float64
	Y float64
}

type Squares struct {
	plotter.XYs
}

type Location struct {
	Lat float64
	Lng float64
}

func main() {
	points := []Point{
		{X: 1, Y: 5},
		{X: 2, Y: 5},
		{X: 3, Y: 5},
		{X: 4, Y: 5},
		{X: 5, Y: 5},
		{X: 6, Y: 5},
		{X: 7, Y: 5},
		{X: 8, Y: 5},
		{X: 9, Y: 5},
		{X: 1, Y: 6},
		{X: 2, Y: 6},
		{X: 3, Y: 6},
		{X: 4, Y: 6},
		{X: 5, Y: 6},
		{X: 6, Y: 6},
		{X: 7, Y: 6},
		{X: 8, Y: 6},
		{X: 9, Y: 6},
		{X: 1, Y: 7},
		{X: 2, Y: 7},
		{X: 3, Y: 7},
		{X: 4, Y: 7},
		{X: 5, Y: 7},
		{X: 6, Y: 7},
		{X: 7, Y: 7},
		{X: 8, Y: 7},
		{X: 9, Y: 7},
	}
	// points := []Point{
	// 	{X: 1, Y: 1},
	// 	{X: 2, Y: 1},
	// 	{X: 1, Y: 2},
	// 	{X: 2, Y: 2},
	// 	{X: 5, Y: 6},
	// 	{X: 6, Y: 5},
	// 	{X: 5, Y: 5},
	// 	{X: 6, Y: 6},
	// 	{X: 9, Y: 8},
	// 	{X: 8, Y: 9},
	// 	{X: 8, Y: 8},
	// 	{X: 9, Y: 9},
	// 	{X: 19, Y: 18},
	// 	{X: 18, Y: 19},
	// 	{X: 18, Y: 18},
	// 	{X: 19, Y: 19},
	// }

	clusterCount := 3

	var colors []color.RGBA
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < clusterCount; i++ {
		rValue := uint8(rand.Intn(255))
		gValue := uint8(rand.Intn(255))
		bValue := uint8(rand.Intn(255))
		aValue := uint8(rand.Intn(255))
		newColor := color.RGBA{R: rValue, G: gValue, B: bValue, A: aValue}
		colors = append(colors, newColor)
	}

	startPointsArr := generateKMeansPlusStartPoint(clusterCount, points)

	// fmt.Println(startPointsArr)

	p := plot.New()
	p.Title.Text = "Squares"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	spoints, _ := plotter.NewScatter(generatePlottingPoints(startPointsArr))
	spoints.Color = color.RGBA{R: 255, G: 33, B: 43, A: 100}
	p.Add(spoints)
	p.Save(4*vg.Inch, 4*vg.Inch, "kMeansStartPoint.png")

	groupMaps := kMeansClustering(points, startPointsArr)
	log.Println(groupMaps[len(groupMaps)-1])
	for i := 0; i < len(groupMaps); i++ {
		groupMap := groupMaps[i]
		counter := 0
		p := plot.New()
		p.Title.Text = "Squares"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		for _, v := range groupMap {
			points, _ := plotter.NewScatter(generatePlottingPoints(v))
			points.Color = colors[counter]
			p.Add(points)
			counter++
		}
		p.Save(4*vg.Inch, 4*vg.Inch, strconv.Itoa(i)+"kmeans.png")
	}
}

func nnClustering(points []Point, startPointsArr []Point) (map[Point][]Point, []Point) {
	groupMap := map[Point][]Point{}
	centroids := []Point{}
	for i := 0; i < len(points); i++ {
		minDistance := math.MaxFloat64
		var minGroupIndex int
		for j := 0; j < len(startPointsArr); j++ {
			xDistance := math.Abs(points[i].X - startPointsArr[j].X)
			yDistance := math.Abs(points[i].Y - startPointsArr[j].Y)
			distance := math.Sqrt(xDistance*xDistance + yDistance*yDistance)
			if distance < minDistance {
				minDistance = distance
				minGroupIndex = j
			}
		}
		groupMap[startPointsArr[minGroupIndex]] = append(groupMap[startPointsArr[minGroupIndex]], points[i])
	}
	for _, v := range groupMap {
		var centroid Point
		var totalX float64
		var totalY float64
		totalPoint := len(v)
		for i := 0; i < totalPoint; i++ {
			totalX = totalX + v[i].X
			totalY = totalY + v[i].Y
		}
		centroid.X = totalX / float64(totalPoint)
		centroid.Y = totalY / float64(totalPoint)
		centroids = append(centroids, centroid)
	}
	return groupMap, centroids
}

func checkIsCentroidsRepeat(centroids []Point, groupMap map[Point][]Point) bool {
	isRepeat := false
	repeatCount := 0
	for i := 0; i < len(centroids); i++ {
		if len(groupMap[centroids[i]]) != 0 {
			repeatCount = repeatCount + 1
		}
	}
	if repeatCount == len(centroids) {
		isRepeat = true
	}
	return isRepeat
}

func kMeansClustering(points []Point, startPointsArr []Point) []map[Point][]Point {
	groupMaps := []map[Point][]Point{}
	newCentroids := make([]Point, len(startPointsArr))
	copy(newCentroids, startPointsArr)
	isCentroidsRepeat := false
	for !isCentroidsRepeat {
		groupMap, centroids := nnClustering(points, newCentroids)
		newCentroids = centroids
		groupMaps = append(groupMaps, groupMap)
		isCentroidsRepeat = checkIsCentroidsRepeat(centroids, groupMap)
	}
	return groupMaps
}

func generateRandomStartPoint(clusterCount int, points []Point) []Point {
	rand.Seed(time.Now().UnixNano())
	randomStartPointArray := make([]Point, clusterCount)
	repeatRecordTable := map[int]bool{}
	for i := 0; i < clusterCount; i++ {
		randomStartPoint := rand.Intn(len(points))
		if repeatRecordTable[randomStartPoint] {
			i = i - 1
		} else {
			repeatRecordTable[randomStartPoint] = true
			randomStartPointArray[i] = points[randomStartPoint]
		}
	}
	return randomStartPointArray
}

func generateKMeansPlusStartPoint(clusterCount int, points []Point) []Point {
	kMeansplusStartPointArray := []Point{}
	seedPoint := calculateMeanPoint(points)
	remainPoints := make([]Point, len(points))
	copy(remainPoints, points)
	for i := 0; i < clusterCount; i++ {
		var distanceArr []float64
		if i == 0 {
			distanceArr = findSeedToAllPointsDistance(seedPoint, remainPoints)
		} else {
			distanceArr = findShortestToAllPointDistance(kMeansplusStartPointArray, remainPoints)
		}
		maxValueIndex := findMax(distanceArr)
		seedPoint = remainPoints[maxValueIndex]
		kMeansplusStartPointArray = append(kMeansplusStartPointArray, remainPoints[maxValueIndex])
		remainPoints = append(remainPoints[:maxValueIndex], remainPoints[maxValueIndex+1:]...)
	}
	return kMeansplusStartPointArray
}

func findMax(arr []float64) int {
	var maxValue float64
	var maxValueIndex int
	for i := 0; i < len(arr); i++ {
		if arr[i] > maxValue {
			maxValue = arr[i]
			maxValueIndex = i
		}
	}
	return maxValueIndex
}

func calculateMeanPoint(points []Point) Point {
	var resultPoint Point
	var totalX float64
	var totalY float64
	var resultX float64
	var resultY float64
	for i := 0; i < len(points); i++ {
		totalX = totalX + points[i].X
		totalY = totalY + points[i].Y
	}
	resultX = totalX / float64(len(points))
	resultY = totalY / float64(len(points))
	resultPoint.X = resultX
	resultPoint.Y = resultY
	return resultPoint
}

func findShortestToAllPointDistance(kMeansplusStartPointArray []Point, remainPoints []Point) []float64 {
	var remainPointsShortestDistanceArr []float64
	for i := 0; i < len(remainPoints); i++ {
		shortestDistance := math.MaxFloat64
		for j := 0; j < len(kMeansplusStartPointArray); j++ {
			distance := findDistance(kMeansplusStartPointArray[j], remainPoints[i])
			if distance < shortestDistance {
				shortestDistance = distance
			}
		}
		remainPointsShortestDistanceArr = append(remainPointsShortestDistanceArr, shortestDistance)
	}
	return remainPointsShortestDistanceArr
}

func findSeedToAllPointsDistance(refPoint Point, points []Point) []float64 {
	distanceArr := make([]float64, len(points))
	var totalDistance float64
	for i := 0; i < len(points); i++ {
		refToPointDistance := findDistance(refPoint, points[i])
		distanceArr[i] = refToPointDistance
		totalDistance = totalDistance + refToPointDistance
	}
	return distanceArr
}

func findDistance(startPoint, endPoint Point) float64 {
	xDistance := startPoint.X - endPoint.X
	yDistance := startPoint.Y - endPoint.Y
	distance := math.Sqrt(xDistance*xDistance + yDistance*yDistance)
	return distance
}

func generatePlottingPoints(points []Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := 0; i < len(points); i++ {
		pts[i].X = points[i].X
		pts[i].Y = points[i].Y
	}
	return pts
}

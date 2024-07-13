
// these functions are for rotating shapes around a center point
func newVector(points []float32, degrees float32) []float32 {
	if math.Mod(float64(len(points)), 3) != 0 {
		panic("Your points dont have a multiple of 3(X/Y/Z) coordinates")
	}
	//todo loop through points in stead of assuming a triangle

	center := []float32{0, 0, 0}
	point1 := points[0:2]
	point2 := points[3:5]
	point3 := points[6:8]

	point1 = rotatePoint(point1, center, degrees)
	point2 = rotatePoint(point2, center, degrees)
	point3 = rotatePoint(point3, center, degrees)
	newSlice := append(point1, point2...)
	newSlice = append(newSlice, point3...)
	return newSlice
}

func rotatePoint(point, center []float32, angleInDegrees float32) []float32 {
	radians := angleInDegrees * (math.Pi / 180)
	cosTheta := float32(math.Cos(float64(radians)))
	sinTheta := float32(math.Sin(float64(radians)))
	newPoint := make([]float32, 3)
	newPoint[0] = (cosTheta*(point[0]-center[0]) -
		sinTheta*(point[1]-center[1]) + center[0])
	newPoint[1] = (sinTheta*(point[0]-center[0]) +
		cosTheta*(point[1]-center[1]) + center[1])
	return newPoint
}

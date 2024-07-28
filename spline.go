package turdgl

// GenerateCatmullRomSpline generates a series of points on a Catmull-Rom spline that passes
// through the given points.
func GenerateCatmullRomSpline(points []Vec, steps int) []Vec {
	n := len(points)
	if n < 4 {
		return nil
	}
	splinePoints := []Vec{}
	// First point
	for j := 0; j <= steps; j++ {
		t := float64(j) / float64(steps)
		splinePoints = append(splinePoints,
			catmullRomSpline(points[0], points[0], points[1], points[2], t))
	}
	// Middle points
	for i := 1; i < n-2; i++ {
		for j := 0; j <= steps; j++ {
			t := float64(j) / float64(steps)
			splinePoints = append(splinePoints,
				catmullRomSpline(points[i-1], points[i], points[i+1], points[i+2], t))
		}
	}
	// Last point
	for j := 0; j <= steps; j++ {
		t := float64(j) / float64(steps)
		splinePoints = append(splinePoints,
			catmullRomSpline(points[n-3], points[n-2], points[n-1], points[n-1], t))
	}
	return splinePoints
}

// catmullRomSpline calculates a point on a Catmull-Rom spline given four control points
// and a parameter t (0 <= t <= 1). Returns the interpolated point on the spline.
func catmullRomSpline(p0, p1, p2, p3 Vec, t float64) Vec {
	t2 := t * t
	t3 := t2 * t

	f0 := -0.5*t3 + t2 - 0.5*t
	f1 := 1.5*t3 - 2.5*t2 + 1.0
	f2 := -1.5*t3 + 2.0*t2 + 0.5*t
	f3 := 0.5*t3 - 0.5*t2

	x := f0*p0.X + f1*p1.X + f2*p2.X + f3*p3.X
	y := f0*p0.Y + f1*p1.Y + f2*p2.Y + f3*p3.Y

	return Vec{x, y}
}

// Copyright 2013, 2014 Peter Vasil. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpx

import (
	"math"
	"sort"
)

const oneDegree = 1000.0 * 10000.8 / 90.0
const earthRadius = 6371 * 1000

// toRad converts degree to radian
func toRad(x float64) float64 {
	return x / 180. * math.Pi
}

// HaversineDistance returns the haversine distance between two points.
//
// Implemented from http://www.movable-type.co.uk/scripts/latlong.html
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRad(lat1 - lat2)
	dLon := toRad(lon1 - lon2)
	thisLat1 := toRad(lat1)
	thisLat2 := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c

	return d
}

// calcMaxSpeed calculates the max speed
func calcMaxSpeed(speedsDistances []speedsAndDistances) float64 {
	lenArrs := len(speedsDistances)

	if len(speedsDistances) < 20 {
		//log.Println("Segment too small to compute speed, size: ", lenArrs)
		return 0.0
	}

	var sumDists float64
	for _, d := range speedsDistances {
		sumDists += d.distance
	}
	averageDist := sumDists / float64(lenArrs)

	var variance float64
	for i := 0; i < len(speedsDistances); i++ {
		variance += math.Pow(speedsDistances[i].distance-averageDist, 2)
	}
	stdDeviation := math.Sqrt(variance)

	// ignore items with distance too long
	var filteredSD []speedsAndDistances
	for i := 0; i < len(speedsDistances); i++ {
		dist := math.Abs(speedsDistances[i].distance - averageDist)
		if dist <= stdDeviation*1.5 {
			filteredSD = append(filteredSD, speedsDistances[i])
		}
	}

	speeds := make([]float64, len(filteredSD))
	for i, sd := range filteredSD {
		speeds[i] = sd.speed
	}

	speedsSorted := sort.Float64Slice(speeds)

	maxIdx := int(float64(len(speedsSorted)) * 0.95)
	if maxIdx >= len(speedsSorted) {
		maxIdx = len(speedsSorted) - 1
	}
	return speedsSorted[maxIdx]
}

// calcUphillDownhill calculates uphill/downhill data
func calcUphillDownhill(elevations []float64) (float64, float64) {
	elevsLen := len(elevations)
	if elevsLen == 0 {
		return 0.0, 0.0
	}

	smoothElevations := make([]float64, elevsLen)

	for i, elev := range elevations {
		var currEle float64
		if 0 < i && i < elevsLen-1 {
			prevEle := elevations[i-1]
			nextEle := elevations[i+1]
			currEle = prevEle*0.3 + elev*0.4 + nextEle*0.3
		} else {
			currEle = elev
		}
		smoothElevations[i] = currEle
	}

	var uphill float64
	var downhill float64

	for i := 1; i < len(smoothElevations); i++ {
		d := smoothElevations[i] - smoothElevations[i-1]
		if d > 0.0 {
			uphill += d
		} else {
			downhill -= d
		}
	}

	return uphill, downhill
}

func distance(lat1, lon1, ele1, lat2, lon2, ele2 float64, threeD, haversine bool) float64 {

	absLat := math.Abs(lat1 - lat2)
	absLon := math.Abs(lon1 - lon2)
	if haversine || absLat > 0.2 || absLon > 0.2 {
		return haversineDistance(lat1, lon1, lat2, lon2)
	}

	coef := math.Cos(toRad(lat1))
	x := lat1 - lat2
	y := (lon1 - lon2) * coef

	distance2d := math.Sqrt(x*x+y*y) * oneDegree

	if !threeD || ele1 == ele2 {
		return distance2d
	}

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow((ele1-ele2), 2))
}

// ElevationAngle calculates the elavation angle
func ElevationAngle(loc1, loc2 *Wpt, radians bool) float64 {
	b := loc2.Ele - loc1.Ele
	a := loc2.Distance2D(loc1)

	if a == 0.0 {
		return 0.0
	}

	angle := math.Atan(b / a)

	if radians {
		return angle
	}

	return 180 * angle / math.Pi
}

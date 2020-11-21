package stravaminhashlsh

import "math"

func deg2Num(lat, lon float64, zoom int32) (x, y int) {
	x = int(math.Floor((lon + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))
	return x, y
}

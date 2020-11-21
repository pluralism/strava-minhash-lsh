package stravaminhashlsh

import "math"

type bresenhamResult struct {
	inverted bool
	array    [][]int
}

func bresenham(x1, y1, x2, y2 int) bresenhamResult {
	result := bresenhamResult{
		inverted: false,
		array:    [][]int{},
	}

	dx := x2 - x1
	dy := y2 - y1

	dx1 := int(math.Abs(float64(dx)))
	dy1 := int(math.Abs(float64(dy)))

	px := 2*dy1 - dx1
	py := 2*dx1 - dy1

	var x, y, xe, ye int

	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
			result.inverted = true
		}

		result.array = append(result.array, []int{x, y})

		for i := 0; x < xe; i++ {
			x++

			if px < 0 {
				px = px + 2*dy1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y++
				} else {
					y--
				}
				px = px + 2*(dy1-dx1)
			}

			if result.inverted {
				result.array = append([][]int{{x, y}}, result.array...)
			} else {
				result.array = append(result.array, []int{x, y})
			}
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
			result.inverted = true
		}

		result.array = append(result.array, []int{x, y})

		for i := 0; y < ye; i++ {
			y++

			if py <= 0 {
				py = py + 2*dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x++
				} else {
					x--
				}
				py = py + 2*(dx1-dy1)
			}

			if result.inverted {
				result.array = append([][]int{{x, y}}, result.array...)
			} else {
				result.array = append(result.array, []int{x, y})
			}
		}
	}

	return result
}

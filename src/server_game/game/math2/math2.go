package math2

import "math"

const (
	PI float32 = 3.14159265358979323846
)

func Angle2Matrix(angles [3]float32) [3][4]float32 {
	angle := angles[2] * (PI * 2 / 360)
	sy := float32(math.Sin(float64(angle)))
	cy := float32(math.Cos(float64(angle)))
	angle = angles[1] * (PI * 2 / 360)
	sp := float32(math.Sin(float64(angle)))
	cp := float32(math.Cos(float64(angle)))
	angle = angles[0] * (PI * 2 / 360)
	sr := float32(math.Sin(float64(angle)))
	cr := float32(math.Cos(float64(angle)))
	return [3][4]float32{
		{cp * cy, cp * sy, -sp, 0.0},
		{sr*sp*cy + cr*(-sy), sr*sp*sy + cr*cy, sr * cp, 0.0},
		{cr*sp*cy + (-sr)*(-sy), cr*sp*sy + (-sr)*cy, cr * cp, 0.0},
	}
}

func VectorRotate(in1 [3]float32, in2 [3][4]float32) [3]float32 {
	return [3]float32{
		in1[0]*in2[0][0] + in1[1]*in2[0][1] + in1[2]*in2[0][2],
		in1[0]*in2[1][0] + in1[1]*in2[1][1] + in1[2]*in2[1][2],
		in1[0]*in2[2][0] + in1[1]*in2[2][1] + in1[2]*in2[2][2],
	}
}

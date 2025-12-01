package vec

type SquareMat2i struct {
	A11, A12, A21, A22 int
}

// https://en.wikipedia.org/wiki/Rotation_matrix
// right-handed coordinate system counterclockwise
func NewRotCCW() SquareMat2i {
	return SquareMat2i{
		A11: 0,
		A12: 1,
		A21: -1,
		A22: 0,
	}
}

// right-handed coordinate system clockwise
func NewRotCW() SquareMat2i {
	return SquareMat2i{
		A11: 0,
		A12: -1,
		A21: 1,
		A22: 0,
	}
}

func (s SquareMat2i) Mul(i Vec2i) Vec2i {
	return Vec2i{
		X: s.A11*i.X + s.A12*i.Y,
		Y: s.A21*i.X + s.A22*i.Y,
	}
}

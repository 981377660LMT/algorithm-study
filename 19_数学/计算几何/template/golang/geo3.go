// 三维几何
// Point: 点
// Line: 直线

package main

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | Float
}
type Point3D[T Number] struct {
	x, y, z T
}

func NewPoint3D[T Number](x, y, z T) Point3D[T] { return Point3D[T]{x, y, z} }

func (p Point3D[T]) Add(q Point3D[T]) Point3D[T] { return Point3D[T]{p.x + q.x, p.y + q.y, p.z + q.z} }
func (p Point3D[T]) Sub(q Point3D[T]) Point3D[T] { return Point3D[T]{p.x - q.x, p.y - q.y, p.z - q.z} }

func (p Point3D[T]) Dot(q Point3D[T]) T { return p.x*q.x + p.y*q.y + p.z*q.z }
func (p Point3D[T]) IsParallel(q Point3D[T]) bool {
	return p.y*q.z == p.z*q.y && p.z*q.x == p.x*q.z && p.x*q.y == p.y*q.x
}

func (p Point3D[T]) Cross(q Point3D[T]) Point3D[T] {
	return Point3D[T]{p.y*q.z - p.z*q.y, p.z*q.x - p.x*q.z, p.x*q.y - p.y*q.x}
}

type Line3D[T Number] struct {
	// a + td
	a, d Point3D[T]
}

func NewLine3D[T Number](a, b Point3D[T]) Line3D[T] {
	return Line3D[T]{a, b.Sub(a)}
}

func (l Line3D[T]) IsParallel(other Line3D[T]) bool {
	n := l.d.Cross(other.d)
	return n.x == 0 && n.y == 0 && n.z == 0
}

func (l Line3D[T]) Contain(p Point3D[T]) bool {
	p = p.Sub(l.a)
	p = p.Cross(l.d)
	return p.x == 0 && p.y == 0 && p.z == 0
}

// 0:无交点, 1:有唯一交点, 2:有两个或两个以上交点
func CountCross[T Integer](L1, L2 Line3D[T]) int32 {
	if L1.IsParallel(L2) {
		if L1.Contain(L2.a) {
			return 2
		}
		return 0
	}
	norm := L1.d.Cross(L2.d)
	if (L1.a.Sub(L2.a)).Dot(norm) == 0 {
		return 1
	}
	return 0
}

// 仅当 CountCross 返回 1 时调用.
func CrossPoint[T Number](line1, line2 Line3D[T]) Point3D[T] {
	d1 := line1.d
	d2 := line2.d
	a := line2.a.Sub(line1.a)
	t1 := func() T {
		for i := 0; i < 3; i++ {
			d1 = NewPoint3D[T](d1.y, d1.z, d1.x)
			d2 = NewPoint3D[T](d2.y, d2.z, d2.x)
			a = NewPoint3D[T](a.y, a.z, a.x)
			det := d1.x*d2.y - d1.y*d2.x
			if det != 0 {
				return T(a.x*d2.y-a.y*d2.x) / T(det)
			}
		}
		panic("unreachable")
	}()
	x := line1.a.x + t1*line1.d.x
	y := line1.a.y + t1*line1.d.y
	z := line1.a.z + t1*line1.d.z
	return NewPoint3D[T](x, y, z)
}

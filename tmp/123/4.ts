import { RectangleSum } from '../../6_tree/树状数组/经典题/RectangleSum-MergeSegtree'

export {}

const INF = 2e9 // !超过int32使用2e15

function foo() {}

// func numberOfPairs(points [][]int) int {
// 	S := NewStaticRectangleCount()
// 	for _, p := range points {
// 		S.AddPoint(p[0], p[1], 1)
// 	}
// 	S.Build()
// 	res := 0
// 	for i, p1 := range points { // liupengsay
// 		x1, y1 := p1[0], p1[1]
// 		for j, p2 := range points { // 小羊肖恩
// 			x2, y2 := p2[0], p2[1]
// 			if i != j && x1 <= x2 && y1 >= y2 && S.Query(x1, x2+1, y2, y1+1) == 2 {
// 				res++
// 			}
// 		}
// 	}
// 	return res
// }
function numberOfPairs(points: number[][]): number {
  const S = new RectangleSum()
  points.forEach(p => {
    S.addPoint(p[0], p[1], 1)
  })
  S.build()
}

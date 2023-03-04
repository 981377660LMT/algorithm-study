 main() {
	grid := NewSegmentTree2D(3, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			grid.Update(i, j, i, j, i*3+j)
			fmt.Println(grid.Get(i, j))
		}
	}
}

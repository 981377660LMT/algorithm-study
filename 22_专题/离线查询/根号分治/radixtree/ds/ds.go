func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }

	tree := NewRadixTree(e, op, 2)
	tree.Build(10, func(i int) int { return i + 1 })

	// Querying the range [2, 5)
	fmt.Println("Query range [2, 5):", tree.QueryRange(2, 5))

	// Get all values
	fmt.Println("Get All:", tree.GetAll())

	// Update value at index 3
	tree.Update(3, 5)
	fmt.Println("After update at index 3:", tree.GetAll())

	// Query All
	fmt.Println("Query all:", tree.QueryAll())

	// Set value at index 6
	tree.Set(6, 10)
	fmt.Println("After set at index 6:", tree.GetAll())

	// Querying the range [2, 5)
	fmt.Println("Query range [2, 5):", tree.QueryRange(2, 5))

	tree.Build(10, func(i int) int { return i + 1 })
	fmt.Println("Query all:", tree.GetAll())
	// min left
	fmt.Println("Min left:", tree.MinLeft(10, func(x int) bool { return x < 27 }))
	// max right
	fmt.Println("Max right:", tree.MaxRight(0, func(x int) bool { return x <= 6 }))
}

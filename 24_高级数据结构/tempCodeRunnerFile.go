
	nums := NewDynamicSequence([]int{})
	time1 := time.Now()
	for i := 0; i < 2e5; i++ {
		nums.Append(i)
		nums.Append(i)
		nums.At(0)
		nums.Insert(i, 100)
		nums.Update(0, i, 10)
		nums.Query(0, i)
		nums.QueryAll()
		nums.Pop(-1)
		nums.Size()
		nums.Erase(0, 1)
		nums.Reverse(0, i)
	}
	fmt.Println(time.Since(time1))

	nums2 := NewDynamicSequence([]int{})
	nums2.Append(1)
	nums2.Append(1)
	nums2.Append(1)
	nums2.Append(1)
	nums2.Set(-1, 2)
	fmt.Println(nums2)
}
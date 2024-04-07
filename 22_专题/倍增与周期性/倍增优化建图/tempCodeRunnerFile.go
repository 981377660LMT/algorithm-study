True(6, func(i int32) bool { return lca.Depth[i] >= 1 })
	expect[int32](step2, 2)
	expect[int32](to2, 1)

			// 	g := groups[groupIds[cur]]
			// 	ptr := &groupsPtr[1][groupIds[cur]]
			// 	for *ptr >= 0 {
			// 		last := g[*ptr]
			// 		*ptr--
			// 		if id := vid(last, true); used[1][id] {
			// 			continue
			// 		} else {
			// 			return int(id)
			// 		}
			// 	}
			// }
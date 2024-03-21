		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			sb[i], sb[j] = sb[j], sb[i]
		}
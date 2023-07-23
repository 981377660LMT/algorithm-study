		first, second := SplitByRank(root.right, k-leftSize-1)
		root.right = first
		return _pushUp(root), second

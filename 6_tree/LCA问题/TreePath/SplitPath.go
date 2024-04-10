// 路径分割/分割路径/路径拆分/拆分路径

package main

// 将路径 `from->to` 以 `separator` 为分隔，按顺序分成两段.
// `separtor` 必须在路径上.
// 返回两段路径的起点和终点.如果某段路径为空，则起点和终点都为-1.
func SplitPath(
	from, to int32, separator int32,
	depth []int32,
	kthAncestorFn func(node, k int32) int32,
	lcaFn func(node1, node2 int32) int32,
) (from1, to1, from2, to2 int32) {
	from1, to1, from2, to2 = -1, -1, -1, -1
	if from == to {
		return
	}

	down, top := from, to
	swapped := false
	if depth[down] < depth[top] {
		down, top = top, down
		swapped = true
	}

	lca := lcaFn(from, to)
	if lca == top {
		// down和top在一条链上.
		if separator == down {
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
		} else {
			from1 = down
			to1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		}
	} else {
		// down和top在lca两个子树上.
		if separator == down {
			from2 = kthAncestorFn(separator, 1)
			to2 = top
		} else if separator == top {
			from1 = down
			to1 = kthAncestorFn(separator, 1)
		} else {
			var jump1, jump2 int32
			if separator == lca {
				jump1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
				jump2 = kthAncestorFn(top, depth[top]-depth[separator]-1)
			} else if lcaFn(separator, down) == separator {
				jump1 = kthAncestorFn(down, depth[down]-depth[separator]-1)
				jump2 = kthAncestorFn(separator, 1)
			} else {
				jump1 = kthAncestorFn(separator, 1)
				jump2 = kthAncestorFn(top, depth[top]-depth[separator]-1)
			}
			from1 = down
			to1 = jump1
			from2 = jump2
			to2 = top
		}
	}

	if swapped {
		from1, to1, from2, to2 = to2, from2, to1, from1
	}
	return
}

func SplitPathByJump(
	from, to, separator int32,
	jump func(start, target, step int32) int32,
) (from1, to1, from2, to2 int32) {
	from1, to1, from2, to2 = -1, -1, -1, -1

	if from == to {
		return
	}

	if separator == from {
		from2 = jump(from, to, 1)
		to2 = to
		return
	}

	if separator == to {
		from1 = from
		to1 = jump(to, from, 1)
		return
	}

	from1 = from
	to1 = jump(separator, from, 1)
	from2 = jump(separator, to, 1)
	to2 = to
	return
}

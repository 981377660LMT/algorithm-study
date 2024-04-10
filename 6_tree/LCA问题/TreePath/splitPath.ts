/**
 * 将路径 `from->to` 以 `separator` 为分隔，按顺序分成两段.
 * `separtor` 必须在路径上.
 * 返回两段路径的起点和终点.如果某段路径为空，则起点和终点都为`undefined`.
 */
function splitPath(
  from: number,
  to: number,
  separator: number,
  treeProps: {
    depth: ArrayLike<number>
    kthAncestorFn: (node: number, k: number) => number
    lcaFn: (node1: number, node2: number) => number
  }
): {
  from1: number | undefined
  to1: number | undefined
  from2: number | undefined
  to2: number | undefined
} {
  let from1: number | undefined = undefined
  let to1: number | undefined = undefined
  let from2: number | undefined = undefined
  let to2: number | undefined = undefined
  if (from === to) {
    return { from1, to1, from2, to2 }
  }

  const { depth } = treeProps

  let down = from
  let top = to
  let swapped = false
  if (depth[down] < depth[top]) {
    const tmp = down
    down = top
    top = tmp
    swapped = true
  }

  const lca = treeProps.lcaFn(from, to)
  if (lca === top) {
    // down和top在一条链上.
    if (separator === down) {
      from2 = treeProps.kthAncestorFn(separator, 1)
      to2 = top
    } else if (separator === top) {
      from1 = down
      to1 = treeProps.kthAncestorFn(down, depth[down] - depth[separator] - 1)
    } else {
      from1 = down
      to1 = treeProps.kthAncestorFn(down, depth[down] - depth[separator] - 1)
      from2 = treeProps.kthAncestorFn(separator, 1)
      to2 = top
    }
  } else {
    // down和top在lca两个子树上.
    if (separator === down) {
      from2 = treeProps.kthAncestorFn(separator, 1)
      to2 = top
    } else if (separator === top) {
      from1 = down
      to1 = treeProps.kthAncestorFn(separator, 1)
    } else {
      let jump1: number
      let jump2: number
      if (separator === lca) {
        jump1 = treeProps.kthAncestorFn(down, depth[down] - depth[separator] - 1)
        jump2 = treeProps.kthAncestorFn(top, depth[top] - depth[separator] - 1)
      } else if (treeProps.lcaFn(separator, down) === separator) {
        jump1 = treeProps.kthAncestorFn(down, depth[down] - depth[separator] - 1)
        jump2 = treeProps.kthAncestorFn(separator, 1)
      } else {
        jump1 = treeProps.kthAncestorFn(separator, 1)
        jump2 = treeProps.kthAncestorFn(top, depth[top] - depth[separator] - 1)
      }
      from1 = down
      to1 = jump1
      from2 = jump2
      to2 = top
    }
  }

  if (swapped) {
    const tmpFrom1 = from1
    from1 = to2
    to2 = tmpFrom1
    const tmpTo1 = to1
    to1 = from2
    from2 = tmpTo1
  }

  return { from1, to1, from2, to2 }
}

/**
 * 将路径 `from->to` 以 `separator` 为分隔，按顺序分成两段.
 * `separtor` 必须在路径上.
 * 返回两段路径的起点和终点.如果某段路径为空，则起点和终点都为`undefined`.
 */
function splitPathByJump(
  from: number,
  to: number,
  separator: number,
  jump: (start: number, target: number, step: number) => number
): {
  from1: number | undefined
  to1: number | undefined
  from2: number | undefined
  to2: number | undefined
} {
  let from1: number | undefined = undefined
  let to1: number | undefined = undefined
  let from2: number | undefined = undefined
  let to2: number | undefined = undefined

  if (from === to) {
    return { from1, to1, from2, to2 }
  }

  if (separator === from) {
    from2 = jump(from, to, 1)
    to2 = to
    return { from1, to1, from2, to2 }
  }

  if (separator === to) {
    from1 = from
    to1 = jump(to, from, 1)
    return { from1, to1, from2, to2 }
  }

  from1 = from
  to1 = jump(separator, from, 1)
  from2 = jump(separator, to, 1)
  to2 = to
  return { from1, to1, from2, to2 }
}

export { splitPath, splitPathByJump }

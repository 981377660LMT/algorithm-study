/**
 * 还原路径/dp复原.
 */
function restorePath(target: number, pre: ArrayLike<number>): number[] {
  const path: number[] = [target]
  while (pre[path[path.length - 1]] !== -1) {
    path.push(pre[path[path.length - 1]])
  }
  return path.reverse()
}

export { restorePath }

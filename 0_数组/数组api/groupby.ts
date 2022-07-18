/**
 * @example
 * ```ts
 * const list = [1, 1, 2, 3, 3, 4, 4, 5, 5, 5]
 * groupBy(list) // [[1, 1], [2], [3, 3], [4, 4], [5, 5, 5]]
 * ```
 */
function groupBy<T>(iterable: Iterable<T>): T[][] {
  const groups: T[][] = []
  let pre: any

  for (const char of iterable) {
    if (char !== pre) {
      const newGroup = [char]
      groups.push(newGroup)
    } else {
      groups[groups.length - 1].push(char)
    }

    pre = char
  }

  return groups
}

if (require.main === module) {
  console.log(groupBy('abbcccdddd'))
}

export { groupBy }

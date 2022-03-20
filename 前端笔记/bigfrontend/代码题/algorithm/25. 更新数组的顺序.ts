/**
 * @param {any[]} items
 * @param {number[]} newOrder
 * @return {void}
 * 使用额外的O(n)空间很简单就能完成该题目，你能不实用额外空间完成该题目吗？
 */
function sort(items: any[], newOrder: number[]): void {
  // 有点像矩阵求逆过程中 两个矩阵同步操作 不看items newOrder变成[ 0, 1, 2, 3, 4, 5 ]就行
  for (let i = 0; i < newOrder.length; i++) {
    while (newOrder[i] !== i) {
      const j = newOrder[i]
      swap(newOrder, i, j)
      swap(items, i, j)
    }
  }

  function swap(arr: any[], i: number, j: number) {
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
}

console.log(sort(['A', 'B', 'C', 'D', 'E', 'F'], [1, 5, 4, 3, 2, 0]))
// 上述例子进行重排过后，应该得到如下结果
// ['F', 'A', 'E', 'D', 'C', 'B']

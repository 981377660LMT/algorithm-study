/**
 * 原地删除数组中满足条件的元素.
 */
function removeInplace(arr: unknown[], shouldRemove: (index: number) => boolean): void {
  let ptr = 0
  for (let i = 0; i < arr.length; i++) {
    if (!shouldRemove(i)) arr[ptr++] = arr[i]
  }
  arr.length = ptr
}

export { removeInplace }

if (require.main === module) {
  const arr = [1, 2, 2, 3, 3, 3, 4, 4, 4, 4]
  removeInplace(arr, i => arr[i] === 2 || arr[i] === 3)
  console.log(arr) // [1, 2, 3, 4]
}

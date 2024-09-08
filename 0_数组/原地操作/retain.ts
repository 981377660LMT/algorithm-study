/**
 * 保留数组中满足条件的元素.
 * 原地操作.
 */
function retain(arr: any[], f: (index: number) => boolean): void {
  let ptr = 0
  for (let i = 0; i < arr.length; i++) {
    if (f(i)) arr[ptr++] = arr[i]
  }
  arr.length = ptr
}

export { retain }

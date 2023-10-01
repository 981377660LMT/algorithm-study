function mergeTwoSortedArray(arr1: number[], arr2: number[]): number[]
function mergeTwoSortedArray(arr1: number[], arr2: number[], mergeTo: number[]): void
function mergeTwoSortedArray(arr1: number[], arr2: number[], mergeTo?: number[]): number[] | void {
  const n1 = arr1.length
  if (!n1) {
    if (mergeTo) {
      for (let i = 0; i < arr2.length; i++) mergeTo[i] = arr2[i]
    } else {
      return arr2
    }
  }

  const n2 = arr2.length
  if (!n2) {
    if (mergeTo) {
      for (let i = 0; i < arr1.length; i++) mergeTo[i] = arr1[i]
    } else {
      return arr1
    }
  }

  const res = mergeTo || Array(n1 + n2)
  let i = 0
  let j = 0
  let k = 0
  while (i < n1 && j < n2) {
    if (arr1[i] < arr2[j]) {
      res[k] = arr1[i]
      i++
    } else {
      res[k] = arr2[j]
      j++
    }
    k++
  }
  while (i < n1) {
    res[k] = arr1[i]
    i++
    k++
  }
  while (j < n2) {
    res[k] = arr2[j]
    j++
    k++
  }

  return mergeTo ? undefined : res
}

function mergeThreeSortedArray(arr1: number[], arr2: number[], arr3: number[]): number[]
function mergeThreeSortedArray(
  arr1: number[],
  arr2: number[],
  arr3: number[],
  mergeTo: number[]
): void
function mergeThreeSortedArray(
  arr1: number[],
  arr2: number[],
  arr3: number[],
  mergeTo?: number[]
): number[] | void {
  const n1 = arr1.length
  if (!n1) {
    if (mergeTo) {
      mergeTwoSortedArray(arr1, arr2, mergeTo)
    } else {
      return mergeTwoSortedArray(arr1, arr2)
    }
  }
  const n2 = arr2.length
  if (!n2) {
    if (mergeTo) {
      mergeTwoSortedArray(arr1, arr3, mergeTo)
    } else {
      return mergeTwoSortedArray(arr1, arr3)
    }
  }
  const n3 = arr3.length
  if (!n3) {
    if (mergeTo) {
      mergeTwoSortedArray(arr2, arr3, mergeTo)
    } else {
      return mergeTwoSortedArray(arr2, arr3)
    }
  }
  const res = mergeTo || Array(n1 + n2 + n3)
  let i1 = 0
  let i2 = 0
  let i3 = 0
  let k = 0
  while (i1 < n1 && i2 < n2 && i3 < n3) {
    if (arr1[i1] < arr2[i2]) {
      if (arr1[i1] < arr3[i3]) {
        res[k] = arr1[i1]
        i1++
      } else {
        res[k] = arr3[i3]
        i3++
      }
    } else if (arr2[i2] < arr3[i3]) {
      res[k] = arr2[i2]
      i2++
    } else {
      res[k] = arr3[i3]
      i3++
    }
    k++
  }
  while (i1 < n1 && i2 < n2) {
    if (arr1[i1] < arr2[i2]) {
      res[k] = arr1[i1]
      i1++
    } else {
      res[k] = arr2[i2]
      i2++
    }
    k++
  }
  while (i1 < n1 && i3 < n3) {
    if (arr1[i1] < arr3[i3]) {
      res[k] = arr1[i1]
      i1++
    } else {
      res[k] = arr3[i3]
      i3++
    }
    k++
  }
  while (i2 < n2 && i3 < n3) {
    if (arr2[i2] < arr3[i3]) {
      res[k] = arr2[i2]
      i2++
    } else {
      res[k] = arr3[i3]
      i3++
    }
    k++
  }
  while (i1 < n1) {
    res[k] = arr1[i1]
    i1++
    k++
  }
  while (i2 < n2) {
    res[k] = arr2[i2]
    i2++
    k++
  }
  while (i3 < n3) {
    res[k] = arr3[i3]
    i3++
    k++
  }
  return mergeTo ? undefined : res
}

function mergeKSortedArray(arrays: number[][]): number[] {
  const n = arrays.length
  if (!n) return []
  if (n === 1) return arrays[0]
  if (n === 2) return mergeTwoSortedArray(arrays[0], arrays[1])

  const merge = (start: number, end: number): number[] => {
    if (start >= end) return []
    if (end - start === 1) return arrays[start]
    const mid = (start + end) >>> 1
    return mergeTwoSortedArray(merge(start, mid), merge(mid, end))
  }

  return merge(0, arrays.length)
}

export { mergeTwoSortedArray, mergeThreeSortedArray, mergeKSortedArray }

if (require.main === module) {
  const tmp = Array(8)
  console.log(mergeTwoSortedArray([1, 3, 5, 7], [2, 4, 6, 8]))
  console.log(tmp)

  const tmp2 = Array(8)
  mergeThreeSortedArray([1, 3, 5, 7], [2, 4, 6, 8], [0, 9, 10, 11], tmp2)
  console.log(tmp2)
}

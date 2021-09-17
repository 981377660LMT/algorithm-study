// 给你个整数数组 arr，其中每个元素都 不相同。

// 请你找到所有具有最小绝对差的元素对，并且按升序的顺序返回。
/**
 * @param {number[]} arr
 * @return {number[][]}
 */
var minimumAbsDifference = function (arr) {
  const res = []
  arr.sort((a, b) => a - b)
  let min = Infinity
  for (let i = 1; i < arr.length; i++) {
    min = Math.min(arr[i] - arr[i - 1], min)
  }
  for (let i = 1; i < arr.length; i++) {
    arr[i] - arr[i - 1] === min && res.push([arr[i - 1], arr[i]])
  }

  return res
}

console.log(minimumAbsDifference([3, 8, -10, 23, 19, -4, -14, 27]))

/**
 * @param {number[]} encoded
 * @param {number} first
 * @return {number[]}
 * 经编码后变为长度为 n - 1 的另一个整数数组 encoded ，其中 encoded[i] = arr[i] XOR arr[i + 1]
 * 例如，arr = [1,0,2,1] 经编码后得到 encoded = [1,2,3] 。
 * 请解码返回原数组 arr 。可以证明答案存在并且是唯一的。
 * @summary
 * res[i] ^ encoded[i]就是原数组的下一个位置。
 */
var decode = function (encoded, first) {
  const res = [first]
  for (let i = 0; i < encoded.length; i++) {
    res.push(encoded[i] ^ res[i])
  }
  return res
}

console.log(decode([1, 2, 3], 1))

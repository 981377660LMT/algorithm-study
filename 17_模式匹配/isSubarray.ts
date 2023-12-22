/**
 * O(n+m) 判断`shorter`是否是`longer`的子串.
 */
function isSubarray<T extends ArrayLike<unknown>>(longer: T, shorter: T): boolean {
  if (shorter.length > longer.length) return false
  if (!shorter.length) return true
  const n = longer.length
  const m = shorter.length
  const st = Array(m + n)
  for (let i = 0; i < m; i++) st[i] = shorter[i]
  for (let i = 0; i < n; i++) st[i + m] = longer[i]
  const z = zAlgo(st)
  for (let i = m; i < n + m; i++) {
    if (z[i] >= m) return true
  }
  return false
}

function zAlgo(seq: ArrayLike<unknown>): Uint32Array {
  const n = seq.length
  if (!n) return new Uint32Array(0)
  const z = new Uint32Array(n)
  let j = 0
  for (let i = 1; i < n; i++) {
    let k = 0
    if (j + z[j] > i) {
      k = Math.min(j + z[j] - i, z[i - j])
    }
    while (i + k < n && seq[k] === seq[i + k]) k++
    if (j + z[j] < i + z[i]) j = i
    z[i] = k
  }
  z[0] = n
  return z
}

export { isSubarray }

if (require.main === module) {
  console.log(isSubarray('abc', 'ab'))
}

// // 类欧几里得算法
// // 商求和(divSum)
// // ∑⌊(ai+b)/m⌋, i in [0,n-1]
// // https://oi-wiki.org/math/euclidean/
// func FloorSum(n, m, a, b int) (res int) {
// 	if a < 0 {
// 		a2 := a%m + m
// 		res -= n * (n - 1) / 2 * ((a2 - a) / m)
// 		a = a2
// 	}
// 	if b < 0 {
// 		b2 := b%m + m
// 		res -= n * ((b2 - b) / m)
// 		b = b2
// 	}
// 	for {
// 		if a >= m {
// 			res += n * (n - 1) / 2 * (a / m)
// 			a %= m
// 		}
// 		if b >= m {
// 			res += n * (b / m)
// 			b %= m
// 		}
// 		yMax := a*n + b
// 		if yMax < m {
// 			break
// 		}
// 		n = yMax / m
// 		b = yMax % m
// 		m, a = a, m
// 	}
// 	return
// }

// // 余数求和(ModSum/remainderSum)
// // ∑k%i (i in [1,n]), 即 sum(k%i for i in [1,n])
// // = ∑k-(k/i)*i
// // = n*k-∑(k/i)*i
// // 对于 [l,r] 范围内的 i，k/i 不变，此时 ∑(k/i)*i = (k/i)*∑i = (k/i)*(l+r)*(r-l+1)/2
// func ModSum(n int, k int) int {
// 	sum := n * k
// 	for l, r := 1, 0; l <= n; l = r + 1 {
// 		h := k / l
// 		if h > 0 {
// 			r = min(k/h, n)
// 		} else {
// 			r = n
// 		}
// 		w := r - l + 1
// 		s := (l + r) * w / 2
// 		sum -= h * s
// 	}
// 	return sum
// }

// // 二维整除分块.
// // ∑{i=1..min(n,m)} floor(n/i)*floor(m/i)
// // https://www.luogu.com.cn/blog/command-block/zheng-chu-fen-kuai-ru-men-xiao-ji
// func FloorSum2D(n, m int) (sum int) {
// 	for l, r := 1, 0; l <= min(n, m); l = r + 1 {
// 		hn, hm := n/l, m/l
// 		r = min(n/hn, m/hm)
// 		w := r - l + 1
// 		sum += hn * hm * w
// 	}
// 	return
// }

/**
 * 类欧几里得算法.
 * 商求和(divSum).
 * `∑⌊(ai+b)/m⌋, i in [0,n-1]`.
 */
function floorSum(n: number, m: number, a: number, b: number): number {
  let res = 0
  if (a < 0) {
    const a2 = (a % m) + m
    res -= ((n * (n - 1)) / 2) * ((a2 - a) / m)
    a = a2
  }
  if (b < 0) {
    const b2 = (b % m) + m
    res -= n * ((b2 - b) / m)
    b = b2
  }

  while (true) {
    if (a >= m) {
      res += ((n * (n - 1)) / 2) * Math.floor(a / m)
      a %= m
    }
    if (b >= m) {
      res += n * Math.floor(b / m)
      b %= m
    }
    const yMax = a * n + b
    if (yMax < m) {
      break
    }
    n = Math.floor(yMax / m)
    b = yMax % m
    const tmp = m
    m = a
    a = tmp
  }

  return res
}

/**
 * 余数求和(ModSum/remainderSum).
 * ! `∑k%i (i in [1,n]), 即 sum(k%i for i in [1,n])`.
 *
 * = ∑k-(k/i)*i
 * = n*k-∑(k/i)*i
 * 对于 [l,r] 范围内的 i，k/i 不变，此时 ∑(k/i)*i = (k/i)*∑i = (k/i)*(l+r)*(r-l+1)/2.
 */
function modSum(n: number, k: number): number {
  let sum = n * k
  for (let l = 1, r = 0; l <= n; l = r + 1) {
    const h = Math.floor(k / l)
    if (h > 0) {
      r = Math.min(Math.floor(k / h), n)
    } else {
      r = n
    }
    const w = r - l + 1
    const s = Math.floor(((l + r) * w) / 2)
    sum -= h * s
  }
  return sum
}

/**
 * 二维整除分块和.
 * `∑{i=1..min(n,m)} floor(n/i)*floor(m/i)`.
 */
function floorSum2D(n: number, m: number): number {
  let sum = 0
  for (let l = 1, r = 0; l <= Math.min(n, m); l = r + 1) {
    const hn = Math.floor(n / l)
    const hm = Math.floor(m / l)
    r = Math.min(Math.floor(n / hn), Math.floor(m / hm))
    const w = r - l + 1
    sum += hn * hm * w
  }
  return sum
}

export { floorSum, modSum, floorSum2D }

if (require.main === module) {
  const navie = (n: number, m: number, a: number, b: number): number => {
    let res = 0
    for (let i = 0; i < n; ++i) {
      res += Math.floor((a * i + b) / m)
    }
    return res
  }

  const modSumNavie = (n: number, k: number): number => {
    let res = 0
    for (let i = 1; i <= n; ++i) {
      res += k % i
    }
    return res
  }

  const floorSum2D = (n: number, m: number): number => {
    let res = 0
    for (let i = 1; i <= n; ++i) {
      for (let j = 1; j <= m; ++j) {
        res += Math.floor(n / i) * Math.floor(m / j)
      }
    }
    return res
  }

  for (let i = 0; i < 100; ++i) {
    const n = Math.floor(Math.random() * 1000)
    const m = Math.floor(Math.random() * 1000)
    const a = -50 + Math.floor(Math.random() * 1000)
    const b = -50 + Math.floor(Math.random() * 1000)
    const res = floorSum(n, m, a, b)
    const ans = navie(n, m, a, b)
    console.assert(res === ans)
  }
  console.log('ok1')

  for (let i = 0; i < 100; ++i) {
    const n = Math.floor(Math.random() * 1000)
    const k = Math.floor(Math.random() * 1000)
    const res = modSum(n, k)
    const ans = modSumNavie(n, k)
    console.assert(res === ans)
  }
  console.log('ok2')

  for (let i = 0; i < 100; ++i) {
    const n = Math.floor(Math.random() * 1000)
    const m = Math.floor(Math.random() * 1000)
    const res = floorSum2D(n, m)
    const ans = floorSum2D(n, m)
    console.assert(res === ans)
  }
  console.log('ok3')
}

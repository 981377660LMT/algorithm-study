// 函数复合

type F = (x: number) => number

function compose(functions: F[]): F {
  return functions.reduceRight(
    (pre, cur) => x => cur(pre(x)),
    x => x
  )
}

/**
 * const fn = compose([x => x + 1, x => 2 * x])
 * fn(4) // 9
 */

export {}

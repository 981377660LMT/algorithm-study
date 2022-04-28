// 二分法它最基本的应用是求一个单调函数的零点。
// 三分法是二分法的变种，他最基本的用途是求单峰函数的极值点。
const EPSILON = 1e-5

declare function f(x: number): number

/**
 *
 * @param left
 * @param right
 * @description
 * 在[left,right]间寻找单峰函数极值点
 * 每次在中点附近取点，那么每次可以减少约二分之一的长度。
 */
function search(left: number, right: number) {
  while (right - left > EPSILON) {
    const mid = (left + right) >> 1
    const f1 = f(mid - EPSILON)
    const f2 = f(mid + EPSILON)
    if (f1 < f2) left = mid
    else right = mid
  }
}

export {}

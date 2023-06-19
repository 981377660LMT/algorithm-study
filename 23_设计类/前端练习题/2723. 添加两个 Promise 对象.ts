// 给定两个 promise 对象 promise1 和 promise2，
// 返回一个新的 promise。promise1 和 promise2 都会被解析为一个数字。
// 返回的 Promise 应该解析为这两个数字的和。
async function addTwoPromises(
  promise1: Promise<number>,
  promise2: Promise<number>
): Promise<number> {
  const [a, b] = await Promise.all([promise1, promise2])
  return a + b
}

/**
 * addTwoPromises(Promise.resolve(2), Promise.resolve(2))
 *   .then(console.log); // 4
 */

export {}

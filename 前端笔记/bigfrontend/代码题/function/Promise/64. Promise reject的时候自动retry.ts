/**
 * @param {() => Promise<any>} fetcher
 * @param {number} maximumRetryCount
 * @return {Promise<any>}
 * @description
 * 因为网络问题API可能会失败。通常情况下我们可以提示错误，然后让用户重试。
   另外一种方案是 遇到网络问题时自动重试。
   请实现一个fetchWithAutoRetry(fetcher, count)，当出错的时候会自动重试，直到最大的重试次数。
   该问题中，你不需要判断错误是否是网络错误，所有rejection都认为网络错误即可。
 */
function fetchWithAutoRetry(fetcher: () => Promise<any>, maximumRetryCount: number): Promise<any> {
  // your code here
  return fetcher().catch(e => {
    if (maximumRetryCount === 0) throw e
    else return fetchWithAutoRetry(fetcher, maximumRetryCount - 1)
  })
}

async function fetchWithAutoRetry2(
  fetcher: () => Promise<any>,
  maximumRetryCount: number
): Promise<any> {
  let tries = 0

  while (true) {
    try {
      return await fetcher()
    } catch (e: any) {
      if (++tries > maximumRetryCount) throw e
    }
  }
}

fetchWithAutoRetry2(
  () =>
    new Promise((resolve, reject) => {
      if (Math.random() > 0.8) console.log(666)
      else reject(2)
    }),
  2
)

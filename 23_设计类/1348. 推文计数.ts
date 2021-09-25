type TweetName = string
type Time = number
type Freq = 'minute' | 'hour' | 'day'

class TweetCounts {
  private tweetMap: Map<TweetName, Time[]>
  constructor() {
    this.tweetMap = new Map()
  }

  /**
   *
   * @param tweetName
   * @param time
   * 用户 tweetName 在 time（以 秒 为单位）时刻发布了一条推文。
   */
  recordTweet(tweetName: string, time: number): void {
    !this.tweetMap.has(tweetName) && this.tweetMap.set(tweetName, [])
    this.tweetMap.get(tweetName)!.push(time)
  }

  /**
   *
   * @param freq
   * @param tweetName
   * @param startTime
   * @param endTime
   * 返回从开始时间 startTime（以 秒 为单位）到结束时间 endTime（以 秒 为单位）内，
   * 每 分 minute，时 hour 或者 日 day （取决于 freq）内指定用户 tweetName 发布的推文总数。
   */
  getTweetCountsPerFrequency(
    freq: Freq,
    tweetName: string,
    startTime: number,
    endTime: number
  ): number[] {
    endTime++
    let length: number
    if (freq === 'minute') length = 60
    else if (freq === 'hour') length = 60 * 60
    else length = 24 * 60 * 60

    // 有点问题
    const res = Array(~~((endTime - startTime - 1) / (length + 1))).fill(0)
    for (const time of this.tweetMap.get(tweetName) || []) {
      if (time >= startTime && time < endTime) {
        res[~~((time - startTime) / length)]++
      }
    }

    return res
  }
}

export {}

/**
 * Your TweetCounts object will be instantiated and called as such:
 * var obj = new TweetCounts()
 * obj.recordTweet(tweetName,time)
 * var param_2 = obj.getTweetCountsPerFrequency(freq,tweetName,startTime,endTime)
 */

//  可以将每个用户的推文时间存储方式换成更有效的平衡二叉树。
//  与暴力法所使用的线性表相比，
//  平衡二叉树保证其中的元素使用二叉树有序排列，其时间复杂度与线性表的区别为：
// 对于插入操作，线性表的时间复杂度为 O(1)，平衡二叉树的时间复杂度为 O(logn)。
// 对于查询操作，线性表的时间复杂度为 O(n)，平衡二叉树的时间复杂度为 O(logn)。
// 即 Map<TweetName,TreeSet<Time>> 其中treeset有bisectLeft/Right查找第几个的api

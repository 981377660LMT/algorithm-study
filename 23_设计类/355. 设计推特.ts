import { Deque } from '../2_queue/Deque/ArrayDeque'
import { PriorityQueue } from '../2_queue/优先级队列'

function* count() {
  let count = 0
  while (true) {
    yield count++
  }
}

class Timer {
  private timer: Generator<number>

  constructor() {
    this.timer = count()
  }

  next(): number {
    return this.timer.next().value
  }
}

type FollowerId = number
type FolloweeId = number
type UserId = number
interface Tweet {
  id: number
  time: number
}

class Twitter {
  private followMap: Map<FollowerId, Set<FolloweeId>>
  private tweetMap: Map<UserId, Deque<Tweet>> // 这里其实最好使用deque
  private timer: Timer

  constructor() {
    this.followMap = new Map()
    this.tweetMap = new Map()
    this.timer = new Timer()
  }

  /**
   *
   * @param userId
   * @param tweetId
   * 创建一条新的推文
   */
  postTweet(userId: number, tweetId: number): void {
    !this.tweetMap.has(userId) && this.tweetMap.set(userId, new Deque(10))
    this.tweetMap.get(userId)!.unshift({ id: tweetId, time: this.timer.next() })
  }

  /**
   *
   * @param userId
   * 检索最近的十条推文。每个推文都必须是由此用户关注的人或者是用户自己发出的。
   * 推文必须按照时间顺序由最近的开始排序。
   * @summary
   * 使用pq合并 k 个有序链表
   */
  getNewsFeed(userId: number): number[] {
    const allTweets = [...(this.followMap.get(userId) || []), userId]
      .filter(id => this.tweetMap.has(id))
      .map(id => this.tweetMap.get(id)!)

    const priorityQueue = new PriorityQueue<Tweet>((a, b) => b.time - a.time)

    allTweets.forEach(tweets => tweets.forEach(tweet => priorityQueue.push(tweet)))

    priorityQueue.heapify()

    const res: number[] = []
    while (priorityQueue.length && res.length < 10) {
      res.push(priorityQueue.shift()!.id)
    }

    return res
  }

  /**
   *
   * @param followerId
   * @param followeeId
   * 关注一个用户
   */
  follow(followerId: number, followeeId: number): void {
    !this.followMap.has(followerId) && this.followMap.set(followerId, new Set())
    this.followMap.get(followerId)!.add(followeeId)
  }

  /**
   *
   * @param followerId
   * @param followeeId
   * 取消关注一个用户
   */
  unfollow(followerId: number, followeeId: number): void {
    if (!this.followMap.has(followerId)) return
    this.followMap.get(followerId)!.delete(followeeId)
  }
}

const twitter = new Twitter()

// 用户1发送了一条新推文 (用户id = 1, 推文id = 5).
twitter.postTweet(1, 5)
twitter.postTweet(1, 3)
twitter.postTweet(1, 101)
twitter.postTweet(1, 13)
twitter.postTweet(1, 10)
twitter.postTweet(1, 2)
twitter.postTweet(1, 94)
twitter.postTweet(1, 505)
twitter.postTweet(1, 333)

// 用户1的获取推文应当返回一个列表，其中包含一个id为5的推文.
console.log(twitter.getNewsFeed(1))

// 用户1关注了用户2.
// console.log(twitter.follow(1, 2))

// // 用户2发送了一个新推文 (推文id = 6).
// console.log(twitter.postTweet(2, 6))

// // // 用户1的获取推文应当返回一个列表，其中包含两个推文，id分别为 -> [6, 5].
// // // 推文id6应当在推文id5之前，因为它是在5之后发送的.
// console.log(twitter.getNewsFeed(1))
// // // 用户1取消关注了用户2.
// twitter.unfollow(1, 2)

// // // 用户1的获取推文应当返回一个列表，其中包含一个id为5的推文.
// // // 因为用户1已经不再关注用户2.
// console.log(twitter.getNewsFeed(1))

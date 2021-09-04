class Twitter {
  constructor() {}

  postTweet(userId: number, tweetId: number): void {}

  getNewsFeed(userId: number): number[] {}

  follow(followerId: number, followeeId: number): void {}

  unfollow(followerId: number, followeeId: number): void {}
}

/**
 * Your Twitter object will be instantiated and called as such:
 * var obj = new Twitter()
 * obj.postTweet(userId,tweetId)
 * var param_2 = obj.getNewsFeed(userId)
 * obj.follow(followerId,followeeId)
 * obj.unfollow(followerId,followeeId)
 */

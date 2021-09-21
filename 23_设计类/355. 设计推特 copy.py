import collections
import heapq
import itertools


class Twitter(object):
    def __init__(self):
        self.timer = itertools.count(step=-1)
        self.tweets = collections.defaultdict(collections.deque)
        self.followees = collections.defaultdict(set)

    def postTweet(self, userId, tweetId):
        self.tweets[userId].appendleft((next(self.timer), tweetId))

    def getNewsFeed(self, userId):
        # 注意| {userId} 需要并上自己的推文
        tweets = heapq.merge(*(self.tweets[u] for u in self.followees[userId] | {userId}))
        # 至多取前10项
        return [t for _, t in itertools.islice(tweets, 10)]

    def follow(self, followerId, followeeId):
        self.followees[followerId].add(followeeId)

    def unfollow(self, followerId, followeeId):
        self.followees[followerId].discard(followeeId)


test = [(-1, 1), (0, 1)]
print([(a, b) for a, b in itertools.islice(test, 3)])


"""
推特的两个主要业务是：

1.发布推文(postTweet) 写
  用户可以向其粉丝发布新消息（平均 4.6k 请求 / 秒，峰值超过 12k 请求 / 秒）。

2.主页时间线(getNewsFeed) 读
  用户可以查阅他们关注的人发布的推文(300k 请求 / 秒）。

处理每秒 12,000 次写入（发推文的速率峰值）还是很简单的。
推特的伸缩性挑战并不是主要来自推特量，而是来自扇出(fan-out,为了服务一个传入请求而需要执行其他服务的请求数量。)

这一对操作有两种实现方式。
1. 读时计算量大:
发布推文时直接将推文存入用户推文的哈希表，当一个用户请求自己的主页时间线时，
首先查找他关注的所有人，查询这些被关注用户发布的推文并按时间顺序合并

2. 写时计算量大:
为每个用户的主页时间线维护一个缓存，就像每个用户的推文收件箱
当一个用户发布推文时，查找所有关注该用户的人，并将新的推文插入到每个人的推文收件箱
因此读取主页时间线的请求开销很小，因为结果已经提前计算好了。

ps:设计成deque便于删除队头已取消关注者的文章/删除队尾过期的文章(如果需要这个功能的话)

推特的第一个版本使用了方法 1,但系统很难跟上主页时间线查询的负载。所以公司转向了方法 2
方法 2 的效果更好，因为发推频率比查询主页时间线的频率几乎低了两个数量级
然而方法 2 的缺点是 一些用户有超过 3000 万的粉丝，这意味着一条推文就可能会导致主页时间线缓存的 3000 万次写入
在推特的例子中，每个用户粉丝数的分布（可能按这些用户的发推频率来加权）是探讨可伸缩性的一个关键负载参数，因为它决定了扇出负载
推特轶事的最终转折：现在已经稳健地实现了方法 2,推特逐步转向了两种方法的混合
即:当用户读取主页时间线时，分别地获取出该用户所关注的每位名流的推文，再与用户的主页时间线缓存合并
"""

from typing import List
from collections import defaultdict, deque
from heapq import merge
from itertools import count, islice


class Twitter:
    def __init__(self):
        self.timer = count(step=-1)
        self.tweet = defaultdict(deque)
        self.following = defaultdict(set)

    def postTweet(self, userId: int, tweetId: int) -> None:
        self.tweet[userId].appendleft((next(self.timer), userId, tweetId))

    def getNewsFeed(self, userId: int) -> List[int]:
        """时间复杂度O(klog(n)) 这里k为10"""
        tweets = merge(
            *(self.tweet[u] for u in self.following[userId] | {userId})
        )  # 注意| {userId} 需要并上自己的推文
        return [t for *_, t in islice(tweets, 10)]  # !至多取前10项

    def follow(self, followerId: int, followeeId: int) -> None:
        self.following[followerId].add(followeeId)

    def unfollow(self, followerId: int, followeeId: int) -> None:
        self.following[followerId].discard(followeeId)


class _Twitter2:
    def __init__(self):
        self.timer = count(step=-1)
        self.tweet = defaultdict(deque)
        self.feedLine = defaultdict(deque)  # 每个用户的推文收件箱
        self.following = defaultdict(set)  # 自己关注了谁
        self.follower = defaultdict(set)  # 谁关注了自己

    def postTweet(self, userId: int, tweetId: int) -> None:
        """当一个用户发布推文时，查找所有关注该用户的人，并将新的推文插入到每个人的推文收件箱"""
        newTweet = (next(self.timer), userId, tweetId)
        self.tweet[userId].appendleft(newTweet)
        for user in self.follower[userId] | {userId}:  # 注意需要并上自己
            self.feedLine[user].appendleft(newTweet)

    def getNewsFeed(self, userId: int) -> List[int]:
        """每个用户的主页时间线维护一个缓存，就像每个用户的推文收件箱;如果用户已经取关,那么就删除"""
        line, res = self.feedLine[userId], []
        while line and len(res) < 10:  # !至多取前10项
            head = line.popleft()
            postUser = head[1]
            if postUser in self.following[userId] | {userId}:  # 注意需要并上自己
                res.append(head)
        line.extendleft(res[::-1])
        return [t for *_, t in res]

    def follow(self, followerId: int, followeeId: int) -> None:
        """注意这种设计 follow时followeeId之前的动态不会出现在自己的line里"""
        self.following[followerId].add(followeeId)
        self.follower[followeeId].add(followerId)

    def unfollow(self, followerId: int, followeeId: int) -> None:
        """unfollow后在line里延迟删除"""
        self.following[followerId].discard(followeeId)
        self.follower[followeeId].discard(followerId)


if __name__ == "__main__":
    test = [(-1, 1), (0, 1)]
    print([(a, b) for a, b in islice(test, 3)])
    pq = [12, 3, 66]
    print(*merge(pq, [2, 3, 4]))

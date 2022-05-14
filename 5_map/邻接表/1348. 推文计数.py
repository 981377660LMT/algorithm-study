from typing import List
from collections import defaultdict
from sortedcontainers import SortedList


class TweetCounts:
    """一家社交媒体公司正试图通过分析特定时间段内出现的推文数量来监控其网站上的活动。"""

    def __init__(self):
        self.tweets = defaultdict(SortedList)

    def recordTweet(self, tweetName: str, time: int) -> None:
        """记录推文发布情况：用户 tweetName 在 time（以 秒 为单位）时刻发布了一条推文。"""
        self.tweets[tweetName].add(time)

    def getTweetCountsPerFrequency(
        self, freq: str, tweetName: str, startTime: int, endTime: int
    ) -> List[int]:
        """每个 时间块 中带有 tweetName 的 tweet 的数量"""
        if freq == "minute":
            seconds = 60
        elif freq == "hour":
            seconds = 60 * 60
        else:
            seconds = 60 * 60 * 24

        res = [0] * ((endTime - startTime) // seconds + 1)
        sl = self.tweets[tweetName]
        left, right = sl.bisect_left(startTime), sl.bisect_right(endTime) - 1
        for pos in range(left, right + 1):
            res[(sl[pos] - startTime) // seconds] += 1
        return res


if __name__ == '__main__':
    tweetCounts = TweetCounts()
    tweetCounts.recordTweet("tweet3", 0)
    tweetCounts.recordTweet("tweet3", 60)
    tweetCounts.recordTweet("tweet3", 10)
    tweetCounts.recordTweet("tweet3", 120)
    tweetCounts.recordTweet("tweet3", 160)
    print(tweetCounts.getTweetCountsPerFrequency("minute", "tweet3", 0, 210))

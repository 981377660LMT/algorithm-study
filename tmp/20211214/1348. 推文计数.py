from typing import List
from collections import defaultdict


# linear scan
class TweetCounts:
    def __init__(self):
        self.tweets = defaultdict(list)

    # 记录推文发布情况：用户 tweetName 在 time（以 秒 为单位）时刻发布了一条推文。
    def recordTweet(self, tweetName: str, time: int) -> None:
        self.tweets[tweetName].append(time)

    def getTweetCountsPerFrequency(
        self, freq: str, tweetName: str, startTime: int, endTime: int
    ) -> List[int]:
        if freq == "minute":
            seconds = 60
        elif freq == "hour":
            seconds = 3600
        else:
            seconds = 86400

        ans = [0] * ((endTime - startTime) // seconds + 1)
        for t in self.tweets[tweetName]:
            if startTime <= t <= endTime:
                ans[(t - startTime) // seconds] += 1
        return ans


from sortedcontainers import SortedList
from collections import defaultdict


class Leaderboard:
    def __init__(self):
        self.record = defaultdict()
        self.scores = SortedList(key=lambda x: -x)

    def addScore(self, playerId: int, score: int) -> None:
        if playerId not in self.record:
            self.scores.add(score)
            self.record[playerId] = score
        else:
            pre_score = self.record[playerId]
            self.scores.discard(pre_score)
            self.scores.add(pre_score + score)
            self.record[playerId] += score

    def top(self, K: int) -> int:
        return sum(self.scores.islice(0, K))

    def reset(self, playerId: int) -> None:
        score = self.record[playerId]
        self.record.pop(playerId)
        self.scores.discard(score)

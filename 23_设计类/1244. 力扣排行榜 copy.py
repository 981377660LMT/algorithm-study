from sortedcontainers import SortedList
from collections import defaultdict


class Leaderboard:
    def __init__(self):
        self.record = defaultdict()
        self.scores = SortedList(key=lambda x: -x[0])

    def addScore(self, playerId: int, score: int) -> None:
        if playerId not in self.record:
            self.scores.add((score, playerId))
            self.record[playerId] = score
        else:
            pre_score = self.record[playerId]
            self.scores.discard((pre_score, playerId))
            self.scores.add((pre_score + score, playerId))
            self.record[playerId] += score

    def top(self, K: int) -> int:
        return sum(score for score, _ in self.scores.islice(0, K))

    def reset(self, playerId: int) -> None:
        score = self.record[playerId]
        self.record.pop(playerId)
        self.scores.discard((score, playerId))


lb = Leaderboard()
lb.addScore(1, 1)
lb.addScore(2, 2)
print(lb.top(2))
lb.reset(1)
print(lb.top(2))


from typing import List
from collections import Counter



class Solution:
    def rankTeams(self, votes: List[str]) -> str:
        counter = Counter()
        for word in votes:
            counter += Counter(word)

        return ''.join([team for team, _ in counter.most_common()])


print(Solution().rankTeams(votes=["ABC", "ACB", "ABC", "ACB", "ACB"]))
# 输出："ACB"
# 解释：A 队获得五票「排位第一」，没有其他队获得「排位第一」，所以 A 队排名第一。
# B 队获得两票「排位第二」，三票「排位第三」。
# C 队获得三票「排位第二」，两票「排位第三」。
# 由于 C 队「排位第二」的票数较多，所以 C 队排第二，B 队排第三。

from collections import defaultdict, deque
from typing import List

INF = int(1e20)


def genNexts(word: str):
    yield from (word[:i] + '*' + word[i + 1 :] for i in range(len(word)))


class Solution:
    def findLadders(self, beginWord: str, endWord: str, wordList: List[str]) -> List[List[str]]:
        dist = defaultdict(lambda: INF)
        dist[beginWord] = 0
        adjMap = defaultdict(set)
        res = []

        for word in wordList:
            for mid in genNexts(word):
                adjMap[mid].add(word)

        # queue直接存路径
        queue = deque([[beginWord]])

        while queue:
            path = queue.popleft()
            cur = path[-1]
            if cur == endWord:
                res.append(path[:])
            else:
                for mid in genNexts(cur):
                    for next in adjMap[mid]:
                        # 因为要找所有的路，所以用等于
                        if dist[cur] + 1 <= dist[next]:
                            queue.append(path + [next])
                            dist[next] = dist[cur] + 1
        return res


print(
    Solution().findLadders(
        beginWord="hit", endWord="cog", wordList=["hot", "dot", "dog", "lot", "log", "cog"]
    )
)


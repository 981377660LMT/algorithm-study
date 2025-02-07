# 3435. 最短公共超序列的字母出现频率
# https://leetcode.cn/problems/frequencies-of-shortest-supersequences/description/
#
# !给你一个字符串数组 words 。请你找到 words 所有 最短公共超序列 ，且确保它们互相之间无法通过排列得到。
# 最短公共超序列 指的是一个字符串，words 中所有字符串都是它的子序列，且它的长度 最短 。
# 请你返回一个二维整数数组 freqs ，表示所有的最短公共超序列，其中 freqs[i] 是一个长度为 26 的数组，
# 它依次表示一个最短公共超序列的所有小写英文字母的出现频率。你可以以任意顺序返回这个频率数组。
# 1 <= words.length <= 256
# words[i].length == 2
# words 中所有字符串由不超过 16 个互不相同的小写英文字母组成。
# words 中的字符串互不相同。
#
# 建图，把字符作为点，如果有序列c1c2，那么c1->c2连一条边.
# !本题等价于MinimumDirectedFeedbackVertexSet，即在有向图中删除最少数量的点，使得剩下的图没有环
# 每个字母在公共超序列中的出现次数要么是 1，要么是 2。
# 枚举出现 2 次的字母（枚举子集），那么其余字母只出现 1 次。
# 如果出现 1 次的字母的约束无环，则说明我们当前枚举的子集是合法的，加入答案候选项中。
# 答案候选项中的大小最小的那些集合，即为最终答案。

from collections import deque
from typing import List

INF = int(1e18)


def enumerateSubset(s: int, allowEmpty=True):
    """降序枚举子集(包含s自身)."""
    t = s
    while t:
        yield t
        t = (t - 1) & s
    if allowEmpty:
        yield 0


class Solution:
    def supersequences(self, words: List[str]) -> List[List[int]]:
        allMask = 0
        graph = [[] for _ in range(26)]
        for w in words:
            u, v = ord(w[0]) - ord("a"), ord(w[1]) - ord("a")
            allMask |= 1 << u | 1 << v
            graph[u].append(v)

        def hasCycle(mask: int) -> bool:
            """判断mask对应的子集是否有环."""
            indeg = [0] * 26
            for i in range(26):
                if mask >> i & 1:
                    for j in graph[i]:
                        if mask >> j & 1:
                            indeg[j] += 1
            queue = deque()
            for i in range(26):
                if mask >> i & 1 and indeg[i] == 0:
                    queue.append(i)
            while queue:
                u = queue.popleft()
                for v in graph[u]:
                    if mask >> v & 1:
                        indeg[v] -= 1
                        if indeg[v] == 0:
                            queue.append(v)
            return any(indeg[i] > 0 for i in range(26) if mask >> i & 1)

        resMasks = []
        minSize = INF
        for m in enumerateSubset(allMask, allowEmpty=True):  # 枚举出现 2 次的字母
            curSize = m.bit_count()
            if curSize <= minSize and not hasCycle(~m):  # 出现 1 次的字母不能有环
                if curSize < minSize:
                    minSize = curSize
                    resMasks.clear()
                resMasks.append(m)

        res = []
        for m in resMasks:
            freq = [(m >> i & 1) + (allMask >> i & 1) for i in range(26)]
            res.append(freq)
        return res


if __name__ == "__main__":
    words = ["aa", "bb", "cc"]
    print(Solution().supersequences(words))

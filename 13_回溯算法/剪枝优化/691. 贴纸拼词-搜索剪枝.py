# param stickers  stickers 长度范围是 [1, 50]
# param target  target 的长度在 [1, 15] 范围内，由小写字母组成
# 拼出目标 target 所需的最小贴纸数量是多少？如果任务不可能，则返回 -1。

from typing import List
from collections import Counter, defaultdict, deque
from functools import lru_cache


class Solution:
    def minStickers(self, stickers: List[str], target: str) -> int:
        """记忆化dfs，不用状压，状压的写法反而更慢"""

        @lru_cache(None)
        def dfs(cur: str) -> int:
            if not cur:
                return 0

            res = int(1e20)
            for select in counters:
                # 排列剪枝成组合
                if cur[0] not in select:
                    continue
                next = replace(select, cur)
                res = min(res, dfs(next) + 1)
            return res

        # 耗尽Counter删除word里的字符
        def replace(counter: Counter, word: str) -> str:
            for char in counter:
                word = word.replace(char, '', counter[char])
            return word

        counters = [Counter(s) for s in stickers]
        res = dfs(target)
        dfs.cache_clear()
        return res if res != int(1e20) else -1

    def minStickers2(self, stickers: List[str], target: str) -> int:
        """bfs求无权图最短路"""
        # 耗尽Counter删除word里的字符
        def replace(counter: Counter, word: str) -> str:
            for char in counter:
                word = word.replace(char, '', counter[char])
            return word

        counters = [Counter(s) for s in stickers]
        dist = defaultdict(lambda: int(1e20))
        queue = deque([(target, 0)])
        while queue:
            cur, curDist = queue.popleft()
            if cur == '':
                return curDist
            for select in counters:
                # 排列剪枝成组合
                if cur[0] not in select:
                    continue
                next = replace(select, cur)
                if dist[next] > curDist + 1:
                    dist[next] = curDist + 1
                    queue.append((next, dist[next]))

        return -1


print(Solution().minStickers(["with", "example", "science"], "thehat"))
print(Solution().minStickers2(["with", "example", "science"], "thehat"))
# 输出：
# 3

# 我们可以使用 2 个 "with" 贴纸，和 1 个 "example" 贴纸。
# 把贴纸上的字母剪下来并重新排列后，就可以形成目标 “thehat“ 了。
# 此外，这是形成目标字符串所需的最小贴纸数量。

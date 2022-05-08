# 你的音乐播放器里有 N 首不同的歌，在旅途中，你的旅伴想要听 L 首歌（不一定不同，即，允许歌曲重复）。请你为她按如下规则创建一个播放列表：
# 每首歌至少播放一次。
# 一首歌只有在其他 K 首歌播放完之后才能再次播放。

# 0 <= K < N <= L <= 100
from functools import lru_cache

MOD = int(1e9 + 7)


class Solution:
    def numMusicPlaylists(self, N: int, L: int, K: int) -> int:
        """一共N首歌曲，目标列表长度是L，相同歌曲的最小间隔是k。
        
        求合法的歌单数
        dfs(index,remain)
        """

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            """当前在index,再选remain类歌"""
            if remain < 0:
                return 0
            if index == L:
                return 1 if remain == 0 else 0

            # 选新歌
            res1 = dfs(index + 1, remain - 1) * remain % MOD

            # 选旧歌
            res2 = dfs(index + 1, remain) * max(0, N - remain - K) % MOD

            return (res1 + res2) % MOD

        res = dfs(0, N)
        dfs.cache_clear()
        return res


print(Solution().numMusicPlaylists(3, 3, 1))
# 输出：6
# 解释：有 6 种可能的播放列表。[1, 2, 3]，[1, 3, 2]，[2, 1, 3]，[2, 3, 1]，[3, 1, 2]，[3, 2, 1].


# 你的音乐播放器里有 N 首不同的歌，在旅途中，你的旅伴想要听 L 首歌（不一定不同，即，允许歌曲重复）。请你为她按如下规则创建一个播放列表：
# 每首歌至少播放一次。
# 一首歌只有在其他 K 首歌播放完之后才能再次播放。

# 第二类斯特林数
class Solution:
    def numMusicPlaylists(self, n: int, goal: int, k: int) -> int:
        MOD = int(1e9 + 7)
        # i长度歌单，j首歌
        dp = [[0] * (n + 1) for _ in range(goal + 1)]
        dp[0][0] = 1

        for i in range(1, goal + 1):
            for j in range(1, n + 1):
                # 分成两种情况
                # 如果当前的歌和`前面的都不一样`，歌单前i-1首歌只包括了j-1首不同的歌曲，
                # 那么当前的选择有dp[i-1][j-1] * (N-j+1)
                # 如果当前的歌和`前面的有重复的`，那最近的K首必然是不能重复的，
                # 所以选择就是dp[i-1][j] * max(0, j-K)
                dp[i][j] = (dp[i - 1][j - 1] * (n - (j - 1)) + dp[i - 1][j] * max(0, j - k)) % MOD

        return dp[-1][-1]


print(Solution().numMusicPlaylists(n=3, goal=3, k=1))
# 输出：6
# 解释：有 6 种可能的播放列表。[1, 2, 3]，[1, 3, 2]，[2, 1, 3]，[2, 3, 1]，[3, 1, 2]，[3, 2, 1].


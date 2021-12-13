# 如果一个颜色片段为 'A' 且 相邻两个颜色 都是颜色 'A' ，那么 Alice `可以`删除该颜色片段。Alice 不可以 删除任何颜色 'B' 片段。
# 如果一个颜色片段为 'B' 且 相邻两个颜色 都是颜色 'B' ，那么 Bob `可以`删除该颜色片段。Bob 不可以 删除任何颜色 'A' 片段。
# Alice 和 Bob 不能 从字符串两端删除颜色片段。
# 如果其中一人无法继续操作，则该玩家 输 掉游戏且另一玩家 获胜 。

# 如果 Alice 获胜，请返回 true，否则 Bob 获胜，返回 false。

# 1 <= colors.length <= 10^5

# 总结：Count "AAA" and "BBB" and compare them
class Solution:
    def winnerOfGame(self, s: str) -> bool:
        a = b = 0

        for i in range(1, len(s) - 1):
            if s[i - 1] == s[i] == s[i + 1]:
                if s[i] == 'A':
                    a += 1
                else:
                    b += 1

        return a > b


print(Solution().winnerOfGame(s="AAABABB"))
# 输出：true
# 解释：
# AAABABB -> AABABB
# Alice 先操作。
# 她删除从左数第二个 'A' ，这也是唯一一个相邻颜色片段都是 'A' 的 'A' 。

# 现在轮到 Bob 操作。
# Bob 无法执行任何操作，因为没有相邻位置都是 'B' 的颜色片段 'B' 。
# 因此，Alice 获胜，返回 true 。


# Alice 和 Bob 玩一个游戏，两人轮流行动，Alice 先手 。
# 每次将问号替换成0-9之间的一个数
# Bob 获胜的条件是 num 中前一半数字的和 等于 后一半数字的和。Alice 获胜的条件是前一半的和与后一半的和 不相等 。

# 2 <= num.length <= 105
# num.length 是 偶数 。

# 1.
# If Bob wants to win,
# the number of '?' in input have to be even.
# 2.
# If we add the same sum or number of digits to left and right,
# this operation won't change the result.

# Bob 总可以保证它和 Alice 在相邻的两次操作中替换的数字之和为 9
class Solution:
    def sumGame(self, A: str) -> bool:
       n = len(num)
        
        def get(s: str) -> (int, int):
            nn = qq = 0
            for ch in s:
                if ch == "?":
                    qq += 1
                else:
                    nn += int(ch)
            return nn, qq
        
        n0, q0 = get(num[:n//2])
        n1, q1 = get(num[n//2:])
        
        return (q0 + q1) % 2 == 1 or n0 - n1 != (q1 - q0) * 9 // 2



print(Solution().sumGame(A="25??"))
# 输出：true
# 解释：Alice 可以将两个 '?' 中的一个替换为 '9' ，Bob 无论如何都无法使前一半的和等于后一半的和。

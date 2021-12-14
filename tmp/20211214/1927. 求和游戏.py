# 玩家可以执行以下操作：
# 选择一个下标 i 满足 num[i] == '?' 。
# 将 num[i] 用 '0' 到 '9' 之间的一个数字字符替代。
# Bob 获胜的条件是 num 中前一半数字的和 等于 后一半数字的和。
# Alice 获胜的条件是前一半的和与后一半的和 不相等 。

# num.length 是 偶数 。


# Bob能在对手干扰下赢的策略只有保证始终能凑出9，只有本身的差距是9的倍数而且两边操作数足够凑出这个差距。
class Solution:
    def sumGame(self, num: str) -> bool:
        """
        :type num: str
        :rtype: bool
        """
        s = a = b = 0
        n = len(num)
        for i, c in enumerate(num):
            if i < n // 2:
                if c == '?':
                    a += 1
                else:
                    s += int(c)
            else:
                if c == '?':
                    b += 1
                else:
                    s -= int(c)

        # Alice 先手， 总能使得和不为9的倍数，从而最终的和不相等
        if (a + b) % 2 == 1:
            return True

        # 两人有相同的次数，差为总能凑数9的倍数的话，Bob胜
        if s % 9 == 0 and s // 9 == (b - a) // 2:
            return False

        return True


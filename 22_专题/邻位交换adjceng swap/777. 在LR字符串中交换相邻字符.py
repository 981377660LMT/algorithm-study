# 一次移动操作指用一个"LX"替换一个"XL"，或者用一个"XR"替换一个"RX"。
# 题目的意思是说 ‘R’只能向右移动，并且只能移向’X’，‘L’只能向左移动，并且只能移向’X’。
# 就是求 两个字符串里, 相同字母的对应位置连线,是否相交, 如果不相交 且满足 <-L, R-> 就可以交换得到.
# """
# R X X L R X R X L
#  \   /   \  |  /
# X R L X X R R L X
# """

"""LR字符串"""

# !1.一个串能被拼成的必要条件是两者去除 - 后的字符串是一样的
# !2.原串的L不能在目标串的左边，R不能在目标串的右边


class Solution:
    def canTransform(self, start: str, end: str) -> bool:
        S1 = [(v, i) for i, v in enumerate(start) if v != "X"]
        S2 = [(v, i) for i, v in enumerate(end) if v != "X"]
        if len(S1) != len(S2):
            return False
        for c1, c2 in zip(S1, S2):
            if c1[0] != c2[0]:
                return False
            if c1[0] == "L":
                if c1[1] < c2[1]:
                    return False
            if c1[0] == "R":
                if c1[1] > c2[1]:
                    return False

        return True


print(Solution().canTransform("RXXLRXRXL", "XRLXXRRLX"))

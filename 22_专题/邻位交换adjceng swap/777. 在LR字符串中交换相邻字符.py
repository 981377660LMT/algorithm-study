# 一次移动操作指用一个"LX"替换一个"XL"，或者用一个"XR"替换一个"RX"。
# 题目的意思是说 ‘R’只能向右移动，并且只能移向’X’，‘L’只能向左移动，并且只能移向’X’。
# 就是求 两个字符串里, 相同字母的对应位置连线,是否相交, 如果不相交 且满足 <-L, R-> 就可以交换得到.
# """
# R X X L R X R X L
#  \   /   \  |  /
# X R L X X R R L X
# """
class Solution:
    def canTransform(self, start: str, end: str) -> bool:
        S = [(v, i) for i, v in enumerate(start) if v != 'X']
        E = [(v, i) for i, v in enumerate(end) if v != 'X']
        if len(S) != len(E):
            return False
        for s, e in zip(S, E):
            if s[0] != e[0]:
                return False
            if s[0] == 'L':
                if s[1] < e[1]:
                    return False
            if s[0] == 'R':
                if s[1] > e[1]:
                    return False

        return True


print(Solution().canTransform("RXXLRXRXL", "XRLXXRRLX"))


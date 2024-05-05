# 0 <= A <= 100
# 0 <= B <= 100

# 贪心：如果a>b or b>a 则连续输出2个a或2个b然后接另外一个字符，直至a==b
# 然后交替输出a,b
# 对于给定的 A 和 B，保证存在满足要求的 S。
class Solution:
    def strWithout3a3b(self, A: int, B: int) -> str:
        if A == 0:
            return 'b' * B
        if B == 0:
            return 'a' * A
        if A == B:
            return 'ab' * A
        if A > B:
            return 'aab' + self.strWithout3a3b(A - 2, B - 1)
        if A < B:
            return 'bba' + self.strWithout3a3b(A - 1, B - 2)


print(Solution().strWithout3a3b(A=1, B=2))
# 输出："abb"
# 解释："abb", "bab" 和 "bba" 都是正确答案。

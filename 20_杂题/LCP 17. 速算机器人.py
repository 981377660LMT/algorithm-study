class Solution:
    def calculate(self, s: str) -> int:
        # x, y = 1, 0
        # for char in s:
        #     if char == 'A':
        #         x = 2 * x + y
        #     elif char == 'B':
        #         y = 2 * y + x
        # return x + y
        return 1 << len(s)


# "A" 运算：使 x = 2 * x + y；
# "B" 运算：使 y = 2 * y + x。
# 请返回最终 x 与 y 的和为多少。
# 目标结果是x+y 出现一个"A"，有x+y=(2x+y)+y=2x+2y 出现一个"B"，有x+y=x+(2y+x)=2x+2y 所以每出现一个A/B，都使x+y的值翻倍 因此结果是2**len(s)

from collections import deque

# 保证 n 可以写成 2^k 的形式，其中 k 是正整数

# 1. 头尾匹配 nlogn
class Solution:
    def findContestMatch2(self, n: int) -> str:
        res = list(map(str, range(1, n + 1)))
        while len(res) > 1:
            next = []
            for i in range(len(res) // 2):
                # i 与 n-1-i配对
                next.append(f'({res[i]},{res[-i-1]})')
            res = next
        return res[0]

    def findContestMatch(self, n: int) -> str:
        res = deque(map(str, range(1, n + 1)))
        while len(res) > 1:
            next = deque()
            while res:
                x, y = res.popleft(), res.pop()
                next.append(f'({x},{y})')
            res = next
        return res[0]


print(Solution().findContestMatch(2))
print(Solution().findContestMatch(4))
# 输出: ((1,4),(2,3))
print(Solution().findContestMatch(8))

# 输入: 8
# 输出: (((1,8),(4,5)),((2,7),(3,6)))
# 解析:
# 第一轮: (1,8),(2,7),(3,6),(4,5)
# 第二轮: ((1,8),(4,5)),((2,7),(3,6))
# 第三轮 (((1,8),(4,5)),((2,7),(3,6)))
# 由于第三轮会决出最终胜者，故输出答案为(((1,8),(4,5)),((2,7),(3,6)))。

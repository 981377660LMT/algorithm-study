import random
from typing import List


class Solution:
    def __init__(self, n: int, blacklist: List[int]):
        """
        预处理思路：
        - 可选整数总共有 M = n - len(blacklist) 个，范围 [0, M-1]。
        - 黑名单中的小于 M 的数会“占用”这个区间，需要映射到尾部的合法数上。
        - 尾部区间是 [M, n-1]，其中剔除黑名单后恰好有 len(blacklist) 个空位。
        - 逐个把黑名单中 < M 的 b 映射到尾部的一个空位 w。
        """
        self.M = n - len(blacklist)
        self.mapping = dict()
        black = set(blacklist)
        w = n - 1
        # 对每个 b < M 的黑名单数，找一个尾部 w 不在黑名单中映射给它
        for b in blacklist:
            if b < self.M:
                # 从尾部向前跳过所有在黑名单里的数
                while w in black:
                    w -= 1
                self.mapping[b] = w
                w -= 1

    def pick(self) -> int:
        """
        1. 在 [0, M-1] 里均匀随机选一个 x；
        2. 如果 x 在映射表里，说明它原本是黑名单数，用映射值返回；
           否则 x 本身就是合法数，直接返回。
        仅调用一次随机函数，且时间 O(1)。
        """
        x = random.randint(0, self.M - 1)
        return self.mapping.get(x, x)


if __name__ == "__main__":
    n = 10
    blacklist = [1, 3, 5, 8]
    sol = Solution(n, blacklist)
    # 统计多次 pick 的分布，验证不落在 blacklist，也近似均匀
    from collections import Counter

    cnt = Counter(sol.pick() for _ in range(100000))
    print("出现次数（示例）:", cnt)
    # 确认没有出现黑名单中的值
    print("有无黑名单出现？", any(b in cnt for b in blacklist))

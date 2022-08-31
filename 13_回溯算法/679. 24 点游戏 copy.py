from itertools import combinations
from typing import List


EPS = 1e-5


class Solution:
    def judgePoint24(self, cards: List[int]) -> bool:
        """每次选出两个数进行做运算"""

        def genNext(a: int, b: int):
            """生成器简化逻辑"""
            yield a + b
            yield a - b
            yield b - a
            yield a * b
            if b != 0:
                yield a / b
            if a != 0:
                yield b / a

        if len(cards) == 1:
            return abs(cards[0] - 24) < EPS
        for i, j in combinations(range(len(cards)), 2):
            num1, num2 = cards[i], cards[j]
            for num3 in genNext(num1, num2):
                newCards = [num3] + [num for k, num in enumerate(cards) if k not in (i, j)]
                if self.judgePoint24(newCards):  # type: ignore
                    return True
        return False


if __name__ == "__main__":
    print(Solution().judgePoint24([1, 2, 3, 4]))

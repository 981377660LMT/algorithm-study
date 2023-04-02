from BIT import BIT2


class Solution:
    def canReach(self, s: str, minJump: int, maxJump: int) -> bool:
        """这个树状数组范围更新可以用差分数组优化，因为query和更新同时进行"""
        if s[-1] == "1":
            return False
        n = len(s)
        bit = BIT2(n + 10)
        for i, char in enumerate(s):
            if char == "0":
                if i == 0 or bit.query(i, i + 1) != 0:
                    bit.add(i + minJump, i + maxJump + 1, 1)
        return bit.query(n - 1, n) != 0

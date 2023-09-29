# https://leetcode.cn/problems/the-wording-game/description/
# 2868. 单词游戏


from typing import List


class Solution:
    def canAliceWin(self, a: List[str], b: List[str]) -> bool:
        def alicePlay() -> bool:
            nonlocal alicePtr, pre
            while alicePtr < len(a):
                cur = a[alicePtr]
                alicePtr += 1
                if (cur[0] == pre[0] and cur > pre) or (ord(cur[0]) == ord(pre[0]) + 1):
                    pre = cur
                    return True
            return False

        def bobPlay() -> bool:
            nonlocal bobPtr, pre
            while bobPtr < len(b):
                cur = b[bobPtr]
                bobPtr += 1
                if (cur[0] == pre[0] and cur > pre) or (ord(cur[0]) == ord(pre[0]) + 1):
                    pre = cur
                    return True
            return False

        alicePtr = 1
        bobPtr = 0
        pre = a[0]
        while True:
            if not bobPlay():
                return True
            if not alicePlay():
                return False


if __name__ == "__main__":
    print(Solution().canAliceWin(["aa", "ba", "ca", "da", "ea"], ["ab", "bb", "db"]))

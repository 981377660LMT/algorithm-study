from typing import List


class Solution:
    def countBits(self, n: int) -> List[int]:
        """
        DP with bit manipulation:
        Let ans[i] = number of 1s in binary of i.
        Observe:
          ans[0] = 0
          For i > 0: ans[i] = ans[i >> 1] + (i & 1)
        This runs in O(n) time and O(n) space.
        """
        res = [0] * (n + 1)
        for i in range(1, n + 1):
            res[i] = res[i >> 1] + (i & 1)
        return res


if __name__ == "__main__":
    sol = Solution()
    print(sol.countBits(5))  # [0,1,1,2,1,2]
    print(sol.countBits(10))  # [0,1,1,2,1,2,2,3,1,2,2]

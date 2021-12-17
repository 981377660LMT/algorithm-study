class Solution:
    def minimumPerimeter(self, neededApples: int) -> int:
        left, right = 1, int(1e100)
        while left <= right:
            mid = (left + right) // 2
            # 参考平方和公式
            if 2 * mid * (mid + 1) * (mid * 2 + 1) >= neededApples:
                right = mid - 1
            else:
                left = mid + 1
        return left * 8


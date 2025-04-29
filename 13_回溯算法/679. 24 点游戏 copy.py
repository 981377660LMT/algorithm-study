from typing import List

EPS = 1e-6


class Solution:
    def judgePoint24(self, nums: List[int]) -> bool:
        """
        回溯搜索所有两两运算组合：
        1. 从当前数字列表中任选 i<j 两个数 a,b；
        2. 对 a,b 应用 +, -, b-a, *, a/b (若 b!=0), b/a (若 a!=0)；
        3. 把结果加入剩余数字列表，递归判断；
        4. 若列表只剩一个数且和 24 误差 < 1e-6，则成功。
        时间 O(1)（常数规模 4!×操作数目），空间 O(1)。
        """

        def dfs(arr: List[float]) -> bool:
            if len(arr) == 1:
                return abs(arr[0] - 24) < EPS
            n = len(arr)

            for i in range(n):
                a = arr[i]
                for j in range(i + 1, n):
                    b = arr[j]

                    rest = [arr[k] for k in range(n) if k != i and k != j]

                    for val in (a + b, a - b, b - a, a * b):
                        if dfs(rest + [val]):
                            return True
                    # 除法要防零
                    if abs(b) > EPS and dfs(rest + [a / b]):
                        return True
                    if abs(a) > EPS and dfs(rest + [b / a]):
                        return True
            return False

        return dfs(list(map(float, nums)))


if __name__ == "__main__":
    sol = Solution()
    tests = [
        ([4, 1, 8, 7], True),  # (8-4)*(7-1)=24
        ([1, 2, 1, 2], False),
        ([3, 3, 8, 8], True),  # 8/(3-8/3)=24
    ]
    for cards, expect in tests:
        res = sol.judgePoint24(cards)
        print(f"{cards} -> {res} (expected {expect})")

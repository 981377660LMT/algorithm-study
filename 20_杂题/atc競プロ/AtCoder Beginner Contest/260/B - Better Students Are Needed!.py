# !1.数学の点が高い方から X 人を合格とする。
# !2.次に、この時点でまだ合格となっていない受験者のうち、英語の点が高い方から Y 人を合格とする。
# !3.次に、この時点でまだ合格となっていない受験者のうち、数学と英語の合計点が高い方から Z 人を合格とする。
# ここまでで合格となっていない受験者は、不合格とする。
# !4.ただし、 1. から 3. までのどの段階についても、同点であった場合は受験生の番号の小さい方を優先します。

# n<=1e5

# 注意需要三次自定义排序
import sys
import os
from typing import List, Tuple


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    def assign(nums: List[Tuple[int, int, int]], visited: List[bool], remain: int) -> List[int]:
        res = []
        for *_, id in nums:
            if not remain:
                break
            if visited[id]:
                continue
            visited[id] = True
            remain -= 1
            res.append(id)
        return res

    n, x, y, z = map(int, input().split())
    sugaku = list(map(int, input().split()))
    eigo = list(map(int, input().split()))
    students = list((a, b, i) for i, (a, b) in enumerate(zip(sugaku, eigo), start=1))
    nums1 = sorted(students, key=lambda x: (-x[0], x[2]))
    nums2 = sorted(students, key=lambda x: (-x[1], x[2]))
    nums3 = sorted(students, key=lambda x: (-x[0] - x[1], x[2]))

    visited = [False] * (n + 10)

    res = assign(nums1, visited, x) + assign(nums2, visited, y) + assign(nums3, visited, z)
    print(*sorted(res), sep="\n")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

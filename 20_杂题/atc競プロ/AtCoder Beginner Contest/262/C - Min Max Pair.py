import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

# !注意comb函数不能在pypy3里使用


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    ok = [i for i, num in enumerate(nums, start=1) if num == i]
    res1 = len(ok) * (len(ok) - 1) // 2
    res2 = 0
    for i1, num1 in enumerate(nums, start=1):
        if 0 <= num1 - 1 < n and i1 != num1:
            num2 = nums[num1 - 1]
            if i1 == num2:
                res2 += 1

    print(res1 + res2 // 2)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

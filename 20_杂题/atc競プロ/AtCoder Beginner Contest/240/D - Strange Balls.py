# 每次往桶里加球
# 如果写着k的球有k个连续 那么就消除
# 求每次操作后桶里球的个数

# 类似于最大栈的思路、相邻元素消除问题
# !每次push多加一个连续长度的信息
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    stack = []
    size = 0

    for num in nums:
        # push
        if stack and stack[-1][0] == num:
            stack[-1][1] += 1
        else:
            stack.append([num, 1])
        size += 1

        # pop
        if stack[-1][0] == stack[-1][1]:
            size -= stack[-1][1]
            stack.pop()

        print(size)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

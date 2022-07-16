# 每个人有姓/名 给每个人起外号 外号是姓或者名
# 注意每个人的外号不能和其余人的姓名或者名字相同
# 问所有人是否能够起外号

# n<=1e5


# !哈希表计数，遍历每个人检查是否能起外号


from collections import defaultdict
import sys
import os


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    # n = int(input())
    # people = []
    # for _ in range(n):
    #     s, t = input().split()
    #     people.append((s, t))

    # for i in range(n):
    #     ok = False
    #     for cand in people[i]:
    #         flag = True
    #         for j in range(n):
    #             if i == j:
    #                 continue
    #             if cand == people[j][0] or cand == people[j][1]:
    #                 flag = False
    #                 break
    #         ok = ok or flag
    #     if not ok:
    #         print("No")
    #         exit(0)
    # print("Yes")

    n = int(input())
    people = []
    counter = defaultdict(int)
    for _ in range(n):
        s, t = input().split()
        people.append((s, t))
        counter[s] += 1
        counter[t] += 1

    for s, t in people:
        counter[s] -= 1
        counter[t] -= 1
        # 看是不是与别人名或姓相同
        if counter[s] >= 1 and counter[t] >= 1:
            print("No")
            exit(0)
        counter[s] += 1
        counter[t] += 1

    print("Yes")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

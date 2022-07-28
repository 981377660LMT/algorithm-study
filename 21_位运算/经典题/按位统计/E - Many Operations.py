"""按位计算+处理前缀的影响"""

# 每一轮有i次操作，每一次操作用一个ti和一个ai来改变x
# 如果 t= 1,  x=  x & a
# 如果 t= 2,  x = x | a
# 如果 t= 3,  x = x XOR a
# 现在进行n轮操作，
# !第i轮操作，将进行第1,2,3...i—1, i个opt
# 询问每一轮操作后的x值，x是一直更新的
# n<=2e5
# ai<=2^30

# !朴素的想法就是, 我们每一次操作完之后, 下一次操作无需从头扫描一遍, 而是记录下前缀的影响(0变为0/1, 1变为0/1)
# 当前操作只需要进行一次当前操作即可, 这样的话就可以变成 O(n)
# 整个操作我们没办法考虑, 那么想到操作都是位运算，根据位运算的特殊性质, 我们可以`按位统计`，记录每个位的变化


from operator import and_, or_, xor
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

CAL = {1: and_, 2: or_, 3: xor}


def main() -> None:
    n, x = map(int, input().split())
    opt = []
    for _ in range(n):
        t, a = map(int, input().split())
        opt.append((t, a))

    res = [0] * n
    for bit in range(30):
        pre = [0, 1]  # !前缀的影响 当前0 => 0 当前1 => 1
        cur = (x >> bit) & 1
        for i in range(n):
            t, a = opt[i]
            a = (a >> bit) & 1
            fromZero, fromOne = CAL[t](0, a), CAL[t](1, a)  # !当前位的影响
            pre = [fromZero if pre[0] == 0 else fromOne, fromZero if pre[1] == 0 else fromOne]
            cur = pre[cur]
            res[i] |= cur << bit
    print(*res, sep="\n")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

# 顺丰主管将对他手下的兼职快递员们进行任务安排。
# 已知快递分为两种类型，第一种类型为同城快送，第二种类型为跨城运输。顺丰主管手下有个兼职快递员，每个快递员的工资分别为，运送快递的顾客满意值为。快递员有两种类型：类型A只能送同城快运，类型B则可以送任意快递。
# 顺丰主管共有个同城快送任务，个跨城运输任务。他需要把这些快递派发给他手下的快递员，每个人最多只能接一个任务。顺丰主管想知道，在支付总工资尽可能小的情况下，最终的顾客满意值之和最大是多少？
# 注：任务必须全部接完。每个快递员只能接一个任务。 输入描述 第一行输入四个正整数，分别代表快递员数量和两种类型任务的数量。接下来的行，每行输入一个字符和两个正整数用来描述每个快递员的类型以及工资和顾客满意值。保证字符是'A'或者'B'。 输出描述 如果无法接完全部任务，则输出-1。否则一个整数，代表最终最大的顾客满意值。
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, x, y = map(int, input().split())
    manA, manB = [], []  # (salary, satisfaction)
    for _ in range(n):
        t, a, b = input().split()
        if t == "A":
            manA.append((int(a), int(b)))
        else:
            manB.append((int(a), int(b)))
    # 先送跨城, 再送同城
    if len(manB) < y:
        print(-1)
        exit(0)
    if len(manA) + len(manB) < x + y:
        print(-1)
        exit(0)
    manB.sort(key=lambda x: (x[0], -x[1]))
    pre = manB[:y]
    suf = manB[y:]
    allA = manA + suf
    allA.sort(key=lambda x: (x[0], -x[1]))

    res1 = sum([x[1] for x in pre])
    res2 = sum([x[1] for x in allA[:x]])
    print(res1 + res2)

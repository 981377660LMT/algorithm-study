import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 给定一个正整数，你能重排这个正整数的数位，使得它变成偶数（不含前导零）吗？
# 要求重排后的正整数和原数不能相等。
# 一共有  次询问。 输入描述 第一行输入一个正整数 ，代表询问次数。接下来的 行，每行输入一个正整数 。保证所有x的长度之和不超过100000 输出描述 输出行，每行代表一次询问。如果存在合法解，请输出一个重排后的正整数，务必保证其为偶数且和原数不等。有多解时输出任意即可。如果不存在合法解，直接输出-1。
if __name__ == "__main__":
    q = int(input())
    for _ in range(q):
        s = input()
        if len(s) == 1:
            print(-1)
            continue
        odd, even = [], []
        for i in s:
            if int(i) % 2 == 0:
                even.append(i)
            else:
                odd.append(i)
        if not even:
            print(-1)
            continue

        # 只有一种偶数的情形
        if len(set(even)) == 1:
            if not odd:
                print(-1)
                continue
            # 至少有一个奇数,且最后一个数为奇数
            if int(s[-1]) % 2 == 1:
                res = []
                todo = None
                for c in s:
                    if todo is None and int(c) % 2 == 0:
                        todo = c
                    else:
                        res.append(c)
                res.append(todo)
                print("".join(res))
                continue
            # 至少有一个奇数,且最后一个数为偶数

        else:
            # >=两个偶数的情形
            res = []
            todo = None
            for c in s:
                if todo is None and int(c) % 2 == 0 and c != s[-1]:
                    todo = c
                else:
                    res.append(c)
            res.append(todo)
            print("".join(res))

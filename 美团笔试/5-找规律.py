res = [0]
for i in range(10):
    res = [num + 1 for num in res] + res + [num + 2 for num in res]


# 。文件编号为1到n n<=1e10
while True:
    try:
        n, k = map(int, input().split())
        res = [0]
        for i in range(n):
            res = [num + 1 for num in res] + res + [num + 2 for num in res]

        print(res[k - 1])

    except EOFError:
        break


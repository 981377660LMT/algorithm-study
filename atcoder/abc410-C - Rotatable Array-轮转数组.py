# abc410-C - Rotatable Array-轮转数组.py
# 给定一个长度为 N 的数组 A，初始 A[i]=i（1 ≤ i ≤ N）。接下来处理 Q 次操作：
#
# 1 p x 把逻辑下标为 p 的元素赋值为 x
# 2 p  输出逻辑下标为 p 的元素
# 3 k  把数组整体向左循环移动 k 次（即把首元素依次移到末尾）。
# 对每条查询 ② 输出对应值。
#
# 维护shift

if __name__ == "__main__":
    N, Q = map(int, input().split())
    A = list(range(1, N + 1))

    lShift = 0
    for _ in range(Q):
        query = list(map(int, input().split()))
        op = query[0]
        if op == 1:
            p, x = query[1], query[2]
            p -= 1
            A[(p + lShift) % N] = x
        elif op == 2:
            p = query[1]
            p -= 1
            print(A[(p + lShift) % N])
        elif op == 3:
            k = query[1]
            lShift = (lShift + k) % N

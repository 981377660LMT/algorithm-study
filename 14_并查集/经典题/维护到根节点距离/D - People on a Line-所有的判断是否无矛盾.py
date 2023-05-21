# https://atcoder.jp/contests/abc087/tasks/arc090_b
# 每个判断为 right-left=dist
# 问所有的判断是否无矛盾


from UnionFindWithDist import UnionFindMapWithDist1

if __name__ == "__main__":
    n, m = map(int, input().split())
    uf = UnionFindMapWithDist1()
    for _ in range(m):
        left, right, weight = map(int, input().split())
        if not uf.union(left, right, weight):
            print("No")
            exit(0)

    print("Yes")

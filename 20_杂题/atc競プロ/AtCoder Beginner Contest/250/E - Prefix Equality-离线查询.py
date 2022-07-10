# 每个查询询问
# !nums1的前x项的集合是否与nums2的前y项的集合相等 (重复元素只算一次)
# n<=2e5
# numsi<=1e9
# x,y<=n


# !. 1.离线查询+双指针
# !固定x之后 就可以尺取寻找y的边界
# 2. 异或哈希


from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    q = int(input())
    Q = defaultdict(list)  # 按照查询左端点保存
    for i in range(q):
        x, y = map(int, input().split())  # 前 x/y 项
        x, y = x - 1, y - 1
        Q[x].append((y, i))

    res = [False] * q

    # !对每个nums1的前缀 nums2对应端点的合法的区间[left,right]唯一确定 用双指针确定
    left, right, cur = 0, 0, 0
    s1, s2 = set(), set()  # s2始终是s1的子集

    isOk = True
    for i in range(n):
        if nums1[i] not in s1:
            s1.add(nums1[i])
            # !可以加入的全加入
            while cur < n and len(s1) > len(s2) and nums2[cur] in s1:
                s2.add(nums2[cur])
                cur += 1

            if len(s1) == len(s2):  # 有解
                isOk = True
                left = right = cur - 1
                # !放入重复元素
                while cur < n and nums2[cur] in s1:
                    right += 1
                    cur += 1
            else:
                isOk = False

        for (qy, qi) in Q[i]:
            if isOk and left <= qy <= right:
                res[qi] = True

    for i in range(q):
        print("Yes" if res[i] else "No")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

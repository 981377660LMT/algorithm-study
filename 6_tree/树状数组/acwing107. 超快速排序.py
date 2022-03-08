# 通过交换两个相邻的序列元素来处理 n 个不同整数的序列，直到序列按升序排序。
# 对于输入序列 9 1 0 5 4，超快速排序生成输出 0 1 4 5 9。
# 您的任务是确定超快速排序需要执行多少交换操作才能对给定的输入序列进行排序。

# 求逆序对数即可:每次交换两个相邻的元素可以减少1个逆序对，最后逆序对为0
def mergesort(A):
    n = len(A)
    if n <= 1:
        return 0, A
    mid = n // 2
    cnt1, left = mergesort(A[:mid])
    cnt2, right = mergesort(A[mid:])
    i, j = 0, 0
    cnt = cnt1 + cnt2
    ans = []
    while i < len(left) and j < len(right):
        if left[i] < right[j]:
            ans.append(left[i])
            i += 1
        else:
            ans.append(right[j])
            j += 1
            cnt += len(left) - i
    if i < len(left):
        ans.extend(left[i:])
    else:
        ans.extend(right[j:])
    return cnt, ans


while True:
    n = int(input())
    if not n:
        break
    A = []
    for _ in range(n):
        A.append(int(input()))
    cnt, B = mergesort(A)
    print(cnt)


# 现在小团想要借书，请你帮忙看看小团能不能借到书，如果可以借到的话在哪一行书架上有这本书。
M, N, Q = map(int, input().split())
Shelf = [False] * (N + 1)
Books = [0] * (M + 1)
for _ in range(Q):
    nums = list(map(int, input().split()))
    opt = nums[0]
    if opt == 1:
        x, y = nums[1:]
        if Books[x] != -1 and not Shelf[y] and not Shelf[Books[x]]:
            Books[x] = y
    elif opt == 2:
        Shelf[nums[1]] = True
    elif opt == 3:
        Shelf[nums[1]] = False
    elif opt == 4:
        res, x = -1, nums[1]
        if Books[x] > 0 and not Shelf[Books[x]]:
            Books[x], res = -1, Books[x]
        print(res)
    elif opt == 5:
        x = nums[1]
        if Books[x] == -1:
            Books[x] = 0


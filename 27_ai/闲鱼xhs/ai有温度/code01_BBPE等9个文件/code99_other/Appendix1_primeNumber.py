def zhishu(num):
    ans = []
    for i in range(2, num):
        for j in range(2, i):
            if i%j == 0:
                break
        else:
            ans.append(i)
    return ans

ans = zhishu(100)




## 剪枝
import math

def zhishu(num):
    count = 0
    ans = []
    for i in range(2, num):
        for j in range(2, int(math.sqrt(i)) + 1):
            if i%j == 0:
                break
            count += 1
        else:
            ans.append(i)
    return ans, count

ans, count = zhishu(100)

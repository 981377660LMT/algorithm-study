from typing import List

# 请返回能够令 a == b 成立的三元组 (i, j , k) 的数目
# a = arr[i] ^ arr[i + 1] ^ ... ^ arr[j - 1]
# b = arr[j] ^ arr[j + 1] ^ ... ^ arr[k]

# 由a=b有a^a=0=a^b 则只需找前缀 prefix[k+1] == prefix[i] 中间j随便取
# 即在前缀数组找两个值相等的对数 O(n^2)
# 可以用哈希表记录优化 O(n)


class Solution:
    def countTriplets(self, arr: List[int]) -> int:
        preXor = [0]
        for num in arr:
            preXor.append(num ^ preXor[-1])

        res = 0
        for i in range(len(preXor)):
            for k in range(i, len(preXor) - 1):
                if preXor[i] == preXor[k + 1]:
                    res += k - i

        return res


print(Solution().countTriplets(arr=[2, 3, 1, 6, 7]))
# 输出：4
# 解释：满足题意的三元组分别是 (0,1,2), (0,2,2), (2,3,4) 以及 (2,4,4)

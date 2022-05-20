1. 所有元素到某个元素总距离之和最小
   排序找到这个元素 O(nlogn)
   `mid=sorted(nums[len(nums)//2])`
   快速选择找到这个元素 O(n)
   ...
   堆找到这个元素 O(nlog(n/2))
   `mid = nsmallest((len(nums) >> 1) + 1, nums)[-1]`

2. 计算中位数(奇数偶数长度一起)
   `mid=(nums[n>>1]+nums[(n-1)>>1])/2`

思路是排序后再求前缀和，最右二分求最后能和哪个数齐平

```Python
nums = sorted(nums)
preSum = [0] + list(accumulate(nums))
nums = [0] + nums
```

`k次加1操作最大化最小值`

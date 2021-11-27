1. 哈希表 lookup + 千纸鹤 accumulate
   `1477. 找两个和为目标值且不重叠的子数组.py`
   `1546. 和为目标值且不重叠的非空子数组的最大数目.py`

```Python
模板
lookup = {0: -1}
for i,running_sum in accumulate:
  if running_sum - target in lookup:
    ...
    lookup[running_sum]=i
```

2. 维护 record 数组，此数组记录每个位置之前的最值

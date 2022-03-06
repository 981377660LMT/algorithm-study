优雅的写法是 sorted(nums,key=func)
然后单独写这个 func 即可

```Python
class Solution:
    def sortJumbled(self, mapping: List[int], nums: List[int]) -> List[int]:
        def cmpFunc(num: int) -> int:
            return int(''.join(str(mapping[int(char)]) for char in str(num)))

        return sorted(nums, key=lambda n:)

```

如果比较麻烦，可以用 `from functools import cmp_to_key` 像 js 一样比较

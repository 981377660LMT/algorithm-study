1. 在有序集合中找出`最小的`大于等于 t 的元素
   最小 => 最左,bisect_left
   注意到这也符合 bisect_left 的 arr[index]>=t
   如果 index==len(arr) 下标越界 就是找不到

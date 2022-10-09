头尾指针模板

1. deque 模拟和双指针是一样的
   `left < right` 等价于 `len(queue) >= 2 (需要两个元素时)`
   `left <= right` 等价于 `len(queue) >= 1 (非空)`
   `left > right 等价于 len(queue) == 0 (空)`

   `s[left]` 等价于 `queue[0]`
   `s[right]` 等价于 `queue[-1]`

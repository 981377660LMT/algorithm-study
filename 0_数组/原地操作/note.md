双指针中的快慢指针。在这里快指针是读指针， 慢指针是写指针
只保留一个元素,删除所有相同元素：慢指针**先移再写**
移除所有数值等于 val 的元素:慢指针**先写再移**

先移再写

```JS
  //  双指针 没见过的就搬过来
  let slowP = 0
  for (let fastP = 0; fastP < nums.length; fastP++) {
    if (nums[fastP] !== nums[slowP]) {
      // 先移后写
      slowP++
      nums[slowP] = nums[fastP]
    }
  }

  // 原地移除 类似于链表的slow.next = null
  nums.length = slowP + 1
  return slowP + 1
```

先写再移

```JS
// 双指针 没见过的就搬过来 对应链表dummy节点做法
const removeElement = (nums: number[], val: number) => {
  const n = nums.length
  let i = 0
  for (let j = 0; j < n; j++) {
    if (nums[j] !== val) {
      nums[i] = nums[j]
      i++
    }
  }

  return i
}
```

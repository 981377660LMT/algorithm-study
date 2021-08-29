添加数组元素：[1,3,4,...,10000]向第二个位置插入 2: 可以将第二个位置 3 赋值为 2，然后在末尾 push(3),起到一个添加元素的效果
删除数组元素：JVM 标记清除垃圾回收算法的技巧：lazy 删除，将要删除的元素标记为删除，等到空间不够时统一删除，大大减少了删除时的数据搬移

**数组的元素在内存地址中是连续的，不能单独删除数组中的某个元素，只能覆盖。**

carry 求和

```JS
/**
 * @param {number[]} digits
 * @return {number[]}
 * @description 给定一个由 整数 组成的 非空 数组所表示的非负整数，在该数的基础上加一。
 */
var plusOne = function (digits: number[]): number[] {
  const n = digits.length
  const res: number[] = []
  let carry = 1
  for (let index = n - 1; index >= 0; index--) {
    const sum = digits[index] + carry
    carry = ~~(sum / 10)
    res.push(sum % 10)
  }
  carry && res.push(1)
  return res.reverse()
}
```

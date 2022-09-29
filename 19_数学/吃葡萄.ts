// 有三种葡萄，每种分别有 a, b, c 颗，现在有三个人，
// 第一个人只吃第一种和第二种葡萄，第二个人只吃第二种和第三种葡萄，第三个人只吃第一种和第三种葡萄。
// 现在给你输入 a, b, c 三个值，请你适当安排，让三个人吃完所有的葡萄，算法返回吃的最多的人最少要吃多少颗葡萄。
const eatGrape = (a: number, b: number, c: number) => {
  const nums = [a, b, c]
  nums.sort((a, b) => a - b)
  const sum = a + b + c
  // 能够构成三角形，可完全平分
  if (nums[0] + nums[1] > nums[2]) {
    return Math.ceil(sum / 3)
  }
  // 不能构成三角形，平分最长边的情况(X 最多吃完 a 和 b，而 c 边需要被 Y 或 Z 平分)
  if (2 * (nums[0] + nums[1]) < nums[2]) {
    return Math.ceil(nums[2] / 2)
  }
  // 不能构成三角形，但依然可以完全平分的情况
  return Math.ceil(sum / 3)
}

// 如果把葡萄的颗数 a, b, c 作为三条线段，它们的大小作为线段的长度，
// 想一想它们可能组成什么几何图形？
// 我们的目的是否可以转化成「尽可能平分这个几何图形的周长」？

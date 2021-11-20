// 如果队列最前面的学生 喜欢 栈顶的三明治，那么会 拿走它 并离开队列。
// 否则，这名学生会 放弃这个三明治 并回到队列的尾部。
function countStudents(students: number[], sandwiches: number[]): number {
  const channel = Array<number>(2).fill(0)

  students.forEach(provider => channel[provider]++)

  for (const consumer of sandwiches) {
    if (channel[consumer] > 0) channel[consumer]--
    else break
  }

  return channel[0] + channel[1]
}

console.log(countStudents([1, 1, 1, 0, 0, 1], [1, 0, 0, 0, 1, 1]))
// 注意到provider(students)其实是不分顺序的
// 只用考虑无法取出数据的情况

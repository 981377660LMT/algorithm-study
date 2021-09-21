// 你的任务是找出最多可以修几门课。

import { MinHeap } from '../../../minheap'

// 贪心的关键点就是：根据谁结束的早来选，选结束的早的，自然能学的课程也就越多
function scheduleCourse(courses: number[][]): number {
  // 按照课程的关闭时间从早到晚排序
  // 使用一个大顶堆来储存已经选择的课程的长度
  // 一旦发现安排了当前课程之后，其结束时间超过了最晚结束时间，那么就从已经安排的课程中，取消掉一门最耗时的课程
  courses.sort((a, b) => a[1] - b[1])
  const pq = new MinHeap((a, b) => b - a)

  let allTime = 0
  for (const course of courses) {
    allTime += course[0]
    pq.push(course[0])
    if (allTime > course[1]) {
      allTime -= pq.shift()!
    }
  }

  return pq.size
}

console.log(
  scheduleCourse([
    [100, 200],
    [200, 1300],
    [1000, 1250],
    [2000, 3200],
  ])
)

export default 1

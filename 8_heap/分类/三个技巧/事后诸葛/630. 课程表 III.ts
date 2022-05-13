// 你的任务是找出最多可以修几门课。

import { MinHeap } from '../../../MinHeap'

// 贪心的关键点就是：根据谁结束的早来选，选结束的早的，自然能学的课程也就越多
// 进来的越菜 出去最菜的
function scheduleCourse(courses: number[][]): number {
  // 按照课程的关闭时间从早到晚排序
  // 使用一个大顶堆来储存已经选择的课程的长度
  // 一旦发现安排了当前课程之后，其结束时间超过了最晚结束时间，那么就从已经安排的课程中，取消掉一门最耗时的课程
  courses.sort((a, b) => a[1] - b[1])
  const pq = new MinHeap<number>((a, b) => b - a)

  let day = 0
  for (const course of courses) {
    day += course[0]
    pq.heappush(course[0])
    if (day > course[1]) {
      day -= pq.heappop()!
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

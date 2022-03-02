模板题 贪心 按结束早的排序

```JS
function maxEvents(events: number[][]): number {
  const sortedEvents = events.slice().sort((a, b) => a[1] - b[1] || a[0] - b[0])
  let pre = sortedEvents[0]
  let res = 1

  for (let i = 1; i < sortedEvents.length; i++) {
    const [_prevStart, prevEnd] = pre
    const [currStart, _currEnd] = sortedEvents[i]
    if (prevEnd <= currStart) {
      res++
      pre = sortedEvents[i]
    }
  }

  return res
}

```

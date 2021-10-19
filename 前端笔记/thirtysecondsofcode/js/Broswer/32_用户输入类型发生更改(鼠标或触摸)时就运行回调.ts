// 使用两个事件监听器
// 假设最初是鼠标输入，并将一个“ touchstart”事件侦听器绑定到文档。
onUserInputChange(type => {
  console.log('The user is now using', type, 'as an input method.')
})

function onUserInputChange(callback: (type: 'mouse' | 'touch') => void) {
  let type: 'mouse' | 'touch' = 'mouse'
  let preTime = 0

  // 假设最初是鼠标输入，并将一个“ touchstart”事件侦听器绑定到文档
  document.addEventListener('touchstart', () => {
    if (type === 'touch') return
    type = 'touch'
    callback(type)
    document.addEventListener('mousemove', mousemoveHandler)
  })

  function mousemoveHandler(): void {
    const now = performance.now()
    // 检测到鼠标移动 监听连续两个“ mousemove”事件在20毫秒内发生的情况。
    if (now - preTime < 20) {
      type = 'mouse'
      callback(type)
      document.removeEventListener('mousemove', mousemoveHandler)
    }
  }
}

// 红灯3秒亮一次，
// 黄灯2秒亮一次，
// 绿灯1秒亮一次；
// 如何让三个灯不断交替重复亮灯？
// （用Promise实现）三个亮灯函数已经存在：
function red() {
  console.log('red')
}
function green() {
  console.log('green')
}
function yellow() {
  console.log('yellow')
}

function light(time: number, callback: () => void) {
  return new Promise<void>(resolve => {
    setTimeout(() => {
      callback()
      resolve()
    }, time)
  })
}

async function main() {
  while (true) {
    await light(100, yellow)
    await light(200, green)
    await light(300, red)
  }
}

main()

export {}

const play = document.getElementById('play')!
const workBtn = document.getElementById('work')!
const interactionBtn = document.getElementById('interaction')!

let iterationCount = 100000000
let value = 0

const expensiveCalculation = function (IdleDeadline: IdleDeadline) {
  // IdleDeadline.timeRemaining()是浏览器空闲时间ms react里留给任务空闲时间是5ms
  while (iterationCount-- > 0 && IdleDeadline.timeRemaining() > 1) {
    value = Math.random() < 0.5 ? value + Math.random() : value + Math.random()
  }
  requestIdleCallback(expensiveCalculation)
}

workBtn.addEventListener('click', function () {
  requestIdleCallback(expensiveCalculation)
})

interactionBtn.addEventListener('click', function () {
  play.style.background = 'palegreen'
})

window.setInterval(() => console.log(iterationCount), 1000)

// requestIdleCallback利用浏览器空闲时间执行任务

export {}

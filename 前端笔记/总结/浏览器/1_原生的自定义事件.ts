let myEvent = new CustomEvent('myEvent', {
  detail: {
    name: 'lindaidai',
  },
})

let btn = document.getElementsByTagName('button')[0]
btn.addEventListener('myEvent', function (e) {
  console.log(e)
  console.log(e.detail)
})

// 事件的触发
setTimeout(() => {
  btn.dispatchEvent(myEvent)
}, 2000)

export {}

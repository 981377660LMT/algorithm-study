const btn = document.getElementsByTagName('button')[0]
btn.addEventListener('myEvent', function (e) {
  console.log(e)
  // @ts-ignore
  console.log(e.detail) // {name: 'lindaidai'}
})

const myEvent = new CustomEvent('myEvent', {
  detail: {
    name: 'lindaidai',
  },
})

// 事件的触发
setTimeout(() => {
  btn.dispatchEvent(myEvent)
}, 500)

export {}

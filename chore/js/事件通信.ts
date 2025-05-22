import EventEmitter from 'events'

let num: number | null = null

const emitter = new EventEmitter()

emitter.on('setFoo', set => {
  set(1)
})

emitter.emit('setFoo', (value: number) => {
  num = value
})

console.log(num)

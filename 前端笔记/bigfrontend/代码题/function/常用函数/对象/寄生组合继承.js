function Parent(name) {
  this.name = name
  this.say = () => {
    console.log(111)
  }
}
Parent.prototype.play = () => {
  console.log(222)
}
function Children(name) {
  Parent.call(this)
  this.name = name
}
Children.prototype = Object.create(Parent.prototype)
Children.prototype.constructor = Children

let child = new Children('111')
console.log(child.name)
child.say()
child.play()

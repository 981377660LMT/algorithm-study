function mergeOptions(options, extendOptions) {}

function Vue() {}

let cid = 0 // 组件的唯一标识
Vue.extend = function (extendOptions) {
  // 创建子类的构造函数 并且调用初始化方法
  const Sub = function VueComponent(options) {
    this._init(options) //调用Vue初始化方法
  }

  Sub.cid = cid++
  Sub.prototype = Object.create(this.prototype) // 子类原型指向父类
  Sub.prototype.constructor = Sub //constructor指向自己
  Sub.options = mergeOptions(this.options, extendOptions) //合并自己的options和父类的options
  return Sub
}

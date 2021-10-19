/**
 *
 * @param {Function} SuperType
 * @param {Function} SubType
 * @returns
 * 寄生组合式继承
 */
const myExtends = (SuperType, SubType) => {
  // your code here
  function myClass(...contructorArgs) {
    // 构造函数初始化
    SuperType.apply(this, contructorArgs)
    SubType.apply(this, contructorArgs)
    // obj.__proto__ = SubType.prototype
    Object.setPrototypeOf(this, SubType.prototype) // 实例方法去SubType.prototype上寻找
  }

  // 找实例方法
  Object.setPrototypeOf(myClass.prototype, SubType.prototype)
  Object.setPrototypeOf(SubType.prototype, SuperType.prototype)

  // 找静态方法
  Object.setPrototypeOf(myClass, SuperType)

  return myClass
}

if (require.main === module) {
  const InheritedSubType = myExtends(SuperType, SubType)
  const instance = new InheritedSubType()

  // 上述代码需要(几乎)和下面的一样
  class SubType extends SuperType {}
  const instance = new SubType()
}

Object.setPrototypeOf

// BFE.dev会用下面的SuperType和 SubType来测试你的代码。
// function SuperType(name) {
//   this.name = name
//   this.forSuper = [1, 2]
//   this.from = 'super'
// }
// SuperType.prototype.superMethod = function() {}
// SuperType.prototype.method = function() {}
// SuperType.staticSuper = 'staticSuper'

// function SubType(name) {
//   this.name = name
//   this.forSub = [3, 4]
//   this.from = 'sub'
// }

// SubType.prototype.subMethod = function() {}
// SubType.prototype.method = function() {}
// SubType.staticSub = 'staticSub'

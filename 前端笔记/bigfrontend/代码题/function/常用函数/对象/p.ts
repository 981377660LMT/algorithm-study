class Dog {
  aa() {}

  static A() {}
}

class Cat extends Dog {}
Dog.prototype.aa()

Cat.A()
Cat.prototype.aa()

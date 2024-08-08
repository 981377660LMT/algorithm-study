/** 两个数据字段的容器：起始血量和攻击字符串. */

interface IBreed {
  getHealth(): number
  getAttack(): number
}

class Breed {
  private readonly _health: number
  private readonly _attack: number

  constructor(health: number, attack: number) {
    this._health = health
    this._attack = attack
  }

  getHealth(): number {
    return this._health
  }

  getAttack(): number {
    return this._attack
  }
}

class Monster {
  private readonly _health: number
  private readonly _breed: IBreed

  // !怪物的品种，取代了之前的子类
  constructor(breed: IBreed) {
    this._breed = breed
    this._health = breed.getHealth()
  }

  getAttack(): number {
    return this._breed.getAttack()
  }
}

// 现在，我们可以直接构造怪物并负责传入它的品种。
// 和常用的OOP语言实现的对象相比这有些退步——我们通常不会分配一块空白内存，然后赋予它类型。
// 相反，我们根据类调用构造器，它负责创建一个新实例。

class Breed2 {
  private readonly _health: number
  private readonly _attack: number

  constructor(health: number, attack: number) {
    this._health = health
    this._attack = attack
  }

  newMonster(): Monster {
    return new Monster(this) // !注意这里
  }

  getHealth(): number {
    return this._health
  }

  getAttack(): number {
    return this._attack
  }
}

export {}

if (require.main === module) {
  const b = new Breed2(100, 20)
  const m1 = new Monster(b) // !1.style1：构造对象然后传入类型对象
  const m2 = b.newMonster() // !2.style2：在类型对象上调用“构造器”函数。由Breed管理monster的创建
}

// 不是简单地调用new,newMonster()函数可以在将控制权传递给Monster初始化之前，从池中或堆中获取内存。
// 通过在唯一有能力创建怪物的Breed函数中放置这些逻辑， 我们保证了所有怪物变量遵守了内存管理规范。

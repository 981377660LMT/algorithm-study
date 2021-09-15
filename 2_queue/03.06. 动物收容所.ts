class AnimalShelf {
  constructor(private catQueue: number[][] = [], private dogQueue: number[][] = []) {}

  // enqueue方法有一个animal参数，
  // animal[0]代表动物编号，animal[1]代表动物种类，其中 0 代表猫，1 代表狗。
  enqueue(animal: number[]): void {
    animal[1] === 0 ? this.catQueue.push(animal) : this.dogQueue.push(animal)
  }

  // 收养人只能收养所有动物中“最老”的动物
  // dequeue*方法返回一个列表[动物编号, 动物种类]，若没有可以收养的动物，则返回[-1,-1]。
  dequeueAny(): number[] {
    const catHead = this.catQueue[0]
    const dogHead = this.dogQueue[0]
    if (!catHead) return this.dogQueue.shift()! || [-1, -1]
    if (!dogHead) return this.catQueue.shift()! || [-1, -1]
    return catHead[0] <= dogHead[0] ? this.catQueue.shift()! : this.dogQueue.shift()!
  }

  dequeueDog(): number[] {
    return this.dogQueue.shift() || [-1, -1]
  }

  dequeueCat(): number[] {
    return this.catQueue.shift() || [-1, -1]
  }
}

/**
 * Your AnimalShelf object will be instantiated and called as such:
 * var obj = new AnimalShelf()
 * obj.enqueue(animal)
 * var param_2 = obj.dequeueAny()
 * var param_3 = obj.dequeueDog()
 * var param_4 = obj.dequeueCat()
 */

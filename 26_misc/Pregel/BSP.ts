export {}

abstract class BarrierStepLoop<E, Id> {
  constructor(protected state: E) {}

  protected abstract plan(currentState: E): Array<() => Id | Promise<Id>>
  protected abstract update(currentState: E, deltas: Id[]): E

  async run(maxSteps: number = 100): Promise<E> {
    let step = 0
    while (step < maxSteps) {
      const tasks = this.plan(this.state)
      if (tasks.length === 0) break
      const deltas = await Promise.all(tasks.map(task => task()))
      this.state = this.update(this.state, deltas)
      step++
    }
    return this.state
  }
}

// --- 示例实现：康威生命游戏 (Game of Life) 的一维简化版 ---
// 规则：如果一个细胞的值 < 10，它会给自己和邻居 +1。直到所有细胞 >= 10。
type CellState = number[]
type CellDelta = { index: number; add: number }[]
class SimpleSimulation extends BarrierStepLoop<CellState, CellDelta> {
  protected plan(currentState: CellState): Array<() => CellDelta> {
    const tasks: Array<() => CellDelta> = []
    currentState.forEach((val, idx) => {
      if (val < 10) {
        tasks.push(() => {
          return [
            { index: idx, add: 1 },
            { index: idx + 1, add: 1 }
          ]
        })
      }
    })
    return tasks
  }

  protected update(currentState: CellState, deltasList: CellDelta[]): CellState {
    const newState = [...currentState]
    deltasList.flat().forEach(change => {
      if (change.index < newState.length) {
        newState[change.index] += change.add
      }
    })
    return newState
  }
}

const sim = new SimpleSimulation([8, 8, 12])
sim.run()

// withTransaction 编程模式
//
// 1. 添加 step -> commit 统一 apply -> state.applyTransaction -> view.update(state)
// 2. 添加 step，再执行的好处：
//    - 指令重排(在线 -> 离线)
//    - 异步api变同步调用
//    - 缺点是无法获取前一次的状态
// 3. withPrevValue 可以获取前一次的状态，但需要做出改变
//    - step 函数
//    - 添加 step 操作内部需要携带状态
// 4. transaction 关注控制，transform 关注计算

// #region Model
type Model = Map<number, number>
// #endregion

// #region transform
interface IStep {
  apply(model: Model): Promise<Model>
}

class CreateStep implements IStep {
  constructor(private readonly _id: number, private readonly _value: number) {}

  async apply(model: Model): Promise<Model> {
    const newModel = new Map(model)
    newModel.set(this._id, this._value)
    return newModel
  }
}

class UpdateStep implements IStep {
  constructor(private readonly _id: number, private readonly _value: number) {}

  async apply(model: Model): Promise<Model> {
    const newModel = new Map(model)
    if (newModel.has(this._id)) {
      newModel.set(this._id, this._value)
      return newModel
    } else {
      throw new Error('not found')
    }
  }
}

class DeleteStep implements IStep {
  constructor(private readonly _id: number) {}

  async apply(model: Model): Promise<Model> {
    const newModel = new Map(model)
    if (newModel.has(this._id)) {
      newModel.delete(this._id)
      return newModel
    } else {
      throw new Error('not found')
    }
  }
}

function isStep(v: any): v is IStep {
  return v instanceof CreateStep || v instanceof UpdateStep || v instanceof DeleteStep
}

// #endregion

// #region transaction

type IStepOrFn = IStep | ((prevValue: Model) => void | Promise<void>)

class Transaction {
  private _model: Model = new Map()

  private readonly _contextStack: IStepOrFn[][] = [[]] // 顶层步骤列表

  private readonly _appliedSteps: IStep[] = []

  async commit(): Promise<void> {
    const topLevelSteps = this._contextStack[0]
    this._model = await this._executeSteps(topLevelSteps, this._model)
  }

  private async _executeSteps(steps: IStepOrFn[], model: Model): Promise<Model> {
    let res = model
    for (const stepOrFn of steps) {
      if (isStep(stepOrFn)) {
        res = await this._applyStep(stepOrFn, res)
        continue
      }
      const newSteps: IStepOrFn[] = []
      this._contextStack.push(newSteps)
      await stepOrFn(res)
      this._contextStack.pop()
      res = await this._executeSteps(newSteps, res)
    }
    return res
  }

  withPrevValue(f: (prevValue: Model) => void | Promise<void>): void {
    this._collect(f)
  }

  // #region api
  create(id: number, value: number): void {
    this._collect(new CreateStep(id, value))
  }

  update(id: number, value: number): void {
    this._collect(new UpdateStep(id, value))
  }

  delete(id: number): void {
    this._collect(new DeleteStep(id))
  }

  private _collect(stepOrFn: IStepOrFn): void {
    this._contextStack[this._contextStack.length - 1].push(stepOrFn)
  }

  private async _applyStep(step: IStep, model: Model): Promise<Model> {
    const res = await step.apply(model)
    this._appliedSteps.push(step)
    return res
  }
  // #endregion
}

async function withTransaction<V>(f: (tr: Transaction) => void | Promise<void>): Promise<void> {
  const tr = new Transaction()
  await f(tr)
  await tr.commit()
}
// #endregion

export {}

if (require.main === module) {
  async function main() {
    await withTransaction(tr => {
      tr.create(101, 2)
      tr.create(100, 100)

      tr.withPrevValue(prevValue => {
        console.log(prevValue, 'prevValue1')
        tr.create(102, 100)
        tr.update(101, 200)
        tr.withPrevValue(prevValue => {
          console.log(prevValue, 'prevValue2')
          tr.update(101, 300)
        })
      })

      tr.withPrevValue(prevValue => {
        console.log(prevValue, 'prevValue3')
        tr.delete(101)
      })
    })
  }

  main().catch(console.error)
}

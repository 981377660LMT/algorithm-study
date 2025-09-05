// withTransaction 编程模式
//
// 1. 添加 step -> commit 统一 apply -> state.applyTransaction -> view.update(state)
// 2. 添加 step，再执行的好处：
//    - 指令重排
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

  private readonly _queue: IStepOrFn[] = []
  private _isInWithPrevValue = false
  private readonly _stagedSteps: IStep[] = []

  private readonly _appliedSteps: IStep[] = []

  async commit(): Promise<void> {
    for (const stepOrFn of this._queue) {
      if (isStep(stepOrFn)) {
        this._model = await this._applyStep(stepOrFn)
        continue
      }
      this._isInWithPrevValue = true
      await stepOrFn(this._model)
      this._isInWithPrevValue = false
      for (const step of this._stagedSteps) {
        this._model = await this._applyStep(step)
      }
      this._stagedSteps.length = 0
    }
  }

  withPrevValue(f: (prevValue: Model) => void | Promise<void>): void {
    if (this._isInWithPrevValue) {
      throw new Error('cannot nest withPrevValue')
    }
    this._queue.push(f)
  }

  // #region api
  create(id: number, value: number): void {
    const step = this._create(id, value)
    this._collect(step)
  }

  private _create(id: number, value: number): CreateStep {
    const step = new CreateStep(id, value)
    return step
  }

  update(id: number, value: number): void {
    const step = this._update(id, value)
    this._collect(step)
  }

  private _update(id: number, value: number): UpdateStep {
    const step = new UpdateStep(id, value)
    return step
  }

  delete(id: number): void {
    const step = this._delete(id)
    this._collect(step)
  }

  private _delete(id: number): DeleteStep {
    const step = new DeleteStep(id)
    return step
  }

  private _collect(step: IStep): void {
    if (this._isInWithPrevValue) {
      this._stagedSteps.push(step)
    } else {
      this._queue.push(step)
    }
  }

  private async _applyStep(step: IStep): Promise<Model> {
    const res = await step.apply(this._model)
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
      })

      tr.withPrevValue(prevValue => {
        console.log(prevValue, 'prevValue2')
        tr.delete(101)
      })
    })
  }

  main().catch(console.error)
}

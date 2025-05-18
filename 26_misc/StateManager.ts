import { create, UseBoundStore } from 'zustand'
import createVanilla, { StoreApi } from 'zustand/vanilla'

type Patch<T> = Partial<{ [P in keyof T]: Patch<T[P]> }>

interface Command<T extends { [key: string]: any }> {
  id?: string
  before: Patch<T>
  after: Patch<T>
}

class StateManager<T extends Record<string, unknown>> {
  /**
   * The initial state.
   */
  private initialState: T

  /**
   * A zustand store that also holds the state.
   */
  private store: StoreApi<T>

  /**
   * The index of the current command.
   */
  protected pointer = -1

  /**
   * The current state.
   */
  private _state: T

  /**
   * The state manager's current status, with regard to restoring persisted state.
   */
  private _status: 'loading' | 'ready' = 'loading'

  /**
   * A stack of commands used for history (undo and redo).
   */
  protected stack: Command<T>[] = []

  /**
   * A snapshot of the current state.
   */
  protected _snapshot: T

  /**
   * A React hook for accessing the zustand store.
   */
  readonly useStore: UseBoundStore<StoreApi<T>>

  /**
   * A promise that will resolve when the state manager has loaded any peristed state.
   */
  ready: Promise<'none' | 'restored' | 'migrated'>

  isPaused = false

  constructor(
    initialState: T,
    id?: string,
    version?: number,
    update?: (prev: T, next: T, prevVersion: number) => T
  ) {
    this._state = deepCopy(initialState)
    this._snapshot = deepCopy(initialState)
    this.initialState = deepCopy(initialState)
    this.store = createVanilla(() => this._state)
    this.useStore = create(this.store)

    this.ready = new Promise<'none' | 'restored' | 'migrated'>(resolve => {
      this._status = 'ready'
      const message = 'none'
      resolve(message)
    }).then(message => {
      if (this.onReady) {
        this.onReady(message)
      }
      return message
    })
  }

  /**
   * Apply a patch to the current state.
   * This does not effect the undo/redo stack.
   * This does not persist the state.
   * @param patch The patch to apply.
   * @param id (optional) An id for the patch.
   */
  private applyPatch = (patch: Patch<T>, id: string, options?: Record<string, any>) => {
    const prev = this._state
    const next = deepMerge(this._state, patch as any)
    if (id === PatchType.UpdateShapeSettings) {
      // 表单配置为覆盖式更新，因为外部表单值变化不会记录字段：undefined，覆盖会导致字段无法删除， 所以需要特殊处理值覆盖
      const shapes = patch?.document?.pages?.[(this as any).currentPageId]?.shapes
      Object.keys(shapes || {}).forEach(shapeId => {
        if (shapes?.[shapeId]?.data) {
          next.document.pages[(this as any).currentPageId].shapes[shapeId].data =
            shapes?.[shapeId]?.data
        }
      })
    }
    const final = this.cleanup(next, prev, patch, id)
    if (this.onStateWillChange) {
      this.onStateWillChange(final, id)
    }
    this._state = final
    this.store.setState(this._state, true)
    if (this.onStateDidChange) {
      this.onStateDidChange(this._state, id, patch, options)
    }
    return this
  }

  // Internal API ---------------------------------

  protected migrate = (next: T): { updateIds: string[]; state: T } => ({
    updateIds: [],
    state: next
  })

  /**
   * Perform any last changes to the state before updating.
   * Override this on your extending class.
   * @param nextState The next state.
   * @param prevState The previous state.
   * @param patch The patch that was just applied.
   * @param id (optional) An id for the just-applied patch.
   * @returns The final new state to apply.
   */
  protected cleanup = (nextState: T, _prevState: T, _patch: Patch<T>, _id?: string): T => nextState

  /**
   * A life-cycle method called when the state is about to change.
   * @param state The next state.
   * @param id An id for the change.
   */
  protected onStateWillChange?: (state: T, id?: string) => void

  /**
   * A life-cycle method called when the state has changed.
   * @param state The next state.
   * @param id An id for the change.
   */
  protected onStateDidChange?: (
    state: T,
    id: string,
    patch: Patch<T> | undefined,
    options?: Record<string, any>
  ) => void

  /**
   * Apply a patch to the current state.
   * This does not effect the undo/redo stack.
   * This does not persist the state.
   * @param patch The patch to apply.
   * @param id (optional) An id for this patch.
   */
  patchState = (patch: Patch<T>, id: PatchType | string): this => {
    if ((this as any).debugPatch) {
      console.log('patchState: ', patch, id)
    }
    const ids = patch.document?.pageStates?.[(this as any).currentPageId]?.selectedIds
    if (ids) {
      ;(this as any).pluginInvoke('onSelectedChange', ids)
    }
    this.applyPatch(patch, id)
    if (this.onPatch) {
      this.onPatch(this._state, patch, id)
    }
    return this
  }

  /**
   * Replace the current state.
   * This does not effect the undo/redo stack.
   * This does not persist the state.
   * @param state The new state.
   * @param id An id for this change.
   */
  protected replaceState = (
    patch: Patch<T> | undefined,
    state: T,
    id: PatchType | string
  ): this => {
    const final = this.cleanup(state, this._state, state, id)
    if (this.onStateWillChange) {
      this.onStateWillChange(final, id)
    }
    this._state = final
    this.store.setState(this._state, true)
    if (this.onStateDidChange) {
      this.onStateDidChange(this._state, id, patch)
    }
    return this
  }

  /**
   * Update the state using a Command.
   * This effects the undo/redo stack.
   * This persists the state.
   * @param command The command to apply and add to the undo/redo stack.
   * @param id (optional) An id for this command.
   */
  protected setState = (command: Command<T>, id = command.id, options?: Record<string, any>) => {
    if (this.pointer < this.stack.length - 1) {
      this.stack = this.stack.slice(0, this.pointer + 1)
    }
    this.stack.push({ ...command, id })
    this.pointer = this.stack.length - 1
    this.applyPatch(command.after, id || 'unknown', options)

    if (this.onCommand) {
      this.onCommand(this._state, command, id)
    }
    return this
  }

  // Public API ---------------------------------

  pause() {
    this.isPaused = true
  }

  resume() {
    this.isPaused = false
  }

  /**
   * A callback fired when the constructor finishes loading any
   * persisted data.
   */
  protected onReady?: (message: 'none' | 'restored' | 'migrated') => void

  /**
   * A callback fired when a patch is applied.
   */
  onPatch?: (state: T, patch: Patch<T>, id?: string) => void

  /**
   * A callback fired when a patch is applied.
   */
  onCommand?: (state: T, command: Command<T>, id?: string) => void

  /**
   * A callback fired when a shape is copied.
   */
  onCopy?: (content?: { shapes: TDShape[]; shapeBounds?: Record<string, TLBounds> }) => void

  /**
   * A callback fired when the state is replaced.
   */
  onReplace?: (state: T) => void

  /**
   * A callback fired when the state is reset.
   */
  onReset?: (state: T) => void

  /**
   * A callback fired when the history is reset.
   */
  onResetHistory?: (state: T) => void

  /**
   * A callback fired when a command is undone.
   */
  onUndo?: (state: T) => void

  /**
   * A callback fired when a command is redone.
   */
  onRedo?: (state: T) => void

  /**
   * Force replace a new undo/redo history. It's your responsibility
   * to make sure that this is compatible with the current state!
   * @param history The new array of commands.
   * @param pointer (optional) The new pointer position.
   */
  replaceHistory = (history: Command<T>[], pointer = history.length - 1): this => {
    this.stack = history
    this.pointer = pointer
    if (this.onReplace) {
      this.onReplace(this._state)
    }
    return this
  }

  /**
   * Reset the history stack (without resetting the state).
   */
  resetHistory = (): this => {
    this.stack = []
    this.pointer = -1
    if (this.onResetHistory) {
      this.onResetHistory(this._state)
    }
    return this
  }

  /**
   * Move backward in the undo/redo stack.
   */
  undo = (): this => {
    if (!this.isPaused) {
      if (!this.canUndo) {
        return this
      }
      const command = this.stack[this.pointer]
      this.pointer--
      this.applyPatch(command.before, 'undo')
    }
    if (this.onUndo) {
      this.onUndo(this._state)
    }
    return this
  }

  /**
   * Move forward in the undo/redo stack.
   */
  redo = (): this => {
    if (!this.isPaused) {
      if (!this.canRedo) {
        return this
      }
      this.pointer++
      const command = this.stack[this.pointer]
      this.applyPatch(command.after, 'redo')
    }
    if (this.onRedo) {
      this.onRedo(this._state)
    }
    return this
  }

  /**
   * Save a snapshot of the current state, accessible at `this.snapshot`.
   */
  setSnapshot = (): this => {
    this._snapshot = { ...this._state }
    return this
  }

  /**
   * Force the zustand state to update.
   */
  forceUpdate = () => {
    this.store.setState(this._state, true)
  }

  /**
   * Get whether the state manager can undo.
   */
  get canUndo(): boolean {
    return this.pointer > -1
  }

  /**
   * Get whether the state manager can redo.
   */
  get canRedo(): boolean {
    return this.pointer < this.stack.length - 1
  }

  /**
   * The current state.
   */
  get state(): T {
    return this._state
  }

  /**
   * The current status.
   */
  get status(): string {
    return this._status
  }

  /**
   * The most-recent snapshot.
   */
  get snapshot(): T {
    return this._snapshot
  }
}

/**
 * Deep copy function for TypeScript.
 * @param T Generic type of target/copied value.
 * @param target Target value to be copied.
 * @see Source project, ts-deeply https://github.com/ykdr2017/ts-deepcopy
 * @see Code pen https://codepen.io/erikvullings/pen/ejyBYg
 */
function deepCopy<T>(target: T): T {
  if (target === null) {
    return target
  }
  if (target instanceof Date) {
    return new Date(target.getTime()) as any
  }

  // First part is for array and second part is for Realm.Collection
  // if (target instanceof Array || typeof (target as any).type === 'string') {
  if (typeof target === 'object') {
    if (typeof target[Symbol.iterator as keyof T] === 'function') {
      const cp = [] as any[]
      if ((target as any as any[]).length > 0) {
        for (const arrayMember of target as any as any[]) {
          cp.push(deepCopy(arrayMember))
        }
      }
      return cp as any as T
    }
    const targetKeys = Object.keys(target)
    const cp = {} as T
    if (targetKeys.length > 0) {
      for (const key of targetKeys) {
        cp[key as keyof T] = deepCopy(target[key as keyof T])
      }
    }
    return cp
  }

  // Means that object is atomic
  return target
}

function deepMerge<T>(target: T, patch: Patch<T>): T {
  const result: T = { ...target }

  const entries = Object.entries(patch) as [keyof T, T[keyof T]][]

  for (const [key, value] of entries) {
    result[key] =
      value === Object(value) && !Array.isArray(value) ? deepMerge(result[key], value) : value
  }

  return result
}

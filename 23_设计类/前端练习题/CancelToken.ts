interface CancellationToken {
  isCancellationRequested: boolean
  throwIfCancellationRequested(): void
  onCancellationRequested(listener: () => void): void
}

export class CancellationTokenSource {
  private _isCancelled = false
  private _listeners: Array<() => void> = []

  get token(): CancellationToken {
    return {
      isCancellationRequested: this._isCancelled,
      throwIfCancellationRequested: () => {
        if (this._isCancelled) throw new Error('OperationCancelled')
      },
      onCancellationRequested: listener => {
        if (this._isCancelled) listener()
        else this._listeners.push(listener)
      }
    }
  }

  cancel() {
    if (!this._isCancelled) {
      this._isCancelled = true
      this._listeners.forEach(l => l())
      this._listeners = []
    }
  }
}

export {}

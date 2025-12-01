import { CancellationToken } from '../types'

export class CancellationTokenSource {
  private _isCancelled = false

  cancel() {
    this._isCancelled = true
  }

  get token(): CancellationToken {
    return {
      isCancellationRequested: this._isCancelled,
      throwIfCancellationRequested: () => {
        if (this._isCancelled) throw new Error('OperationCancelled')
      }
    }
  }
}

class CancelToken {
  public isCancelled = false

  cancel() {
    this.isCancelled = true
  }

  throwIfCancelled() {
    if (this.isCancelled) {
      throw new Error('CancelledError')
    }
  }
}

export {}

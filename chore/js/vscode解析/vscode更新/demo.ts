class TempFile implements Disposable {
  constructor(path: string) {
    // do some work...
  }

  [Symbol.dispose](): void {
    throw new Error('Method not implemented.')
  }
}

export function doSomeWork() {
  // using file = new TempFile(".some_temp_file");

  // use file...

  if (1) {
    // do some more work...
    return
  }
}

function do2() {
  const cleanup = new DisposableStack()
  cleanup.adopt(value, onDispose)
  cleanup.de()
}

doSomeWork()

export {}

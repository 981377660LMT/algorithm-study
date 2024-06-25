function run(fn: () => void, cb: (onInvalidate: (fn: () => void) => void) => void): () => void {
  let cleanup: () => void | undefined

  /** save the cleanup function. */
  const onInvalidate = (fn: () => void): void => {
    cleanup = fn
  }

  const job = () => {
    fn()
    if (cleanup) {
      cleanup()
    }
    cb(onInvalidate)
  }

  return job
}

export {}

if (require.main === module) {
}

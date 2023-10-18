interface IPipe<T> {
  value: T
  to: <R>(func: (value: T) => R) => IPipe<R>
}

function pipe<T>(value: T): IPipe<T> {
  return {
    value,
    to: func => pipe(func(value))
  }
}

export { pipe }

if (require.main === module) {
  // demo:
  const res = pipe(1).to(JSON.stringify).to<number>(JSON.parse).value
  console.log(res)
}

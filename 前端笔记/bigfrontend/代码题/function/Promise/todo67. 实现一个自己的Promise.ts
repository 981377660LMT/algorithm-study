class MyPromise<T> {
  constructor(executor: (resolve: (value: T) => void, reject: (reason?: any) => void) => void) {
    // your code here
  }

  then(onFulfilled: (value: T) => void, onRejected: (reason?: any) => void): MyPromise<T> {
    // your code here
  }

  catch(onRejected: (reason?: any) => void): MyPromise<T> {
    // your code here
  }

  static resolve<T>(value?: T): MyPromise<T> {
    // your code here
  }

  static reject<T = never>(reason?: any): MyPromise<T> {
    // your code here
  }
}

if (require.main === module) {
  const myPromise = new MyPromise()
  Promise.resolve().catch()
  Promise.reject()
}

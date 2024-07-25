// 如何确保如果在错误恢复流程中再次发生错误，这个新的错误也能被相同的逻辑处理
// 递归错误处理.

// export function catchError<T, O extends ObservableInput<any>>(
//   selector: (err: any, caught: Observable<T>) => O
// ): OperatorFunction<T, T | ObservedValueOf<O>> {
//     return operate((source, subscriber) => {
//     let innerSub: Subscription | null = null;
//     let syncUnsub = false;
//     let handledResult: Observable<ObservedValueOf<O>>;

//     innerSub = source.subscribe(
//       createOperatorSubscriber(subscriber, undefined, undefined, (err) => {
//         handledResult = innerFrom(selector(err, catchError(selector)(source)));
//         if (innerSub) {
//           innerSub.unsubscribe();
//           innerSub = null;
//           handledResult.subscribe(subscriber);
//         } else {
//           // We don't have an innerSub yet, that means the error was synchronous
//           // because the subscribe call hasn't returned yet.
//           syncUnsub = true;
//         }
//       })
//     );

//     if (syncUnsub) {
//       // We have a synchronous error, we need to make sure to
//       // finalize right away. This ensures that callbacks in the `finalize` operator are called
//       // at the right time, and that finalization occurs at the expected
//       // time between the source error and the subscription to the
//       // next observable.
//       innerSub.unsubscribe();
//       innerSub = null;
//       handledResult!.subscribe(subscriber);
//     }
//   });
// }

export {}

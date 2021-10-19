// 创建一个新的 MutationObserver，并为指定元素上的每个变异运行提供的回调。
const observeMutations = (
  element: Node,
  callback: {
    (...data: any[]): void
    (message?: any, ...optionalParams: any[]): void
    (arg0: MutationRecord): void
  },
  options?: undefined
) => {
  const observer = new MutationObserver(mutations => mutations.forEach(m => callback(m)))
  observer.observe(
    element,
    Object.assign(
      {
        childList: true,
        attributes: true,
        attributeOldValue: true,
        characterData: true,
        characterDataOldValue: true,
        subtree: true,
      },
      options
    )
  )
  return observer
}
const obs = observeMutations(document, console.log)
// Logs all mutations that happen on the page
obs.disconnect()
// Disconnects the observer and stops logging mutations on the page

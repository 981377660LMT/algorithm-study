const outer = () => {
  const inner = () => console.trace('Hello')
  inner()
}

outer()

Array.from({ length: 4 }).forEach(
  () => console.count('items') // Call the counter labelled 'items'
)

/*
  items: 1
  items: 2
  items: 3
  items: 4
*/
console.countReset('items') // Reset the counter labelled 'items'

console.log(
  'CSS can make %cyour console logs%c %cawesome%c!', // String to format
  // Each string is the CSS to apply for each consecutive %c
  'color: #fff; background: #1e90ff; padding: 4px', // Apply styles
  '', // Clear any styles
  'color: #f00; font-weight: bold', // Apply styles
  '' // Clear any styles
)

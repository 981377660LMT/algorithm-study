// 1.Using the History API ：History API 只允许同源的 url，因此无法导航到完全不同的网站。
// 使用 history.pushState ()或 history.replaceState ()来修改浏览器中的 URL
// Current URL: https://my-website.com/page_a
const nextURL = 'https://my-website.com/page_b'
const nextTitle = 'My new page title'
const nextState = { additionalInformation: 'Updated the URL with JS' }

// This will create a new entry in the browser's history, without reloading
window.history.pushState(nextState, nextTitle, nextURL)

// This will replace the current entry in the browser's history, without reloading
window.history.replaceState(nextState, nextTitle, nextURL)

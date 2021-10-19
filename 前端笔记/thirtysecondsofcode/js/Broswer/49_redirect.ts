const redirect = (url: string, asLink = true) =>
  asLink ? (window.location.href = url) : window.location.replace(url)

redirect('https://google.com')

// href是Hypertext Reference的缩写  意思是指定超链接目标的URL

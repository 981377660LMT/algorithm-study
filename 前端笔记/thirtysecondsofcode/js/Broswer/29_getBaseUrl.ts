const getBaseURL = (url: string) => url.replace(/[?#].*$/, '')

console.log(getBaseURL('http://url.com/page?name=Adam&surname=Smith'))

const loves = $x("//li[@class='list-group-item question']//span")
const deleteButton = $x("//div[@class='lc-alert-button-group__2THN ']")[0].firstElementChild

function sleep(ms = 500) {
  return new Promise(resolve =>
    setTimeout(() => {
      resolve()
    }, ms)
  )
}

async function main() {
  for (let i = 0; i < loves.length; i++) {
    loves[i].click()
    await sleep(50)
    deleteButton.click()
    await sleep(50)
  }
}

main()

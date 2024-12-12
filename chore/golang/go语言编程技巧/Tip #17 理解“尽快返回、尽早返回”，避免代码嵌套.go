// https://colobu.com/gotips/017.html

package main

func main() {

}

// !1. 提前处理错误，别让他们碍事。
func demo1() {
	var fileExists func(filename string) bool
	var readFile func(filename string) ([]byte, error)

	// Bad:
	{
		if fileExists("filename") {
			content, readErr := readFile("filename")
			if readErr != nil {
				// handle error
			} else {
				// process content
				_ = content
			}
		} else {
			// handle error
		}
	}

	// !Good:
	{
		if !fileExists("filename") {
			// handle file not exist error
			return
		}

		content, readErr := readFile("filename")
		if readErr != nil {
			// handle read error
			return
		}

		// process content
		_ = content
	}

}

// !2. 把初始化和错误检测分开
// 就算user仅仅用在else的作用域里，我也建议把初始化和错误检测分开。
// 这样可以避免深层的嵌套以及可以简化错误的处理。
func demo2() {

	var fetchUser func(id int) (string, error)

	// Bad:
	func() error {
		{
			if user, err := fetchUser(1); err != nil {
				return err
			} else {
				// process user
				_ = user
			}

			// continue to process user
			return nil
		}
	}()

	// Good:
	func() error {
		user, err := fetchUser(1)
		if err != nil {
			return err
		}

		// process user
		_ = user

		// continue to process user
		return nil
	}()
}

// !3. 提取出一个只返回error的函数，将处理逻辑限制在函数内部，不影响外部逻辑
func demo3() {
	var fetchUser func(id int) (string, error)

	// Good:
	doSomethingWithUser := func(id string) error {
		user, err := fetchUser(1)
		if err != nil {
			return err
		}

		// process user
		_ = user

		return nil
	}

	if err := doSomethingWithUser("1"); err != nil {
		// handle error
	}

	// continue to process user
}

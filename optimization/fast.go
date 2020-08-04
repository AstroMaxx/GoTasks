package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.WriteString(out, "found users:\n")
	seenBrowsers := make([]string, 0, 114)
	var browser string
	var browsers []string
	var notSeenBefore, isAndroid, isMSIE bool
	var i = -1
	//lines := make([][]byte, 0, 1000)
	//lines = bytes.Split(fileContents, []byte("\n"))

	//users := make([]User, 0, 1000)

	lines := bufio.NewScanner(file)
	user := User{}
	
	for lines.Scan() {
		line := lines.Bytes()
		// fmt.Printf("%v %v\n", err, line)
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		//users = append(users, user)

		//for i, user := range users {
		i++
		isAndroid = false
		isMSIE = false

		browsers = user.Browsers

		for _, browserRaw := range browsers {
			browser = browserRaw
			notSeenBefore = false
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore = true

				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

			for _, browserRaw := range browsers {
				browser = browserRaw
				notSeenBefore = false
				if strings.Contains(browser, "MSIE") {
					isMSIE = true
					notSeenBefore = true

					for _, item := range seenBrowsers {
						if item == browser {
							notSeenBefore = false
						}
					}
					if notSeenBefore {
						// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
						seenBrowsers = append(seenBrowsers, browser)
					}
				}
			}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		io.WriteString(out, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email+">\n")
	}

	io.WriteString(out, "\nTotal unique browsers " + strconv.Itoa(len(seenBrowsers)) + "\n")
}

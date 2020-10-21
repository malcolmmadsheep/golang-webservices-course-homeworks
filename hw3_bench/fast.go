package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var foundUsersBuff bytes.Buffer
	var user User

	seenBrowsers := make([]string, 0, 128)
	scanner := bufio.NewScanner(file)
	i := -1

	for scanner.Scan() {
		user.UnmarshalJSON(scanner.Bytes())
		i++
		isAndroid := false
		isMSIE := false

		if user.Browsers == nil || len(user.Browsers) == 0 {
			continue
		}

		for _, browser := range user.Browsers {
			thisIsAndroid := strings.Contains(browser, "Android")
			thisIsMSIE := strings.Contains(browser, "MSIE")

			if thisIsAndroid {
				isAndroid = true
			}

			if thisIsMSIE {
				isMSIE = true
			}

			if !(thisIsAndroid || thisIsMSIE) {
				continue
			}

			notSeenBefore := true

			for _, seenBrowser := range seenBrowsers {
				if browser == seenBrowser {
					notSeenBefore = false
				}
			}

			if notSeenBefore {
				seenBrowsers = append(seenBrowsers, browser)
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		foundUsersBuff.WriteByte('[')
		foundUsersBuff.WriteString(strconv.Itoa(i))
		foundUsersBuff.WriteString("] ")
		foundUsersBuff.WriteString(user.Name)
		foundUsersBuff.WriteString(" <")
		writeEmail(&foundUsersBuff, user.Email)
		foundUsersBuff.WriteString(">\n")

		user.Clear()
	}

	fmt.Fprintln(out, "found users:")
	foundUsersBuff.WriteTo(out)
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func writeEmail(out *bytes.Buffer, email string) {
	for j := 0; j < len(email); j++ {
		switch email[j] {
		case '@':
			out.WriteString(" [at] ")
			break
		default:
			out.WriteByte(email[j])
			break
		}
	}
}

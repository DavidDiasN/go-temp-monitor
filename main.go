package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	files := []string{}

	for i := 1; i < 8; i++ {
		files = append(files, fmt.Sprintf("temp%d_input", i))
	}

	sum := 0
	successes := 0
	max := 0
	failures := 0

	for {

		for _, fname := range files {

			file, err := os.Open(fmt.Sprintf("/sys/devices/platform/coretemp.0/hwmon/hwmon4/%s", fname))

			if err != nil {
				fmt.Println("open error")
				failures++
				continue
			}
			defer file.Close()

			data := make([]byte, 8)
			_, err = file.Read(data)
			if err != nil {
				fmt.Println("Read error")
				failures++
				continue
			}
			dataString := strings.ReplaceAll(string(data), "\n", "")
			dataString = strings.ReplaceAll(dataString, "\x00", "")

			num, err := strconv.Atoi(dataString)
			if err != nil {

				fmt.Println("Conversion error")
				fmt.Println(err)
				failures++
				continue
			}

			sum += num
			successes++
			if num > max {
				max = num
			}
		}
		file, err := os.OpenFile("/home/david/.statStore/cpu_temps/i3_info", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error oepning file:", err)
			return
		}
		defer file.Close()
		average := 0
		if successes != 0 {
			average = sum / successes
		}

		_, err = file.WriteString(fmt.Sprintf("%d\n%d\n%d", average, max, failures))
		if err != nil {
			fmt.Println("Error writing to file: ", err)
		}

		time.Sleep(5 * time.Second)
	}
}

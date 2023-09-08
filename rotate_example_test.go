package logrotate_test

import (
	"fmt"
	"github.com/gtoxlili/logrotate"
	"time"
)

func ExampleNewWriter() {

	rw, _ := logrotate.NewWriter(
		logrotate.Secondly,
		"tmp",
		"app.log")
	defer rw.Close()

	for i := 0; i < 300; i++ {
		time.Sleep(time.Millisecond * 100)
		rw.Write([]byte(fmt.Sprintf("now is %s\n", time.Now().String())))
	}
	fmt.Println("done")

	// Output: done
}

# Logrotate

Logrotate is an Go library that provides log rotation functionality. It offers high performance by utilizing mmap for log file writing and implements various optimizations for time-based operations. It implements the io.WriteCloser interface, making it easy to integrate with any logging library. It is recommended to use Logrotate with the [Diode](https://github.com/cloudfoundry/go-diodes) library for achieving thread-safe and non-blocking log rotation.

## Usage

Here's an example of how to use Logrotate with Zerolog:

```go
package main

import (
	"github.com/gtoxlili/logrotate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	rw, _ := logrotate.NewWriter(logrotate.Hourly, "/path/to/logs", "myapp.log")

	wr := diode.NewWriter(rw, 1000, 100*time.Millisecond, func(missed int) {
		log.Error().Msgf("Dropped %d messages", missed)
	})

	log.Logger = zerolog.New(wr)
}
```

## License

Logrotate is licensed under the [MIT License](LICENSE).

## Contributing

Contributions to Logrotate are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

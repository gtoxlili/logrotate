# Logrotate

Logrotate is an open-source library that provides log rotation functionality. It offers high performance by utilizing mmap for log file writing and implements various optimizations for time-based operations. It implements the io.ReadCloser interface, making it easy to integrate with any logging library. It is recommended to use Logrotate with the [Diode](https://github.com/cloudfoundry/go-diodes) library for achieving thread-safe and non-blocking log rotation.


## Features

- High Performance: Logrotate leverages `mmap` for efficient log file writing, resulting in improved performance.
- Time-Based Rotation: Logrotate supports rotation based on different time units such as minutes, hours, and days. It automatically rotates logs based on the specified threshold.
- Customizable Filename Format: Logrotate allows users to define a custom filename format for rotated log files, making it flexible to match specific requirements.
- Integration with Any Logging Library: Logrotate implements the io.ReadCloser interface, enabling easy integration with any logging library.


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

	log.Logger = zerolog.New(wr).With().Timestamp().Logger()
}
```

## Performance Considerations

Logrotate is designed to provide high performance for log rotation tasks. By utilizing `mmap`, it minimizes the overhead of file I/O operations. Additionally, Logrotate optimizes time-based operations by using a monotonic clock and efficiently determining when log rotation is needed.


## License

Logrotate is licensed under the [MIT License](LICENSE).

## Contributing

Contributions to Logrotate are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

# Go

cd go

go run -ldflags="-s -w" .

**Results:**

100000 threads finished 10000 iterrations each in 47.01 seconds: map[string]SomeData

100000 threads finished 10000 iterrations each in 58.05 seconds: map[string]*SomeData

RAM: 1.5Gb - 4Gb

# Rust

cd rust

cargo run --release

**Results:**

100000 tasks finished 10000 iterrations each in 68.1729521s

RAM: 35Mb - 60Mb

# Test Description

We run 100'000 tasks, in each task 10'000 small structs created, inserted into a hash-map, and after that retrieved from the hash-map by the key.

**Go:**

```
cd go
go run -ldflags="-s -w" .
```

**Rust:**

```
cd rust
cargo run --release
```

# Test Results

Windows 10 Pro, Intel(R) Core(TM) i7-9850H CPU @2.60GHz

**Go (goroutines):**
 - With pure goroutines: finished in 46.61s, task avg 16.77s, min 0.00s, max 46.31s, RAM: 2000Mb - 4000Mb
 - With CPU workers: finished in 69.23s, task avg 0.0079s, min 0.0000s, max 0.0972s, RAM 200-1000Mb (1000Mb only at start, tend to go down to 200Mb when running)

**Rust (tokio tasks):**
 - With default memalloc: finished in 67.67s, task avg 6ms, min 3ms, max 53ms, RAM: 35Mb - 60Mb
 - With mimalloc: finished in 48.65s, task avg 4ms, min 3ms, max 59ms, RAM: 78Mb

![Chart](assets/chart1.png)

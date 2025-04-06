# Test Description

We run 100'000 tasks, in each task 10'000 small structs created, inserted into a hash-map, and after that retrieved from the hash-map by the key.

**Key points:**

Rust was 30% slower with the default malloc, but almost identical to Go with mimalloc. While the biggest difference was massive RAM usage by Go: 2-4Gb vs Rust only 30-60Mb. But why? Is that simply because GC can't keep up with so many goroutines allocating structs?

Notice that on average Rust finishes a task in 0.006s (max in 0.053s), while Go's average task duration is 16s! Massive differrence! If both finished all tasks at roughtly the same time that means that Go is trying to execute thousands of them in parallel (sharing the same 12 CPU threads), while Rust is most likely running only 12 tasks at once. This explains why Rust's average task duration is so short. 

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

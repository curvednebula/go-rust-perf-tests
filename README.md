# Test Description

We run 100'000 tasks, in each task 10'000 small structs created, inserted into a hash-map, and after that retrieved from the hash-map by the key.

**Key points:**

Rust was 30% slower with the default malloc, but almost identical to Go with mimalloc. While the biggest difference was massive RAM usage by Go: 2-4Gb vs Rust is only 35-60Mb. But why? Is that simply because GC can't keep up with so many goroutines allocating structs?

Notice that on average Rust finished a task in 0.006s (max in 0.053s), while Go's average task duration was 16s! A massive differrence! If both finished all tasks at roughtly the same time that could only mean that Go is executing thousands of tasks in parallel sharing limited amount of CPU threads available, but Rust is running only couple of them at once. This explains why Rust's average task duration is so short.

Since Go runs so many tasks in paralell it keeps thousands of hash maps filled with thousands of structs in the RAM. GC can't even free this memory because application is still using it. Rust on the other hand only creates couple of hash maps at once.

## Optimization 1 ##

To solve the problem I've created a simple utility: CPU workers. It limits number of parallel tasks executed to be not more than the number of CPU threads. With this optimization Go's memory usage dropped to 1000Mb at start and it goes down to 200Mb as test runs. This is at least 4 times better than before. Interestingly enough with this optimization execution time increased from 46 sec to 70 sec, which almost identical to Rust's time with default malloc.



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

Platform: Windows 10 Pro, Intel(R) Core(TM) i7-9850H CPU @2.60GHz

## Start 10 tasks each millisec

Start 10 tasks each millisec, wich is ~10'000 tasks per second

**Go (goroutines):**
 - Goroutines only: finished in 46.8229s, task avg 0.1611s, min 0.0000s, max 3.0188s, RAM: up to 400Mb
 - Goroutines + CPU workers: finished in 74.9265s, task avg 0.0090s, min 0.0000s, max 0.0906s, RAM: up to 500Mb
 - CPU workers only: finished in 62.1834s, task avg 0.0074s, min 0.0005s, max 0.1315s, RAM: up to 35Mb

**Rust (tokio tasks):**
 - With std malloc: finished in 66.4028708s, task avg 6ms, min 4ms, max 43ms, RAM: up to 57Mb
 - With mimalloc: finished in 47.7435647s, task avg 4ms, min 3ms, max 39ms, RAM: up to 77Mb

 ![Chart](assets/10-tasks-per-ms.png)

## Instant Burst

Just start all 100'000 tasks as quick as possible.

**Go (goroutines):**
 - Goroutines only: finished in 46.61s, task avg 16.77s, min 0.00s, max 46.31s, RAM: 2000Mb - 4000Mb
 - Goroutines + CPU workers: finished in 69.23s, task avg 0.0079s, min 0.0000s, max 0.0972s, RAM 200-1000Mb (1000Mb at start, tend to go down to 200Mb when running)
 - CPU workers only: finished in 60.7386s, task avg 0.0072s, min 0.0022s, max 0.0399s, RAM: up to 35Mb

**Rust (tokio tasks):**
 - With std memalloc: finished in 67.67s, task avg 6ms, min 3ms, max 53ms, RAM: 35Mb - 60Mb
 - With mimalloc: finished in 48.65s, task avg 4ms, min 3ms, max 59ms, RAM: 78Mb

![Chart](assets/instant-burst.png)


use mimalloc::MiMalloc;
use std::{collections::HashMap, thread, time::Duration};
use tokio::{task::JoinSet, time::Instant};

#[global_allocator]
static GLOBAL: MiMalloc = MiMalloc;

const TASKS_NUM: u32 = 100_000;
const ITEMS_NUM: u32 = 10_000;
const TASKS_IN_BUNCH: u32 = 10;
const TIME_BETWEEN_BUNCHES_MS: u64 = 1;

struct SomeData {
    name: String,
    num: u32,
}

async fn run_test() {
    let start = Instant::now();
    let mut join_set = JoinSet::new();

    for task_idx in 0..TASKS_NUM {
        join_set.spawn(async move {
            let task_start = Instant::now();
            let mut map = HashMap::new();
            let mut _sum: u64 = 0;

            for j in 0..ITEMS_NUM {
                let name = j.to_string(); // same performance: format!("{}", j);

                map.insert(
                    name.clone(),
                    SomeData {
                        name: name.clone(),
                        num: j,
                    },
                );

                let val = map.get(&name);
                if let Some(value) = val {
                    if value.name == name {
                        _sum += value.num as u64;
                    }
                }
            }
            return task_start.elapsed();
        });

        if task_idx % TASKS_IN_BUNCH == 0 {
            thread::sleep(Duration::from_millis(TIME_BETWEEN_BUNCHES_MS));
        }
    }

    let mut num_results = 0;
    let mut all_tasks_time= Duration::ZERO;
    let mut min_time = Duration::MAX;
    let mut max_time = Duration::ZERO;

    while let Some(Ok(task_time)) = join_set.join_next().await {
        all_tasks_time += task_time;
        if min_time > task_time {
            min_time = task_time;
        }
        if max_time < task_time {
            max_time = task_time;
        }
        num_results += 1;
    }

    assert!(num_results == TASKS_NUM);

    let total_duration = start.elapsed();
    let avg_time = all_tasks_time / num_results;

    println!(
        "- finished in {:?}, task avg {:?}, min {:?}, max {:?}",
        total_duration, avg_time, min_time, max_time
    );
}

#[tokio::main]
async fn main() {
    run_test().await;
    run_test().await;
    run_test().await;
}

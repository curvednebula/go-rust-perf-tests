use mimalloc::MiMalloc;
use std::{collections::HashMap, thread, time::Duration, u128};
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

#[tokio::main]
async fn main() {
    let start = Instant::now();
    let mut join_set = JoinSet::new();

    for task_idx in 0..TASKS_NUM {
        join_set.spawn(async move {
            let task_start = Instant::now();
            let mut map = HashMap::new();
            let mut _sum: u64 = 0;

            for j in 0..ITEMS_NUM {
                let name = format!("name-{}", j);

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
    let mut all_tasks_time: u128 = 0;
    let mut min_time: u128 = u128::MAX;
    let mut max_time: u128 = u128::MIN;

    while let Some(res) = join_set.join_next().await {
        let val = res.unwrap();

        let task_time = val.as_millis();
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
    let avg_time = all_tasks_time / (num_results as u128);

    println!(
        "{} tasks, {} items: finished in {:?}, task avg {:?}, min {:?}, max {:?}",
        num_results,
        ITEMS_NUM,
        total_duration,
        Duration::from_millis(avg_time as u64),
        Duration::from_millis(min_time as u64),
        Duration::from_millis(max_time as u64)
    );
}

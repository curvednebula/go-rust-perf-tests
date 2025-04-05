use std::{collections::HashMap, time::Duration, u128};
use tokio::{task::JoinSet, time::Instant};

const TASKS_NUM: u32 = 100_000;
const VALUES_NUM: u32 = 10_000;

struct SomeData {
    name: String,
    num: u32,
}

#[tokio::main]
async fn main() {
    let start = Instant::now();
    let mut set = JoinSet::new();

    for _ in 0..TASKS_NUM {
        set.spawn(async move {
            let mut map = HashMap::new();
            let mut sum: u64 = 0;

            for j in 0..VALUES_NUM {
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
                        sum += value.num as u64;
                    }
                }
            }
            return start.elapsed();
        });
    }

    let mut all_tasks_time: u128 = 0;
    let mut min_time: u128 = u128::MAX;
    let mut max_time: u128 = u128::MIN;
    let mut num_results = 0;

    while let Some(res) = set.join_next().await {
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

    let duration = start.elapsed();
    let avg_task_completed_in = all_tasks_time / (num_results as u128);

    println!(
        "{} tasks, {} iterrations in each: finished in {:?}, one task avg {:?}, min {:?}, max {:?}",
        num_results,
        VALUES_NUM,
        duration,
        Duration::from_millis(avg_task_completed_in as u64),
        Duration::from_millis(min_time as u64),
        Duration::from_millis(max_time as u64)
    );
}

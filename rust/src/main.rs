use futures::future::join_all;
use std::{collections::HashMap, time::Duration, u128};
use tokio::time::Instant;

const TASKS_NUM: u32 = 100_000;
const VALUES_NUM: u32 = 10_000;

struct SomeData {
    name: String,
    num: u32,
}

#[tokio::main]
async fn main() {
    let start = Instant::now();

    let tasks: Vec<_> = (0..TASKS_NUM)
        .map(|_| {
            tokio::spawn(async move {
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
            })
        })
        .collect();

    let results = join_all(tasks).await;

    let mut all_tasks_time: u128 = 0;
    let mut min_time: u128 = u128::MAX;
    let mut max_time: u128 = u128::MIN;

    for result in &results {
        match result {
            Ok(val) => {
                let task_time = val.as_millis();
                all_tasks_time += task_time;
                if min_time > task_time {
                    min_time = task_time;
                }
                if max_time < task_time {
                    max_time = task_time;
                }
            }
            Err(err) => eprintln!("Error: {:?}", err),
        }
    }

    let duration = start.elapsed();

    let avg_task_completed_in = all_tasks_time / (results.len() as u128);

    println!(
        "{} tasks, {} iterrations in each: finished in {:?}, one task avg {:?}, min {:?}, max {:?}",
        TASKS_NUM,
        VALUES_NUM,
        duration,
        Duration::from_millis(avg_task_completed_in as u64),
        Duration::from_millis(min_time as u64),
        Duration::from_millis(max_time as u64)
    );
}

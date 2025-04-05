use futures::future::join_all;
use std::collections::HashMap;
use tokio::time::Instant;

const THREADS_NUM: u32 = 100_000;
const VALUES_NUM: u32 = 10_000;

struct SomeData {
    name: String,
    num: u32,
}

#[tokio::main]
async fn main() {
    let start = Instant::now();

    let tasks: Vec<_> = (0..THREADS_NUM)
        .map(|_| {
            tokio::spawn(async move {
                let mut map = HashMap::new();
                let mut sum: u32 = 0;

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
                        //
                    }
                }
            })
        })
        .collect();

    join_all(tasks).await;

    let duration = start.elapsed();

    println!(
        "{} tasks finished {} iterrations each in {:?}",
        THREADS_NUM, VALUES_NUM, duration
    );
}

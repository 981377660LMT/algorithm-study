use std::collections::HashMap;

struct UndergroundSystem {
    times: HashMap<String, HashMap<String, (i32, i32)>>,
    passengers: HashMap<i32, (String, i32)>,
}

impl UndergroundSystem {
    fn new() -> Self {
        Self {
            times: HashMap::new(),
            passengers: HashMap::new(),
        }
    }

    fn check_in(&mut self, id: i32, station_name: String, t: i32) {
        self.passengers.insert(id, (station_name, t));
    }

    fn check_out(&mut self, id: i32, station_name: String, t: i32) {
        if let Some((start_station, start_t)) = self.passengers.remove(&id) {
            let (sum, count) = self
                .times
                .entry(start_station)
                .or_default()
                .entry(station_name)
                .or_default();
            *sum += t - start_t;
            *count += 1;
        }
    }

    fn get_average_time(&mut self, start_station: String, end_station: String) -> f64 {
        let (sum, count) = self
            .times
            .entry(start_station)
            .or_default()
            .entry(end_station)
            .or_default();

        *sum as f64 / *count as f64
    }
}

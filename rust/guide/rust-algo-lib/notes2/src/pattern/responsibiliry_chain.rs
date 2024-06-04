// 责任链允许多个对象来对请求进行处理，而无需让发送者类与具体接收者类相耦合。
// 链可在运行时由遵循标准处理者接口的任意处理者动态生成。
//
// 例子：
// Patient -> Reception -> Doctor -> Medical -> Cashier
// !Box 允许以任意组合方式进行链式调用.
// let mut reception = Reception::new(doctor);
// let mut reception = Reception::new(cashier);

mod patient {
    #[derive(Default)]
    pub struct Patient {
        pub name: String,
        pub registration_done: bool,
        pub doctor_checkup_done: bool,
        pub medicine_done: bool,
        pub payment_done: bool,
    }
}

mod department {
    use super::patient::Patient;

    pub trait IDepartment {
        fn handle(&mut self, patient: &mut Patient);
        fn next(&mut self) -> &mut Option<Box<dyn IDepartment>>;

        fn excute(&mut self, patient: &mut Patient) {
            self.handle(patient);
            if let Some(next) = self.next() {
                next.excute(patient);
            }
        }
    }

    /// Helps to wrap an object into a boxed type.
    pub fn into_next(department: impl IDepartment + 'static) -> Option<Box<dyn IDepartment>> {
        Some(Box::new(department))
    }
}

mod cashier {
    use super::department::IDepartment;
    use super::patient::Patient;

    #[derive(Default)]
    pub struct Cashier {
        next: Option<Box<dyn IDepartment>>,
    }

    impl IDepartment for Cashier {
        fn handle(&mut self, patient: &mut Patient) {
            if !patient.payment_done {
                println!("Cashier: Patient payment done");
                patient.payment_done = true;
            }
        }

        fn next(&mut self) -> &mut Option<Box<dyn IDepartment>> {
            &mut self.next
        }
    }
}

mod doctor {
    use super::department::{into_next, IDepartment};
    use super::patient::Patient;

    pub struct Doctor {
        next: Option<Box<dyn IDepartment>>,
    }

    impl Doctor {
        pub fn new(next: impl IDepartment + 'static) -> Self {
            Self {
                next: into_next(next),
            }
        }
    }

    impl IDepartment for Doctor {
        fn handle(&mut self, patient: &mut Patient) {
            if !patient.doctor_checkup_done {
                println!("Doctor: Patient doctor checkup done");
                patient.doctor_checkup_done = true;
            }
        }

        fn next(&mut self) -> &mut Option<Box<dyn IDepartment>> {
            &mut self.next
        }
    }
}

mod medical {
    use super::department::{into_next, IDepartment};
    use super::patient::Patient;

    pub struct Medical {
        next: Option<Box<dyn IDepartment>>,
    }

    impl Medical {
        pub fn new(next: impl IDepartment + 'static) -> Self {
            Self {
                next: into_next(next),
            }
        }
    }

    impl IDepartment for Medical {
        fn handle(&mut self, patient: &mut Patient) {
            if !patient.medicine_done {
                println!("Medical: Patient medicine done");
                patient.medicine_done = true;
            }
        }

        fn next(&mut self) -> &mut Option<Box<dyn IDepartment>> {
            &mut self.next
        }
    }
}

mod reception {
    use super::department::{into_next, IDepartment};
    use super::patient::Patient;

    pub struct Reception {
        next: Option<Box<dyn IDepartment>>,
    }

    impl Reception {
        pub fn new(next: impl IDepartment + 'static) -> Self {
            Self {
                next: into_next(next),
            }
        }
    }

    impl IDepartment for Reception {
        fn handle(&mut self, patient: &mut Patient) {
            if !patient.registration_done {
                println!("Reception: Patient registration done");
                patient.registration_done = true;
            }
        }

        fn next(&mut self) -> &mut Option<Box<dyn IDepartment>> {
            &mut self.next
        }
    }
}

fn main() {
    use cashier::Cashier;
    use department::IDepartment;
    use doctor::Doctor;
    use medical::Medical;
    use patient::Patient;
    use reception::Reception;

    let cashier = Cashier::default();
    let medical = Medical::new(cashier);
    let doctor = Doctor::new(medical);
    let mut reception = Reception::new(doctor);

    let mut patient = Patient {
        name: "John".into(),
        ..Patient::default()
    };

    reception.excute(&mut patient);
}

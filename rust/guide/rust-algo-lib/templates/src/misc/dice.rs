#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Dice<T> {
    pub top: T,
    pub bottom: T,
    pub front: T,
    pub back: T,
    pub right: T,
    pub left: T,
}

impl<T: Clone> Dice<T> {
    pub fn new(top: T, bottom: T, front: T, back: T, right: T, left: T) -> Self {
        Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        }
    }

    pub fn rot_left(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(right, left, front, back, bottom, top)
    }

    pub fn rot_right(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(left, right, front, back, top, bottom)
    }

    pub fn rot_front(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(back, front, top, bottom, right, left)
    }

    pub fn rot_back(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(front, back, bottom, top, right, left)
    }

    pub fn rot_clockwise(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(top, bottom, right, left, back, front)
    }

    pub fn rot_counterclockwise(&self) -> Self {
        let Dice {
            top,
            bottom,
            front,
            back,
            right,
            left,
        } = self.clone();
        Dice::new(top, bottom, left, right, front, back)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test() {
        let dice = Dice::new(1, 6, 2, 5, 3, 4);

        assert_eq!(dice.rot_right(), dice.rot_left().rot_left().rot_left());
        assert_eq!(dice.rot_front(), dice.rot_back().rot_back().rot_back());
        assert_eq!(
            dice.rot_clockwise(),
            dice.rot_counterclockwise()
                .rot_counterclockwise()
                .rot_counterclockwise()
        )
    }
}

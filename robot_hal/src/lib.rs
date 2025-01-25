//! Hardware Abstraction Layer for Robot Control
//! This is the main library file that exports all modules

mod error;
mod motor;
mod sensor;

pub use error::HardwareError;
pub use motor::{Motor, MotorDirection, N20Motor};
pub use sensor::{Sensor, DummySensor};
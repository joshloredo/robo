//! Error types for hardware operations
use std::error::Error;
use std::fmt;

#[derive(Debug)]
pub enum HardwareError {
    GpioError(String),
    SensorError(String),
    MotorError(String),
}

impl fmt::Display for HardwareError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            HardwareError::GpioError(msg) => write!(f, "GPIO Error: {}", msg),
            HardwareError::SensorError(msg) => write!(f, "Sensor Error: {}", msg),
            HardwareError::MotorError(msg) => write!(f, "Motor Error: {}", msg),
        }
    }
}

impl Error for HardwareError {}

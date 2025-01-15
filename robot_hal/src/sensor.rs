//! Sensor trait and implementations
use crate::error::HardwareError;

pub trait Sensor {
    type Reading;
    
    fn init(&mut self) -> Result<(), HardwareError>;
    fn read(&mut self) -> Result<Self::Reading, HardwareError>;
    fn cleanup(&mut self) -> Result<(), HardwareError>;
}
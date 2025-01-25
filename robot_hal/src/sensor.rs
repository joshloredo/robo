//! Sensor trait and implementations
use crate::error::HardwareError;
use log::info;

pub trait Sensor {
    type Reading;
    
    fn init(&mut self) -> Result<(), HardwareError>;
    fn read(&mut self) -> Result<Self::Reading, HardwareError>;
    fn cleanup(&mut self) -> Result<(), HardwareError>;
    fn say_hello(&mut self) -> Result<(), HardwareError>;
}

// Basic sensor implementation for testing
pub struct DummySensor;

impl DummySensor {
    pub fn new() -> Self {
        DummySensor
    }
}

impl Sensor for DummySensor {
    type Reading = f32;
    
    fn init(&mut self) -> Result<(), HardwareError> {
        Ok(())
    }
    
    fn read(&mut self) -> Result<Self::Reading, HardwareError> {
        Ok(0.0)
    }
    
    fn cleanup(&mut self) -> Result<(), HardwareError> {
        Ok(())
    }
    
    fn say_hello(&mut self) -> Result<(), HardwareError> {
        info!("Hello from DummySensor!");
        Ok(())
    }
}
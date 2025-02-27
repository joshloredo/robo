//! Motor control implementations
#[cfg(feature = "hardware")]
use rppal::gpio::{Gpio, OutputPin};
use crate::error::HardwareError;
use log::info;

#[derive(Debug, Clone, Copy)]
pub enum MotorDirection {
    Forward,
    Backward,
    Stop,
}

pub trait Motor {
    fn set_speed(&mut self, speed: f32) -> Result<(), HardwareError>;
    fn set_direction(&mut self, direction: MotorDirection) -> Result<(), HardwareError>;
    fn stop(&mut self) -> Result<(), HardwareError>;
    fn say_hello(&mut self) -> Result<(), HardwareError>;
}

// Hardware implementation for Raspberry Pi
#[cfg(feature = "hardware")]
pub struct N20Motor {
    pwm_pin: OutputPin,
    dir_pin1: OutputPin,
    dir_pin2: OutputPin,
    current_direction: MotorDirection,
    current_speed: f32,
}

#[cfg(feature = "hardware")]
impl N20Motor {
    pub fn new(pwm_pin_num: u8, dir_pin1_num: u8, dir_pin2_num: u8) -> Result<Self, HardwareError> {
        let gpio = Gpio::new().map_err(|e| HardwareError::GpioError(e.to_string()))?;
        
        let pwm_pin = gpio.get(pwm_pin_num)
            .map_err(|e| HardwareError::GpioError(e.to_string()))?
            .into_output();
            
        let dir_pin1 = gpio.get(dir_pin1_num)
            .map_err(|e| HardwareError::GpioError(e.to_string()))?
            .into_output();
            
        let dir_pin2 = gpio.get(dir_pin2_num)
            .map_err(|e| HardwareError::GpioError(e.to_string()))?
            .into_output();

        Ok(N20Motor {
            pwm_pin,
            dir_pin1,
            dir_pin2,
            current_direction: MotorDirection::Stop,
            current_speed: 0.0,
        })
    }
}

#[cfg(feature = "hardware")]
impl Motor for N20Motor {
    fn set_speed(&mut self, speed: f32) -> Result<(), HardwareError> {
        if speed < 0.0 || speed > 1.0 {
            return Err(HardwareError::MotorError("Speed must be between 0.0 and 1.0".into()));
        }
        self.current_speed = speed;
        Ok(())
    }

    fn set_direction(&mut self, direction: MotorDirection) -> Result<(), HardwareError> {
        match direction {
            MotorDirection::Forward => {
                self.dir_pin1.set_high();
                self.dir_pin2.set_low();
            }
            MotorDirection::Backward => {
                self.dir_pin1.set_low();
                self.dir_pin2.set_high();
            }
            MotorDirection::Stop => {
                self.dir_pin1.set_low();
                self.dir_pin2.set_low();
            }
        }
        
        self.current_direction = direction;
        Ok(())
    }

    fn stop(&mut self) -> Result<(), HardwareError> {
        self.set_direction(MotorDirection::Stop)?;
        self.set_speed(0.0)?;
        Ok(())
    }

    fn say_hello(&mut self) -> Result<(), HardwareError> {
        info!("Hello from Hardware N20Motor!");
        Ok(())
    }
}

// Simulation implementation
#[cfg(feature = "simulation")]
pub struct N20Motor {
    current_direction: MotorDirection,
    current_speed: f32,
}

#[cfg(feature = "simulation")]
impl N20Motor {
    pub fn new(_pwm_pin_num: u8, _dir_pin1_num: u8, _dir_pin2_num: u8) -> Result<Self, HardwareError> {
        Ok(N20Motor {
            current_direction: MotorDirection::Stop,
            current_speed: 0.0,
        })
    }
}

#[cfg(feature = "simulation")]
impl Motor for N20Motor {
    fn set_speed(&mut self, speed: f32) -> Result<(), HardwareError> {
        if speed < 0.0 || speed > 1.0 {
            return Err(HardwareError::MotorError("Speed must be between 0.0 and 1.0".into()));
        }
        self.current_speed = speed;
        info!("Simulation: Setting motor speed to {}", speed);
        Ok(())
    }

    fn set_direction(&mut self, direction: MotorDirection) -> Result<(), HardwareError> {
        self.current_direction = direction;
        info!("Simulation: Setting motor direction to {:?}", direction);
        Ok(())
    }

    fn stop(&mut self) -> Result<(), HardwareError> {
        self.set_direction(MotorDirection::Stop)?;
        self.set_speed(0.0)?;
        info!("Simulation: Motor stopped");
        Ok(())
    }

    fn say_hello(&mut self) -> Result<(), HardwareError> {
        info!("Hello from Simulated N20Motor!");
        Ok(())
    }
}
# Software Requirements
- 1. The system must be able to run on a Raspberry Pi 5
- 2. The system must be extensible to other hardware platforms
- 3. The system must be able to be built and run on a macOS machine
- 4. The system must provide a way to simulate the hardware for testing and development
- 5. The system must be configurable via a configuration file
- 6. The system must not use any external libraries or dependencies that are not available on the Raspberry Pi 5
- 7. The system must not panic or crash under any circumstances
- 8. The system will be written in Rust
- 9. The system should be modular and easy to understand and modify
- 10. The system can use any open source libraries or dependencies that are available on the Raspberry Pi 5
- 11. The system must log all errors and warnings to the console
- 12. The system should have a build system that is easy to understand and modify
- 13. The system should have a way to run the system in a simulated environment
- 14. The system should have a way to run the system in a real environment
- 15. The system should have a build system that generates a binary that can be run on the Raspberry Pi 5
- 16. The compilation process should be easy to understand and modify
- 17. The compilation process should log all errors and warnings to the console

# Hardware Requirements

## Core Computing
- Raspberry Pi 4B (2GB RAM minimum)
- 32GB Class 10 microSD Card
- USB-C Power Supply (RPi Official)

## Bill of Materials (Per Robot)

### Adafruit Components
| Item | Part Number | Quantity | Estimated Price (USD) |
|------|-------------|----------|-------------------|
| Raspberry Pi 4B (2GB) | 4292 | 1 | $45.00 |
| MPU6050 6-DoF IMU | 3886 | 1 | $7.50 |
| 12V DC Power Supply | 798 | 1 | $8.95 |
| Terminal Block | 677 | 1 | $5.95 |
| Breadboard-Friendly 2.1mm DC Jack | 373 | 1 | $2.95 |
| USB-C Power Supply (RPi) | 4298 | 1 | $7.95 |
| Logic Level Converter | 757 | 1 | $3.95 |

### Amazon/Generic Electronics
| Item | Specification | Quantity | Estimated Price (USD) |
|------|---------------|----------|-------------------|
| N20 DC Motors w/Encoders | 12V 100RPM | 2 | $24.00 |
| TB6612FNG Motor Driver | Dual H-Bridge | 1 | $6.95 |
| HC-SR04 Ultrasonic Sensors | 5V | 4 | $12.00 |
| 18650 Li-ion Battery Holder | 3-cell | 1 | $8.95 |
| 18650 Batteries | Protected cells | 3 | $21.00 |
| Battery Protection Circuit | 3S BMS | 1 | $8.95 |
| Jumper Wires Pack | M-M, F-F, M-F | 1 | $7.95 |
| microSD Card | 32GB Class 10 | 1 | $8.95 |

### Hardware (McMaster/Local)
| Item | Specification | Quantity | Estimated Price (USD) |
|------|---------------|----------|-------------------|
| M3 Screws Assortment | Various lengths | 1 pack | $8.00 |
| M3 Nuts | Standard | 1 pack | $4.00 |
| M3 Standoffs Kit | Various lengths | 1 kit | $12.00 |
| Ball Caster | 3/8" ball | 1 | $3.95 |

### 3D Printing Materials
| Item | Specification | Quantity (kg) | Estimated Price (USD) |
|------|---------------|---------------|-------------------|
| PETG Filament | 0.5kg used | 0.5 | $12.50 |
| TPU Filament | 0.1kg used | 0.1 | $3.00 |

### Optional Upgrades
| Item | Vendor | Quantity | Estimated Price (USD) |
|------|---------|----------|-------------------|
| RPLiDAR A1M8 | Amazon/AliExpress | 1 | $99.00 |
| Raspberry Pi Camera V2 | Adafruit #3099 | 1 | $29.95 |
| OLED Display 128x64 | Adafruit #326 | 1 | $19.95 |

## Cost Summary
- Basic Configuration: ~$224.45
- With Optional Upgrades: ~$373.35

# Software Requirements

## Operating System
- Raspberry Pi OS (64-bit) Latest Version
- Alternative: Ubuntu Server 22.04 LTS (64-bit)

## Core Languages
- Rust (Latest Stable)

## Rust Dependencies
- embedded-hal
- rppal
- tokio
- serialport
- nalgebra (for robotics calculations)

## Development Tools
- Visual Studio Code with extensions:
  - rust-analyzer
  - Go
  - Remote SSH
- git
- cargo
- cross (for cross-compilation)

## Build Dependencies
- build-essential
- cmake
- pkg-config
- libssl-dev

## Monitoring Tools
- Prometheus (optional)
- Grafana (optional)
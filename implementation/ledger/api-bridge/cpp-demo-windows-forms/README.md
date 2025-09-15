# Web4 Race Car Battery Management API Bridge - Windows Forms Demo

A Windows Forms-based C++ demo application for the Web4 Race Car Battery Management API Bridge, specifically designed for RAD Studio and Windows development environments. This demo provides a modern, user-friendly interface for testing and demonstrating the API Bridge functionality.

## Features

- **Modern Windows Forms Interface**: Native Windows UI with tabbed interface
- **Dual Interface Support**: Both REST and gRPC API interfaces
- **Privacy-Focused Features**: Anonymous component registration and verification
- **Complete API Coverage**: All blockchain operations including:
  - Account Management
  - Component Registry (Legacy)
  - Privacy-Focused Component Operations
  - LCT (Linked Context Token) Management
  - Pairing Process with Split-Key Generation
  - Pairing Queue Operations
  - Trust Tensor Operations
  - Energy Operations
- **Real-time Monitoring**: Battery status streaming via gRPC
- **Performance Testing**: REST vs gRPC comparison tools
- **RAD Studio Compatible**: Designed specifically for Windows development environments
- **Real Blockchain Integration**: No mock responses, all operations use actual blockchain transactions
- **Configurable UI**: Theme support, window positioning, and behavior customization

## Screenshots

The application features a modern tabbed interface with:
- Connection status indicators
- Tabbed sections for each API module
- Real-time log display
- Performance monitoring tools
- Configuration management

## Prerequisites

### Required Dependencies

- **Windows 10/11** (64-bit)
- **Visual Studio 2019/2022** or **RAD Studio 11/12**
- **C++17 compatible compiler** (MSVC)
- **CMake 3.16+**
- **Git**

### Optional Dependencies (for full functionality)

- **gRPC** - For gRPC interface support
- **Protobuf** - For protocol buffer support
- **OpenSSL** - For SSL/TLS support

## Building the Application

### Using Visual Studio

1. **Install Dependencies**:
   ```powershell
   # Install vcpkg (if not already installed)
   git clone https://github.com/Microsoft/vcpkg.git
   cd vcpkg
   .\bootstrap-vcpkg.bat
   .\vcpkg integrate install
   
   # Install required packages
   .\vcpkg install nlohmann-json:x64-windows
   .\vcpkg install grpc:x64-windows
   .\vcpkg install protobuf:x64-windows
   ```

2. **Build with CMake**:
   ```powershell
   mkdir build
   cd build
   cmake .. -DCMAKE_TOOLCHAIN_FILE=C:/vcpkg/scripts/buildsystems/vcpkg.cmake
   cmake --build . --config Release
   ```

3. **Alternative: Build with Visual Studio**:
   ```powershell
   cmake .. -G "Visual Studio 16 2019" -A x64
   cmake --build . --config Release
   ```

### Using RAD Studio

1. **Open the project**:
   - Launch RAD Studio
   - Open the `APIBridgeDemoWindowsForms.cbproj` project file
   - Configure the project settings for your target platform

2. **Build the project**:
   - Select Build → Build All Projects
   - The executable will be created in the output directory

## Configuration

The application uses a JSON configuration file (`config.json`) that can be customized:

```json
{
  "api": {
    "rest": {
      "endpoint": "http://localhost:8080",
      "timeout": 30
    },
    "grpc": {
      "endpoint": "localhost:9092",
      "timeout": 30
    }
  },
  "ui": {
    "window": {
      "width": 1200,
      "height": 800
    },
    "theme": {
      "use_dark_theme": false,
      "accent_color": "#0078D4"
    }
  }
}
```

## Usage

### Starting the Application

1. **Launch the executable**:
   ```powershell
   .\bin\APIBridgeDemoWindowsForms.exe
   ```

2. **Connect to API**:
   - Click the "Connect" button to establish connection to the API Bridge
   - The status bar will show connection status for both REST and gRPC

3. **Navigate through tabs**:
   - **Account**: Manage blockchain accounts
   - **Component**: Register and verify components
   - **Privacy**: Test privacy-focused features
   - **LCT**: Manage Linked Context Tokens
   - **Pairing**: Test pairing processes
   - **Pairing Queue**: Manage pairing requests
   - **Trust**: Test trust tensor operations
   - **Energy**: Test energy operations
   - **Performance**: Compare REST vs gRPC performance
   - **Logs**: View application logs

### Key Features

#### Connection Management
- **Connect/Disconnect**: Easily connect to and disconnect from the API
- **Status Indicators**: Real-time status of REST and gRPC connections
- **Auto-reconnect**: Automatic reconnection on connection loss

#### Tabbed Interface
- **Organized Sections**: Each API module has its own tab
- **Context-sensitive Controls**: Buttons and controls change based on connection status
- **Real-time Updates**: Tab content updates automatically

#### Logging System
- **Real-time Log Display**: View logs as they happen
- **Log Levels**: Filter logs by level (INFO, WARNING, ERROR, SUCCESS, DEBUG)
- **Export Functionality**: Export logs to file
- **Log Rotation**: Automatic log file rotation

#### Performance Testing
- **REST vs gRPC Comparison**: Compare performance between interfaces
- **Streaming Tests**: Test real-time data streaming
- **Benchmark Tools**: Measure response times and throughput

## Architecture

### Main Components

1. **MainForm**: Main application window and UI management
2. **ConfigManager**: Configuration file handling and settings management
3. **LogManager**: Logging system with file and console output
4. **RESTClient**: HTTP client for REST API communication
5. **GRPCClient**: gRPC client for streaming and binary communication

### UI Structure

```
MainForm
├── StatusBar (Connection status, REST/gRPC status)
├── TabControl
│   ├── Account Tab
│   ├── Component Tab
│   ├── Privacy Tab
│   ├── LCT Tab
│   ├── Pairing Tab
│   ├── Pairing Queue Tab
│   ├── Trust Tab
│   ├── Energy Tab
│   ├── Performance Tab
│   └── Logs Tab
└── Common Controls (Connect, Disconnect, Refresh)
```

## Development

### Project Structure

```
cpp-demo-windows-forms/
├── main.cpp                 # Application entry point
├── MainForm.h/cpp          # Main window and UI management
├── ConfigManager.h/cpp     # Configuration management
├── LogManager.h/cpp        # Logging system
├── RESTClient.h/cpp        # REST API client
├── GRPCClient.h/cpp        # gRPC client
├── CMakeLists.txt          # Build configuration
├── config.json             # Application configuration
└── README.md               # This file
```

### Adding New Features

1. **New Tab**: Add a new tab to the MainForm class
2. **New API Operations**: Extend the RESTClient or GRPCClient classes
3. **Configuration**: Add new settings to ConfigManager
4. **Logging**: Use LogManager for consistent logging

### Debugging

- **Debug Output**: Use `DebugOutput()` function for debug messages
- **Log Levels**: Set log level in configuration for detailed logging
- **Error Handling**: All errors are logged and displayed to user

## Troubleshooting

### Common Issues

1. **Connection Failed**:
   - Verify API Bridge server is running
   - Check endpoint configuration in `config.json`
   - Ensure firewall allows connections

2. **Build Errors**:
   - Verify all dependencies are installed
   - Check CMake configuration
   - Ensure Visual Studio/RAD Studio is properly configured

3. **gRPC Not Available**:
   - Install gRPC and Protobuf dependencies
   - Check CMake configuration
   - Verify network connectivity

### Performance Issues

- **Slow Response**: Check network latency and server performance
- **Memory Usage**: Monitor memory usage in Task Manager
- **UI Responsiveness**: Check for blocking operations in UI thread

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the repository
- Contact: support@web4.com
- Documentation: https://web4.com/docs

## Version History

- **v1.0.0**: Initial Windows Forms implementation
  - Basic UI with tabbed interface
  - REST and gRPC client integration
  - Configuration management
  - Logging system
  - Performance testing tools 
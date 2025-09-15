cmake_minimum_required(VERSION 3.16)
project(APIBridgeDemoWindowsForms VERSION 1.0.0 LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Set output directories
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

# Windows-specific settings
if(WIN32)
    # Windows-specific settings
    add_definitions(-D_WIN32_WINNT=0x0601)
    add_definitions(-DWIN32_LEAN_AND_MEAN)
    add_definitions(-DUNICODE)
    add_definitions(-D_UNICODE)
    
    # For RAD Studio compatibility
    if(CMAKE_CXX_COMPILER_ID STREQUAL "MSVC")
        add_compile_options(/W4 /utf-8 /EHsc)
        add_compile_definitions(_CRT_SECURE_NO_WARNINGS)
    endif()
else()
    message(FATAL_ERROR "This project is designed for Windows only")
endif()

# Find required packages
find_package(PkgConfig REQUIRED)

# Find httplib (header-only HTTP library)
find_package(httplib QUIET)
if(NOT httplib_FOUND)
    # If not found, we'll use FetchContent to download it
    include(FetchContent)
    FetchContent_Declare(
        httplib
        GIT_REPOSITORY https://github.com/yhirose/cpp-httplib.git
        GIT_TAG v0.14.1
    )
    FetchContent_MakeAvailable(httplib)
endif()

# Find nlohmann/json
find_package(nlohmann_json QUIET)
if(NOT nlohmann_json_FOUND)
    FetchContent_Declare(
        nlohmann_json
        GIT_REPOSITORY https://github.com/nlohmann/json.git
        GIT_TAG v3.11.2
    )
    FetchContent_MakeAvailable(nlohmann_json)
endif()

# Find gRPC
find_package(gRPC QUIET)
if(NOT gRPC_FOUND)
    # For demo purposes, we'll provide a fallback implementation
    # In production, install gRPC properly
    message(WARNING "gRPC not found. gRPC functionality will be limited.")
    set(GRPC_AVAILABLE FALSE)
else()
    set(GRPC_AVAILABLE TRUE)
endif()

# Find Protobuf
find_package(Protobuf QUIET)
if(NOT Protobuf_FOUND)
    message(WARNING "Protobuf not found. gRPC functionality will be limited.")
    set(PROTOBUF_AVAILABLE FALSE)
else()
    set(PROTOBUF_AVAILABLE TRUE)
endif()

# Create source groups for better IDE organization
file(GLOB_RECURSE HEADER_FILES "*.h")
file(GLOB_RECURSE SOURCE_FILES "*.cpp")

# Main executable
add_executable(APIBridgeDemoWindowsForms WIN32
    main.cpp
    MainForm.cpp
    ConfigManager.cpp
    LogManager.cpp
    RESTClient.cpp
    GRPCClient.cpp
)

# Set target properties
set_target_properties(APIBridgeDemoWindowsForms PROPERTIES
    OUTPUT_NAME "APIBridgeDemoWindowsForms"
    DEBUG_POSTFIX "_d"
    RELEASE_POSTFIX ""
    WIN32_EXECUTABLE TRUE
)

# Include directories
target_include_directories(APIBridgeDemoWindowsForms PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_BINARY_DIR}
)

# Link libraries
target_link_libraries(APIBridgeDemoWindowsForms PRIVATE
    httplib::httplib
    nlohmann_json::nlohmann_json
)

# Add gRPC support if available
if(GRPC_AVAILABLE AND PROTOBUF_AVAILABLE)
    target_link_libraries(APIBridgeDemoWindowsForms PRIVATE
        gRPC::grpc++
        gRPC::grpc++_reflection
        protobuf::libprotobuf
    )
    
    # Add gRPC definitions
    target_compile_definitions(APIBridgeDemoWindowsForms PRIVATE GRPC_ENABLED=1)
else()
    target_compile_definitions(APIBridgeDemoWindowsForms PRIVATE GRPC_ENABLED=0)
endif()

# Windows-specific libraries
target_link_libraries(APIBridgeDemoWindowsForms PRIVATE
    ws2_32
    iphlpapi
    crypt32
    comctl32
    user32
    gdi32
    shell32
    ole32
    uuid
    oleaut32
    advapi32
)

# Installation
install(TARGETS APIBridgeDemoWindowsForms
    RUNTIME DESTINATION bin
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
)

# Copy configuration files
install(FILES
    config.json
    DESTINATION bin
)

# Create a package target for distribution
include(CPack)

# Set CPack properties
set(CPACK_PACKAGE_NAME "APIBridgeDemoWindowsForms")
set(CPACK_PACKAGE_VERSION ${PROJECT_VERSION})
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "Web4 Race Car Battery Management API Bridge Demo - Windows Forms")
set(CPACK_PACKAGE_VENDOR "Web4 Team")
set(CPACK_PACKAGE_CONTACT "support@web4.com")

set(CPACK_GENERATOR "NSIS;ZIP")
set(CPACK_NSIS_DISPLAY_NAME "Web4 API Bridge Demo - Windows Forms")
set(CPACK_NSIS_PACKAGE_NAME "Web4 API Bridge Demo - Windows Forms")
set(CPACK_NSIS_CONTACT "support@web4.com")
set(CPACK_NSIS_URL "https://web4.com")
set(CPACK_NSIS_ENABLE_NSIS64 "ON")

# Testing
enable_testing()

# Add a simple test
add_test(NAME APIBridgeDemoWindowsFormsTest
    COMMAND APIBridgeDemoWindowsForms --test
    WORKING_DIRECTORY ${CMAKE_BINARY_DIR}/bin
)

# Documentation
find_package(Doxygen QUIET)
if(DOXYGEN_FOUND)
    set(DOXYGEN_IN ${CMAKE_CURRENT_SOURCE_DIR}/Doxyfile.in)
    set(DOXYGEN_OUT ${CMAKE_CURRENT_BINARY_DIR}/Doxyfile)
    
    configure_file(${DOXYGEN_IN} ${DOXYGEN_OUT} @ONLY)
    
    add_custom_target(docs ALL
        COMMAND ${DOXYGEN_EXECUTABLE} ${DOXYGEN_OUT}
        WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
        COMMENT "Generating API documentation with Doxygen"
        VERBATIM
    )
endif()

# Print configuration summary
message(STATUS "=== APIBridgeDemoWindowsForms Configuration ===")
message(STATUS "C++ Standard: ${CMAKE_CXX_STANDARD}")
message(STATUS "Build Type: ${CMAKE_BUILD_TYPE}")
message(STATUS "gRPC Available: ${GRPC_AVAILABLE}")
message(STATUS "Protobuf Available: ${PROTOBUF_AVAILABLE}")
message(STATUS "Install Prefix: ${CMAKE_INSTALL_PREFIX}")
message(STATUS "Output Directory: ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}")
message(STATUS "===============================================") 
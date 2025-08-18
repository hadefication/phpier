# Feature Specification: redis-memcached-cli

## Overview
Add CLI commands to access Redis and Memcached services directly from the phpier CLI tool. This feature provides developers with convenient access to Redis CLI and Memcached telnet interface without needing to remember Docker container names or connection details.

## Requirements

### Core Functionality
- `phpier redis` - Launch Redis CLI connected to the Redis/Valkey container
- `phpier memcached` - Launch telnet interface connected to Memcached container
- Both commands should work when services are running
- Commands should provide helpful error messages when services are not available
- Support for passing additional arguments to the underlying tools

### User Interface Requirements
- Commands should follow existing phpier CLI patterns
- Provide clear error messages when containers are not running
- Support `--help` flag for usage information
- Maintain consistency with other phpier proxy commands

### Integration Requirements
- Work with existing Docker Compose setup
- Compatible with current service naming conventions
- Support both Redis and Valkey containers (Redis-compatible)
- Work across all supported PHP versions and project configurations

### Service Requirements
- Detect if Redis/Valkey service is running before attempting connection
- Detect if Memcached service is running before attempting connection
- Handle cases where services are defined but not started

## Implementation Notes

### Technical Considerations
- Use Docker exec to access container CLIs
- Follow existing command proxy patterns from other tools (php, composer, etc.)
- Redis/Valkey container typically named `redis` or `valkey` in compose files
- Memcached container typically named `memcached` in compose files
- Need to handle both service names for Redis (redis/valkey compatibility)

### Files to Modify
- `cmd/redis.go` - New Redis command
- `cmd/memcached.go` - New Memcached telnet command  
- `cmd/root.go` - Register new commands
- Update command registration and help text

### Integration Points
- Docker service detection utilities
- Existing container execution patterns
- Error handling for service unavailability
- Command argument passing mechanisms

### Implementation Patterns
- Follow existing proxy command structure
- Use Cobra command patterns consistent with other commands
- Implement proper error handling with meaningful messages
- Support argument forwarding to underlying tools

### Testing Strategy
- Test with Redis/Valkey containers running and stopped
- Test with Memcached containers running and stopped  
- Test argument forwarding functionality
- Test error handling when services are unavailable
- Test across different project configurations

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete
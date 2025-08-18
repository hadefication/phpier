# PHPier Troubleshooting Guide

## Common Issues and Solutions

### Docker Issues

```bash
# Check Docker is running
docker --version
docker-compose --version

# Check if ports are available
docker ps
netstat -tulpn | grep :80
```

### Permission Issues

```bash
# Fix file permissions
sudo chown -R $USER:$USER .
chmod -R 755 .
```

### Container Issues

```bash
# View project logs
cd your-project/.phpier
docker-compose logs

# View specific service logs
docker-compose logs app

# Rebuild project from scratch
phpier down --remove-volumes
phpier up --build

# Reset global services
phpier stop
phpier start --force
```

### Domain Access Issues

```bash
# Test Traefik routing
curl -H "Host: your-project.localhost" http://localhost

# Check DNS resolution
ping your-project.localhost
```

## Supported PHP Versions

phpier supports all major PHP versions with appropriate tooling:

| PHP Version | Status | Extensions Included |
|-------------|--------|-------------------|
| 5.6 | Legacy Support | Basic web development extensions |
| 7.3 | Legacy Support | Full extension set |
| 7.4 | Active Support | Full extension set + modern tools |
| 8.0+ | Active Support | Latest extensions + performance optimizations |

Each container includes Composer, NVM (Node.js), and appropriate tooling for the PHP version.

## Service Diagnostics

### Check Service Status

```bash
# Global services status
phpier global up

# Project service status
phpier up

# All Docker containers
docker ps -a
```

### Port Conflicts

```bash
# Check what's using common ports
lsof -i :80
lsof -i :3306
lsof -i :8080

# Kill conflicting processes
sudo kill -9 <PID>
```

### Network Issues

```bash
# List Docker networks
docker network ls

# Inspect phpier networks
docker network inspect phpier-global

# Clean up orphaned networks
docker network prune
```

### Volume Issues

```bash
# List Docker volumes
docker volume ls

# Inspect project volumes
docker volume inspect <project-name>_data

# Remove specific volumes
docker volume rm <volume-name>

# Clean up orphaned volumes
docker volume prune
```

## Performance Issues

### Slow Container Startup

1. **Check disk space**: `df -h`
2. **Prune Docker**: `docker system prune -a`
3. **Restart Docker**: Restart Docker Desktop/daemon
4. **Use faster storage**: Move to SSD if using HDD

### Memory Issues

```bash
# Check Docker memory usage
docker stats

# Increase Docker memory limits in Docker Desktop
# Or adjust PHP memory limits in .phpier/docker/php/php.ini
```

## Reset Procedures

### Complete Reset

```bash
# Stop all phpier services
phpier stop --force

# Remove all project containers and volumes
phpier down --remove-volumes

# Clean up Docker completely
docker system prune -a --volumes

# Restart fresh
phpier start
phpier up --build
```

### Reset Single Project

```bash
# From project directory
phpier down --remove-volumes
rm -rf .phpier/
phpier init <php-version>
phpier up --build
```

## Getting Help

### Debug Information

When reporting issues, include:

```bash
# System information
uname -a
docker --version
docker-compose --version
phpier version

# Service status
docker ps -a
phpier status  # if available

# Recent logs
docker-compose logs --tail=50
```

### Log Locations

- **Project logs**: `.phpier/` directory - `docker-compose logs`
- **Global service logs**: Docker logs for global containers
- **Application logs**: Usually in your project's `storage/logs` or similar

### Common Error Messages

#### "Port already in use"
- Another service is using the same port
- Use `lsof -i :<port>` to find the conflicting process
- Stop the conflicting service or change phpier ports

#### "Network not found"
- Global services aren't running
- Run `phpier start` to start global services

#### "Container name already exists"
- Previous containers weren't cleaned up properly
- Run `phpier down` or `docker rm <container-name>`

#### "Permission denied"
- File permission issues
- Fix with `sudo chown -R $USER:$USER .` and `chmod -R 755 .`

## Support Resources

- 🐛 **Issues**: [GitHub Issues](https://github.com/your-org/phpier/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/your-org/phpier/discussions)  
- 📖 **Documentation**: [Wiki](https://github.com/your-org/phpier/wiki)
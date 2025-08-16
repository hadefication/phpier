# phpier
A CLI tool to manage PHP development using Docker.

Supports PHP 5.6, 7.3., 7.4, 8.0, 8.1, 8.2, 8.3, 8.4

Use traefik for folder base domain i.e. <directory>.localhost when dev up is running.

PHPMyAdmin if MySQL is installed


# Containers
- App container that is pre-installed with all the tools commonly used for PHP development such as php, php-fpm, common php extensions, composer, nvm, nginx
- Database (select from mysql, postgress, mariadb) containers
- Valkey (Redis) support
- Memcached support
- Mailpit

## CLI
- phpier init <version> -- Initialize a phpier environment using <version> of PHP.
- phpier up -- Run the docker-compose run command
- phpier down -- Run the docker-compose down command
- phpier build -- Build or rebuild services again
- phpier php -- Proxy the php in the app container
- phpier <app_container_tools> -- Proxt the tools in the app container

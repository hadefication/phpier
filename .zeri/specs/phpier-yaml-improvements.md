# Feature Specification: Shared Services Architecture

## Overview
This specification outlines a major architectural shift for phpier. The goal is to move from a project-specific, all-in-one stack to a model where non-app services (databases, caches, tools) are shared across multiple projects. Each project will consist of a lightweight, standalone app container that connects to this central, shared services stack.

This change will significantly reduce resource consumption, simplify management, and speed up the initialization of new projects.

## Requirements

### 1. Global Shared Services Stack
- A new set of CLI commands (`phpier global up/down/status`) will be created to manage a persistent, global stack of services.
- This global stack will include:
  - Traefik (as the central reverse proxy)
  - Database services (MySQL, PostgreSQL, MariaDB)
  - Cache services (Redis, Memcached)
  - Tooling services (Mailpit, PHPMyAdmin, etc.)
- All global services will run on a dedicated, shared Docker network (e.g., `phpier_global`).
- The configuration for this global stack will be stored in a central location (e.g., `$HOME/.phpier/global-compose.yml`).

### 2. Project-Specific App Containers
- The `phpier init` command will be modified to generate a `docker-compose.yml` for the **app container only** (PHP + Nginx).
- This app container will be configured to connect to the `phpier_global` network.
- The project's `.phpier.yaml` will contain the necessary connection details for the shared services (e.g., database host, Redis port). The database host will point to the container name in the global stack (e.g., `phpier-mysql`).

### 3. Traefik for Service Discovery and Routing
- The global Traefik instance will automatically discover and route traffic to project-specific app containers.
- Routing will be based on the project's directory name, as is currently the case (e.g., `my-project.localhost`).
- Traefik will need to be configured to monitor Docker for new containers on the `phpier_global` network.

### 4. Configuration Changes (`.phpier.yaml`)
- The `.phpier.yaml` file will be simplified. It will no longer define the versions or settings for shared services.
- It will need a new section to specify connection details for the shared services. For example:
  ```yaml
  services:
    database:
      host: phpier-mysql
      port: 3306
      # username, password, database name are still project-specific
      username: my_project_user
      password: my_project_password
      database: my_project_db
  ```

## Implementation Notes
- **New CLI Commands:** Introduce `phpier global <sub-command>` to manage the lifecycle of the shared services stack.
- **Docker Compose Refactoring:** The template engine will need to generate two different types of `docker-compose.yml` files: one for the global stack and one for individual projects.
- **Networking:** The key to this architecture is the shared Docker network. The global stack will create the network, and project containers will connect to it as an external network.
- **Backwards Compatibility:** This is a breaking change. We will need to provide a clear migration path for users with existing projects. A `phpier migrate` command could be considered.

## TODO
- [x] Design and plan implementation of the global services stack.
- [x] Implement `phpier global up` and `phpier global down` commands.
- [x] Modify the `phpier init` command to generate project-specific app container configurations.
- [x] Update the configuration structure (`internal/config/config.go`) to reflect the new shared model.
- [x] Adjust the Docker Compose templates (`internal/templates`) for both global and project-specific stacks.
- [x] Ensure Traefik correctly discovers and routes to the app containers.
- [x] Build and test the shared services architecture successfully.
- [ ] Add comprehensive tests for the new architecture.
- [ ] Update documentation to explain the new shared services model and migration steps.
- [ ] Review and refine the implementation.
- [ ] Mark specification as complete.

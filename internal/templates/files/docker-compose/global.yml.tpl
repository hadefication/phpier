name: phpier

services:
  traefik:
    image: traefik:v2.10
    container_name: phpier-traefik
    restart: unless-stopped
    ports:
      - "{{.Global.Traefik.Port}}:80"
      - "{{.Global.Traefik.SSLPort}}:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik/traefik.yml:/etc/traefik/traefik.yml:ro
      - ./traefik/dynamic:/etc/traefik/dynamic:ro
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.traefik.rule=Host(`phpier-traefik.{{.Global.Traefik.Domain}}`)"
      - "traefik.http.routers.traefik.entrypoints=web"
      - "traefik.http.services.traefik.loadbalancer.server.port=8080"

  {{if eq .Global.Services.Database.Type "mysql"}}
  mysql:
    image: mysql:{{.Global.Services.Database.Version}}
    container_name: phpier-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "{{.Global.Services.Database.Port}}:3306"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.mysql.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.mysql.entrypoints=mysql"
      - "traefik.tcp.services.mysql.loadbalancer.server.port=3306"
  {{end}}

  {{if eq .Global.Services.Database.Type "postgresql"}}
  postgres:
    image: postgres:{{.Global.Services.Database.Version}}
    container_name: phpier-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: phpier
      POSTGRES_USER: phpier
      POSTGRES_PASSWORD: phpier
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "{{.Global.Services.Database.Port}}:5432"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.postgres.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.postgres.entrypoints=postgres"
      - "traefik.tcp.services.postgres.loadbalancer.server.port=5432"
  {{end}}

  {{if eq .Global.Services.Database.Type "mariadb"}}
  mariadb:
    image: mariadb:{{.Global.Services.Database.Version}}
    container_name: phpier-mariadb
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: phpier
      MYSQL_USER: phpier
      MYSQL_PASSWORD: phpier
    volumes:
      - mariadb_data:/var/lib/mysql
    ports:
      - "{{.Global.Services.Database.Port}}:3306"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.mariadb.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.mariadb.entrypoints=mariadb"
      - "traefik.tcp.services.mariadb.loadbalancer.server.port=3306"
  {{end}}

  {{if serviceEnabled "redis" .Global}}
  redis:
    image: redis:alpine
    container_name: phpier-redis
    restart: unless-stopped
    volumes:
      - redis_data:/data
    ports:
      - "{{.Global.Services.Cache.Redis.Port}}:6379"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.redis.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.redis.entrypoints=redis"
      - "traefik.tcp.services.redis.loadbalancer.server.port=6379"
  {{end}}

  {{if serviceEnabled "mailpit" .Global}}
  mailpit:
    image: axllent/mailpit
    container_name: phpier-mailpit
    restart: unless-stopped
    ports:
      - "1026:1025"
      - "8026:8025"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.mailpit.rule=Host(`phpier-mailpit.{{.Global.Traefik.Domain}}`)"
      - "traefik.http.routers.mailpit.entrypoints=web"
      - "traefik.http.services.mailpit.loadbalancer.server.port=8025"
  {{end}}

  # Web interfaces for database and cache management
  {{if eq .Global.Services.Database.Type "mysql"}}
  adminer:
    image: adminer:4.8.1
    container_name: phpier-adminer
    restart: unless-stopped
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.adminer.rule=Host(`phpier-mysql.{{.Global.Traefik.Domain}}`)"
      - "traefik.http.routers.adminer.entrypoints=web"
      - "traefik.http.services.adminer.loadbalancer.server.port=8080"
  {{end}}

  {{if serviceEnabled "redis" .Global}}
  redis-commander:
    image: rediscommander/redis-commander:latest
    container_name: phpier-redis-commander
    restart: unless-stopped
    environment:
      - REDIS_HOSTS=redis:phpier-redis:6379
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.redis-commander.rule=Host(`phpier-redis.{{.Global.Traefik.Domain}}`)"
      - "traefik.http.routers.redis-commander.entrypoints=web"
      - "traefik.http.services.redis-commander.loadbalancer.server.port=8081"
  {{end}}

networks:
  {{.Global.Network}}:
    name: {{.Global.Network}}
    driver: bridge

volumes:
  {{if eq .Global.Services.Database.Type "mysql"}}
  mysql_data:
  {{end}}
  {{if eq .Global.Services.Database.Type "postgresql"}}
  postgres_data:
  {{end}}
  {{if eq .Global.Services.Database.Type "mariadb"}}
  mariadb_data:
  {{end}}
  {{if serviceEnabled "redis" .Global}}
  redis_data:
  {{end}}

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

  {{if .Global.Services.Databases.MySQL.Enabled}}
  mysql:
    image: mysql:{{.Global.Services.Databases.MySQL.Version}}
    container_name: phpier-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: {{.Global.Services.Databases.MySQL.Password}}
      MYSQL_DATABASE: {{.Global.Services.Databases.MySQL.Database}}
      MYSQL_USER: {{.Global.Services.Databases.MySQL.Username}}
      MYSQL_PASSWORD: {{.Global.Services.Databases.MySQL.Password}}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "{{.Global.Services.Databases.MySQL.Port}}:3306"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.mysql.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.mysql.entrypoints=mysql"
      - "traefik.tcp.services.mysql.loadbalancer.server.port=3306"
  {{end}}

  {{if .Global.Services.Databases.PostgreSQL.Enabled}}
  postgres:
    image: postgres:{{.Global.Services.Databases.PostgreSQL.Version}}
    container_name: phpier-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: {{.Global.Services.Databases.PostgreSQL.Database}}
      POSTGRES_USER: {{.Global.Services.Databases.PostgreSQL.Username}}
      POSTGRES_PASSWORD: {{.Global.Services.Databases.PostgreSQL.Password}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "{{.Global.Services.Databases.PostgreSQL.Port}}:5432"
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.tcp.routers.postgres.rule=HostSNI(`*`)"
      - "traefik.tcp.routers.postgres.entrypoints=postgres"
      - "traefik.tcp.services.postgres.loadbalancer.server.port=5432"
  {{end}}

  {{if .Global.Services.Databases.MariaDB.Enabled}}
  mariadb:
    image: mariadb:{{.Global.Services.Databases.MariaDB.Version}}
    container_name: phpier-mariadb
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: {{.Global.Services.Databases.MariaDB.Password}}
      MYSQL_DATABASE: {{.Global.Services.Databases.MariaDB.Database}}
      MYSQL_USER: {{.Global.Services.Databases.MariaDB.Username}}
      MYSQL_PASSWORD: {{.Global.Services.Databases.MariaDB.Password}}
    volumes:
      - mariadb_data:/var/lib/mysql
    ports:
      - "{{.Global.Services.Databases.MariaDB.Port}}:3306"
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
  {{if or .Global.Services.Databases.MySQL.Enabled .Global.Services.Databases.MariaDB.Enabled}}
  adminer:
    image: adminer:latest
    container_name: phpier-adminer
    restart: unless-stopped
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.adminer.rule=Host(`phpier-adminer.{{.Global.Traefik.Domain}}`)"
      - "traefik.http.routers.adminer.entrypoints=web"
      - "traefik.http.services.adminer.loadbalancer.server.port=8080"
  {{end}}

networks:
  {{.Global.Network}}:
    driver: bridge

volumes:
  {{if .Global.Services.Databases.MySQL.Enabled}}
  mysql_data:
  {{end}}
  {{if .Global.Services.Databases.PostgreSQL.Enabled}}
  postgres_data:
  {{end}}
  {{if .Global.Services.Databases.MariaDB.Enabled}}
  mariadb_data:
  {{end}}
  {{if serviceEnabled "redis" .Global}}
  redis_data:
  {{end}}

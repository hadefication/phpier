name: {{.Config.Docker.ProjectName}}

services:
  app:
    build:
      context: ..
      dockerfile: .phpier/Dockerfile.php
    container_name: {{.Config.Docker.ProjectName}}-app
    restart: unless-stopped
    volumes:
      - ..:/var/www/html
      - ./docker/php/php.ini:/usr/local/etc/php/php.ini
      - ./docker/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./docker/nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
      - ./docker/supervisor/supervisord.conf:/etc/supervisor/conf.d/supervisord.conf:ro
    networks:
      - {{.Config.Docker.Network}}
    environment:
      - PHP_VERSION={{.Config.PHP.Version}}
    working_dir: /var/www/html
{{- if serviceEnabled "traefik" .Config }}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}.rule=Host(`{{.Config.Docker.ProjectName}}.{{.Config.Traefik.Domain}}`)"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}.entrypoints=web"
      - "traefik.http.services.{{.Config.Docker.ProjectName}}.loadbalancer.server.port=80"
      - "traefik.docker.network={{.Config.Docker.Network}}"
{{- else }}
    ports:
      - "{{.Config.Services.Webserver.Port}}:80"
{{- end }}

{{- if eq .Config.Services.Database.Type "mysql" }}
  mysql:
    image: mysql:{{.Config.Services.Database.Version}}
    container_name: {{.Config.Docker.ProjectName}}-mysql
    restart: unless-stopped
    command: --default-authentication-plugin=mysql_native_password
    environment:
      MYSQL_ROOT_PASSWORD: {{.Config.Services.Database.Password}}
      MYSQL_DATABASE: {{.Config.Services.Database.Database}}
      MYSQL_USER: {{.Config.Services.Database.Username}}
      MYSQL_PASSWORD: {{.Config.Services.Database.Password}}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "{{.Config.Services.Database.Port}}:3306"
    networks:
      - {{.Config.Docker.Network}}
{{- end }}

{{- if eq .Config.Services.Database.Type "postgresql" }}
  postgres:
    image: postgres:{{.Config.Services.Database.Version}}
    container_name: {{.Config.Docker.ProjectName}}-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: {{.Config.Services.Database.Database}}
      POSTGRES_USER: {{.Config.Services.Database.Username}}
      POSTGRES_PASSWORD: {{.Config.Services.Database.Password}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "{{.Config.Services.Database.Port}}:5432"
    networks:
      - {{.Config.Docker.Network}}
{{- end }}

{{- if eq .Config.Services.Database.Type "mariadb" }}
  mariadb:
    image: mariadb:{{.Config.Services.Database.Version}}
    container_name: {{.Config.Docker.ProjectName}}-mariadb
    restart: unless-stopped
    command: --default-authentication-plugin=mysql_native_password
    environment:
      MYSQL_ROOT_PASSWORD: {{.Config.Services.Database.Password}}
      MYSQL_DATABASE: {{.Config.Services.Database.Database}}
      MYSQL_USER: {{.Config.Services.Database.Username}}
      MYSQL_PASSWORD: {{.Config.Services.Database.Password}}
    volumes:
      - mariadb_data:/var/lib/mysql
    ports:
      - "{{.Config.Services.Database.Port}}:3306"
    networks:
      - {{.Config.Docker.Network}}
{{- end }}

{{- if serviceEnabled "redis" .Config }}
  redis:
    image: redis:alpine
    container_name: {{.Config.Docker.ProjectName}}-redis
    restart: unless-stopped
    ports:
      - "{{.Config.Services.Cache.Redis.Port}}:6379"
    networks:
      - {{.Config.Docker.Network}}
{{- end }}

{{- if serviceEnabled "memcached" .Config }}
  memcached:
    image: memcached:alpine
    container_name: {{.Config.Docker.ProjectName}}-memcached
    restart: unless-stopped
    ports:
      - "{{.Config.Services.Cache.Memcached.Port}}:11211"
    networks:
      - {{.Config.Docker.Network}}
{{- end }}

{{- if serviceEnabled "phpmyadmin" .Config }}
  phpmyadmin:
    image: phpmyadmin/phpmyadmin
    container_name: {{.Config.Docker.ProjectName}}-phpmyadmin
    restart: unless-stopped
    environment:
      PMA_HOST: {{if eq .Config.Services.Database.Type "mysql"}}{{.Config.Docker.ProjectName}}-mysql{{else if eq .Config.Services.Database.Type "mariadb"}}{{.Config.Docker.ProjectName}}-mariadb{{end}}
      PMA_USER: {{.Config.Services.Database.Username}}
      PMA_PASSWORD: {{.Config.Services.Database.Password}}
{{- if not (serviceEnabled "traefik" .Config) }}
    ports:
      - "8080:80"
{{- end }}
    networks:
      - {{.Config.Docker.Network}}
{{- if serviceEnabled "traefik" .Config }}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-phpmyadmin.rule=Host(`phpmyadmin.{{.Config.Traefik.Domain}}`)"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-phpmyadmin.entrypoints=web"
      - "traefik.http.services.{{.Config.Docker.ProjectName}}-phpmyadmin.loadbalancer.server.port=80"
{{- end }}
    depends_on:
{{- if eq .Config.Services.Database.Type "mysql" }}
      - mysql
{{- else if eq .Config.Services.Database.Type "mariadb" }}
      - mariadb
{{- end }}
{{- end }}

{{- if serviceEnabled "pgadmin" .Config }}
  pgadmin:
    image: dpage/pgadmin4
    container_name: {{.Config.Docker.ProjectName}}-pgadmin
    restart: unless-stopped
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
{{- if not (serviceEnabled "traefik" .Config) }}
    ports:
      - "8081:80"
{{- end }}
    networks:
      - {{.Config.Docker.Network}}
{{- if serviceEnabled "traefik" .Config }}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-pgadmin.rule=Host(`pgadmin.{{.Config.Traefik.Domain}}`)"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-pgadmin.entrypoints=web"
      - "traefik.http.services.{{.Config.Docker.ProjectName}}-pgadmin.loadbalancer.server.port=80"
{{- end }}
    depends_on:
      - postgres
{{- end }}

{{- if serviceEnabled "mailpit" .Config }}
  mailpit:
    image: axllent/mailpit
    container_name: {{.Config.Docker.ProjectName}}-mailpit
    restart: unless-stopped
{{- if not (serviceEnabled "traefik" .Config) }}
    ports:
      - "8025:8025"
      - "{{.Config.Services.Tools.Mailpit.Port}}:1025"
{{- else }}
    ports:
      - "{{.Config.Services.Tools.Mailpit.Port}}:1025"  # SMTP port always exposed
{{- end }}
    networks:
      - {{.Config.Docker.Network}}
{{- if serviceEnabled "traefik" .Config }}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-mailpit.rule=Host(`mailpit.{{.Config.Traefik.Domain}}`)"
      - "traefik.http.routers.{{.Config.Docker.ProjectName}}-mailpit.entrypoints=web"
      - "traefik.http.services.{{.Config.Docker.ProjectName}}-mailpit.loadbalancer.server.port=8025"
{{- end }}
{{- end }}

{{- if serviceEnabled "traefik" .Config }}
  traefik:
    image: traefik:v2.10
    container_name: {{.Config.Docker.ProjectName}}-traefik
    restart: unless-stopped
    ports:
      - "{{.Config.Traefik.Port}}:80"
      - "{{.Config.Traefik.SSLPort}}:443"
      - "8080:8080"  # Traefik dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./.phpier/traefik/traefik.yml:/etc/traefik/traefik.yml:ro
      - ./.phpier/traefik/dynamic:/etc/traefik/dynamic:ro
      - ./.phpier/traefik/data:/data
      - ./.phpier/traefik/logs:/var/log
    networks:
      - {{.Config.Docker.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.traefik.rule=Host(`traefik.{{.Config.Traefik.Domain}}`)"
      - "traefik.http.routers.traefik.entrypoints=web"
      - "traefik.http.services.traefik.loadbalancer.server.port=8080"
{{- end }}

networks:
  {{.Config.Docker.Network}}:
    driver: bridge

volumes:
{{- if eq .Config.Services.Database.Type "mysql" }}
  mysql_data:
{{- end }}
{{- if eq .Config.Services.Database.Type "postgresql" }}
  postgres_data:
{{- end }}
{{- if eq .Config.Services.Database.Type "mariadb" }}
  mariadb_data:
{{- end }}
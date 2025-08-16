name: {{.Project.Name}}

services:
  php:
    build:
      context: ..
      dockerfile: .phpier/Dockerfile.php
    container_name: {{.Project.Name}}-php
    restart: unless-stopped
    volumes:
{{- range $volume := .Project.App.Volumes}}
      - {{$volume}}
{{- end}}
      - ./.phpier/logs/nginx:/var/log/nginx
      - ./.phpier/logs/php:/var/log/php
      - ./.phpier/logs/supervisor:/var/log/supervisor
{{- if .Project.App.Environment}}
    environment:
{{- range $env := .Project.App.Environment}}
      - {{$env}}
{{- end}}
{{- end}}
    networks:
      - {{.Global.Network}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Project.Name}}.rule={{getHostRule .Project .Global}}"
      - "traefik.http.routers.{{.Project.Name}}.entrypoints=web"
      - "traefik.http.services.{{.Project.Name}}.loadbalancer.server.port=80"
      - "traefik.docker.network={{.Global.Network}}"

networks:
  {{.Global.Network}}:
    external: true
    name: phpier_{{.Global.Network}}

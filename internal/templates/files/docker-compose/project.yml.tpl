name: {{.Project.Name}}

services:
  app:
    build:
      context: .
      dockerfile: .phpier/Dockerfile.php
    image: phpier-{{.Project.Name}}:{{.Project.PHP}}
    container_name: {{.Project.Name}}-app
    restart: unless-stopped
    volumes:
{{- range $volume := .Project.App.Volumes}}
      - {{$volume}}
{{- end}}
      - ./.phpier/logs/nginx:/var/log/nginx
      - ./.phpier/logs/php:/var/log/php
      - ./.phpier/logs/supervisor:/var/log/supervisor
    environment:
      - WWWUSER=${WWWUSER}
{{- if .Project.App.Environment}}
{{- range $env := .Project.App.Environment}}
      - {{$env}}
{{- end}}
{{- end}}
    networks:
      - {{.Global.Network}}
    labels:
      # Traefik configuration
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Project.Name}}.rule={{getHostRule .Project .Global}}"
      - "traefik.http.routers.{{.Project.Name}}.entrypoints=web"
      - "traefik.http.services.{{.Project.Name}}.loadbalancer.server.port=80"
      - "traefik.docker.network={{.Global.Network}}"
      # Phpier metadata
      - "phpier.project.name={{.Project.Name}}"
      - "phpier.project.php={{.Project.PHP}}"
      - "phpier.project.node={{.Project.Node}}"
      - "phpier.managed=true"

networks:
  {{.Global.Network}}:
    external: true
    name: phpier_{{.Global.Network}}

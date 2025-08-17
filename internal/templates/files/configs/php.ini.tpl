; Custom PHP Configuration for phpier
; This file contains sensible defaults - edit directly to customize

; Memory and execution limits
memory_limit = 256M
max_execution_time = 300
max_input_time = 300
max_input_vars = 3000

; File upload limits
upload_max_filesize = 64M
post_max_size = 64M

; Error handling and logging
display_errors = On
display_startup_errors = On
error_reporting = E_ALL
log_errors = On
error_log = /var/log/php_errors.log

; Development settings
html_errors = On

; Date and timezone
date.timezone = UTC

; Session configuration
session.save_handler = files
session.save_path = /tmp
session.gc_maxlifetime = 3600
session.cookie_lifetime = 0

; OPcache settings (if enabled)
opcache.enable = 1
opcache.enable_cli = 1
opcache.memory_consumption = 128
opcache.interned_strings_buffer = 8
opcache.max_accelerated_files = 4000
opcache.revalidate_freq = 2
opcache.fast_shutdown = 1
opcache.validate_timestamps = 1

; Security settings
expose_php = Off
allow_url_fopen = On
allow_url_include = Off

; MySQL default settings
mysqli.default_port = 3306
mysqli.default_socket = /var/run/mysqld/mysqld.sock

; PostgreSQL default settings
pgsql.allow_persistent = On
pgsql.auto_reset_persistent = Off

; Output buffering
output_buffering = 4096
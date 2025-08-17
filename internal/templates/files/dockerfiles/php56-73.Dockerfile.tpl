FROM php:{{.Project.PHP}}-fpm

# Set working directory
WORKDIR /var/www/html

# Fix Debian repositories for older PHP versions (Stretch/Buster/Bullseye EOL)
RUN sed -i 's|http://deb.debian.org/debian|http://archive.debian.org/debian|g' /etc/apt/sources.list \
    && sed -i 's|http://security.debian.org/debian-security|http://archive.debian.org/debian|g' /etc/apt/sources.list \
    && sed -i 's|bullseye-security|bullseye|g' /etc/apt/sources.list \
    && sed -i '/stretch-updates/d' /etc/apt/sources.list \
    && sed -i '/buster-updates/d' /etc/apt/sources.list \
    && sed -i '/bullseye-updates/d' /etc/apt/sources.list

# Install system dependencies for older PHP versions  
RUN apt-get update && apt-get install -y --allow-unauthenticated \
    git \
    curl \
    unzip \
    nginx \
    supervisor \
    gosu \
    libpng-dev \
    libjpeg-dev \
    libfreetype6-dev \
    libonig-dev \
    libxml2-dev \
    libzip-dev \
    libicu-dev \
    libpq-dev \
    libmagickwand-dev \
    libcurl4-openssl-dev \
    libssl-dev \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# Configure PHP extensions for older versions
RUN docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/

# Install core PHP extensions available in older versions
RUN docker-php-ext-install -j$(nproc) \
    bcmath \
    calendar \
    curl \
    dom \
    exif \
    ftp \
    gd \
    intl \
    mbstring \
    mysqli \
    opcache \
    pdo \
    pdo_mysql \
    pdo_pgsql \
    pgsql \
    soap \
    sockets \
    tokenizer \
    xml \
    zip

# Install PECL extensions compatible with older PHP versions
RUN pecl install redis-4.3.0 \
    && docker-php-ext-enable redis

# Install Composer (version compatible with PHP version)
{{- if or (eq .Project.PHP "5.6") (eq .Project.PHP "7.0") (eq .Project.PHP "7.1") (eq .Project.PHP "7.2") (eq .Project.PHP "7.3") }}
# Use Composer 2.2.x for older PHP versions (last version supporting PHP 7.2+)
COPY --from=composer:2.2 /usr/bin/composer /usr/bin/composer
{{- else }}
# Use latest Composer for modern PHP versions (7.4+)
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer
{{- end }}

# Node.js installation skipped for PHP 5.6 to avoid compatibility issues with Debian Stretch
# If you need Node.js with PHP 5.6, consider using a newer PHP version or manual installation

# Copy custom PHP configuration
COPY .phpier/docker/php/php.ini /usr/local/etc/php/conf.d/custom.ini

# Configure Nginx
COPY .phpier/docker/nginx/nginx.conf /etc/nginx/nginx.conf
COPY .phpier/docker/nginx/default.conf /etc/nginx/sites-available/default
RUN ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

# Configure Supervisor
COPY .phpier/docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy entrypoint script and make it executable
COPY .phpier/docker/entrypoint.sh /usr/local/bin/start
RUN chmod +x /usr/local/bin/start

# Create phpier user for permission mapping
RUN useradd -ms /bin/bash -u 1337 phpier

# Create www-data user directories
RUN mkdir -p /var/www/html && chown www-data:www-data /var/www/html

# Expose port
EXPOSE 80

# Set entrypoint
ENTRYPOINT ["start"]

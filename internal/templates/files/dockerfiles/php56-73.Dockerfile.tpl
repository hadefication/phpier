FROM php:{{.Project.PHP}}-fpm

# Set working directory
WORKDIR /var/www/html

# Install system dependencies for older PHP versions
RUN apt-get update && apt-get install -y \
    git \
    curl \
    unzip \
    nginx \
    supervisor \
    libpng-dev \
    libjpeg-dev \
    libfreetype6-dev \
    libonig-dev \
    libxml2-dev \
    libzip-dev \
    libicu-dev \
    libpq-dev \
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

# Install Composer compatible with older PHP versions
COPY --from=composer:2.2 /usr/bin/composer /usr/bin/composer

# Install Node.js (if configured)
{{- if shouldInstallNode .Project.Node }}
{{- $nodeVersion := resolveNodeVersion .Project.Node }}
{{- if eq $nodeVersion "lts" }}
# Install Node.js LTS
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest
{{- else if eq $nodeVersion "16" }}
# Install Node.js 16.x (latest available)
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest
{{- else if eq $nodeVersion "18" }}
# Install Node.js 18.x (latest available)
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest
{{- else if eq $nodeVersion "20" }}
# Install Node.js 20.x (latest available)
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest
{{- else if eq $nodeVersion "22" }}
# Install Node.js 22.x (latest available)
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest
{{- else }}
# Install specific Node.js version: {{ $nodeVersion }}
{{- $majorVersion := index (split $nodeVersion ".") 0 }}
RUN curl -fsSL https://deb.nodesource.com/setup_{{ $majorVersion }}.x | bash - \
    && apt-get install -y nodejs={{ $nodeVersion }}-1nodesource1 \
    && npm install -g npm@latest
{{- end }}
{{- else }}
# Node.js installation skipped (node: none)
{{- end }}

# Copy custom PHP configuration
COPY .phpier/docker/php/php.ini /usr/local/etc/php/conf.d/custom.ini

# Configure Nginx
COPY .phpier/docker/nginx/nginx.conf /etc/nginx/nginx.conf
COPY .phpier/docker/nginx/default.conf /etc/nginx/sites-available/default
RUN ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

# Configure Supervisor
COPY .phpier/docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy startup script and make it executable
COPY .phpier/docker/startup.sh /usr/local/bin/startup.sh
RUN chmod +x /usr/local/bin/startup.sh

# Create www-data user directories
RUN mkdir -p /var/www/html && chown www-data:www-data /var/www/html

# Expose port
EXPOSE 80

# Use startup script instead of direct supervisor
CMD ["/usr/local/bin/startup.sh"]

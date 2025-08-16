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

# Install Node.js compatible with older PHP versions
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=14.21.3
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm use $NODE_VERSION \
    && nvm alias default $NODE_VERSION

# Add Node.js to PATH
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

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

FROM mcr.microsoft.com/devcontainers/go:1.21 as base

### Installing dependencies

# Setup Node.js 18.x
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install ng command
RUN npm install -g @angular/cli

# Install gcloud comamnd, via package manager
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update -y && \
    apt-get install google-cloud-cli -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

FROM base

# Frontend port
EXPOSE 4200
# Backend Port
EXPOSE 8080

# Declare the environment variables
ENV KHI_FRONTEND_STATIC_FILE_FOLDER, \
    KHI_SERVER_HOST, \
    KHI_SERVER_PORT, \
    KHI_VIEWER_MODE

FROM node:20-alpine
WORKDIR /app

# Copy dependency manifests first and install dependencies
COPY package*.json ./
RUN npm ci

# Copy the rest of the app
COPY server.js .

# Ensure logs folder exists
RUN mkdir -p /app/logs

EXPOSE 8080
CMD ["node", "server.js"]
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json first to leverage Docker cache
COPY package*.json ./
RUN npm ci

# Copy the rest of the application code
COPY . .

# Build the SvelteKit application
RUN npm run build

# Prune dev dependencies from node_modules
RUN npm prune --production

FROM node:22-alpine

WORKDIR /app

# Copy only the essential files from the builder stage
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .

# Set environment variable for production
ENV NODE_ENV production

# Expose the port where the Node.js server will listen (default for adapter-node is 3000)
EXPOSE 3000

# Command to start the SvelteKit application
CMD ["node", "build"]
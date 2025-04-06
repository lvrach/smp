# Build stage
FROM oven/bun:1 AS builder
WORKDIR /app

# Copy package files
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install

# Copy source files
COPY . .

# Build the application
RUN bun run build

# Runtime stage
FROM oven/bun:1-slim
WORKDIR /app

# Copy only the built application from builder
COPY --from=builder /app/build/index.js ./build/index.js

# Set the command to run the application
CMD ["bun", "build/index.js"]

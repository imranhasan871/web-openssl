# Use Node.js 18 Alpine
FROM node:18-alpine

# Install dependencies for compatibility
RUN apk add --no-cache libc6-compat

# Set working directory
WORKDIR /app

# Copy package files first
COPY package.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Set environment variables
ENV NODE_ENV=development

# Start the development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Expose port
EXPOSE 3000
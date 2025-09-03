# Use the official Node.js runtime as the base image
FROM node:18-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install ALL dependencies (including dev dependencies for build)
RUN npm ci

# Copy the rest of the application code
COPY . .

# Set environment variables for the build process
ARG NEXT_PUBLIC_GOOGLE_MAPS_API_KEY
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_ENVIRONMENT

# Expose ARGs as ENV variables for the build process
ENV NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=$NEXT_PUBLIC_GOOGLE_MAPS_API_KEY
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_ENVIRONMENT=$NEXT_PUBLIC_ENVIRONMENT

# --- DIAGNOSTIC STEP ---
RUN echo "Debug: NEXT_PUBLIC_API_URL is set to '$NEXT_PUBLIC_API_URL'"
# --- END DIAGNOSTIC STEP ---

# Build the Next.js application
RUN npm run build

# Remove dev dependencies but keep TypeScript for Next.js config
RUN npm prune --production && npm install typescript

# Expose the port the app runs on
EXPOSE 3000

# Start the application
CMD ["npm", "start"]

# Use the latest Node.js LTS image from Docker Hub
FROM node:latest

# Set working directory inside the container
WORKDIR /usr/src/app

# Copy package.json if you have dependencies (optional)
# COPY package*.json ./
# RUN npm install --production

# Copy all project files into the container
COPY . .

# Run your main.js file with Node.js
CMD ["node", "entrypoint.js"]
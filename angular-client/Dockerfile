# Create image based on the official Node 8 image from docker
FROM node:9

# Create a directory where our app will be placed
RUN mkdir -p /usr/src/app

# Change directoyr so that our commands run insdide the new directory
WORKDIR /usr/src/app

# Copy dependency definitions
COPY package.json /usr/src/app

# Install dependencies
RUN npm install

# Get all the code needed to run the app
COPY . /usr/src/app

# Expose the port the app runs in
EXPOSE 4200

# Serve the app
CMD [ "npm", "start" ]

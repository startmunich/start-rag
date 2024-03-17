# Use an official redis runtime as the base image
FROM redis

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Make port 5000 (default for flask) available to the world outside this container
EXPOSE 6379
    

# Run when the container launches
CMD ["redis-server", "--save", "60", "1", "--loglevel", "warning"]
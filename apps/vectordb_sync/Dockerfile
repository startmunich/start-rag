# Use an official Python runtime as the base image
FROM python:3.11-slim

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed dependencies specified in requirements.txt
RUN pip install -r requirements.txt

# Make port 5000 (default for flask) available to the world outside this container
EXPOSE 5000
    

# Run when the container launches
CMD ["python", "vectordb_sync.py"]
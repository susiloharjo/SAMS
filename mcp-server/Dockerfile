# Use an official Python runtime as a parent image
FROM python:3.11-slim

# Set the working directory in the container
WORKDIR /app

# Install uv package manager
RUN pip install uv

# Copy the requirements file into the container at /app
# We use pyproject.toml which is uv's standard
COPY pyproject.toml ./

# Install any needed packages specified in pyproject.toml
RUN uv pip install --system -r pyproject.toml

# Copy the rest of the application's code into the container
COPY . .

# Make port 8081 available to the world outside this container
EXPOSE 8081

# Run server.py when the container launches
CMD ["uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8081"]

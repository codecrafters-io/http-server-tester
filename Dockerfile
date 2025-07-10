FROM caddy:2.9.1-alpine

# Install Python
RUN apk add --no-cache python3

# Create the data directory
RUN mkdir -p /tmp/data/codecrafters.io/http-server-tester/

# Copy the Caddyfile and Python script
COPY Caddyfile /etc/caddy/Caddyfile
COPY file_writer.py /file_writer.py

# Expose the ports
EXPOSE 4221 4222

# Run both servers
CMD sh -c "python3 /file_writer.py & caddy run --config /etc/caddy/Caddyfile" 
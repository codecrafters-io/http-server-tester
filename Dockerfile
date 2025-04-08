FROM caddy:2.9.1-alpine

# Create the data directory
RUN mkdir -p /tmp/data/codecrafters.io/http-server-tester/

# Copy the Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Expose the port
EXPOSE 4221

# Run Caddy
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"] 
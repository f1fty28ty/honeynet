FROM python:3.11-alpine
WORKDIR /app

# Copy honeypot script
COPY server.py .

# Create logs directory for persistent file logging
RUN mkdir -p /app/logs

EXPOSE 3306

CMD ["python", "server.py"]
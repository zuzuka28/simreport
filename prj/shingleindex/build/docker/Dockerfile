FROM python:3.11-slim as builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libffi-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN pip install --upgrade pip \
    && pip install --no-cache-dir -r requirements.txt

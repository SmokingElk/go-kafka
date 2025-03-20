#!/bin/bash

echo "Waiting for Kafka to be ready..."

while ! nc -z kafka1 29091; do
  echo "Waiting for kafka1:29091..."
  sleep 1
done

while ! nc -z kafka2 29092; do
  echo "Waiting for kafka2:29092..."
  sleep 1
done

while ! nc -z kafka3 29093; do
  echo "Waiting for kafka3:29093..."
  sleep 1
done

echo "Kafka is ready!"
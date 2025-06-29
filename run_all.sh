#!/bin/bash

# Run backend
echo "Starting backend..."
cd backend
go run server.go &
cd ..

# Run frontend
echo "Starting frontend..."
cd frontend
npm install
npm run dev

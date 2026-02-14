FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
# We use --host 0.0.0.0 so it's accessible outside the container
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

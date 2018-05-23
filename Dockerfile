FROM node:9-alpine

EXPOSE 8000

# COPY . /home/node/app
#WORKDIR /home/node/app
#RUN npm install
WORKDIR /home/node/app
CMD npm start
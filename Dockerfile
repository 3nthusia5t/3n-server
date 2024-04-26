FROM golang:1.20


ENV APP_REPO "https://"${APP_TOKEN}"@github.com/3nthusia5t/3n-app.git"
ENV SERVER_REPO "https://"${APP_TOKEN}"@github.com/3nthusia5t/3n-server.git"
ENV ARTICLES_REPO "https://"${APP_TOKEN}"@github.com/3nthusia5t/3n-articles.git"

# Install tools
RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y npm
RUN apt-get install -y cron
RUN apt-get install -y certbot


WORKDIR /app
RUN git clone ${APP_REPO}
RUN git clone ${SERVER_REPO}
RUN git clone ${ARTICLES_REPO}

#Compile the backend service
WORKDIR /app/3n-server/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/3n-server.bin


#Compile the frontend
WORKDIR /app/3n-app
RUN npm install
RUN npm run build

#Set up a cronjob
RUN mv /app/3n-server/utils/cron/cronfile /etc/cron.d/cronfile
RUN chmod 0644 /etc/cron.d/cronfile
RUN chmod 777 /app/3n-server/utils/scripts/update.sh
RUN chmod 777 /app/3n-server/utils/cron/routine.sh
RUN crontab /etc/cron.d/cronfile
#Final setup
WORKDIR /app

#Load articles 
RUN /app/3n-server.bin transcompile --src /app/3n-articles/markdown --dst /app/3n-articles/html
RUN /app/3n-server.bin update -a /app/3n-articles/html

#Prepare entry point
RUN mv /app/3n-server/entry.sh /app
RUN chmod 777 /app/entry.sh

#Expose port 80
EXPOSE 80
EXPOSE 100
EXPOSE 443
#Serve website
CMD ["/app/entry.sh"]
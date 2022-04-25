FROM redis:latest
VOLUME /data
EXPOSE 6379
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
# syntax=docker/dockerfile:1

# 1.==============================================================
# The build container will only compile flutter project
FROM fischerscode/flutter:stable AS build
WORKDIR /flutter
# switch to root user to change folder permission for flutter CLI
# https://stackoverflow.com/a/24555761/12328041
USER root
RUN chown -R flutter:flutter .
USER flutter
# permission needs to be in accordance with flutter user
COPY --chown=flutter:flutter . .
# get the dependancies before building
RUN flutter pub get
# build flutter web platform
RUN flutter build -v web

# 2.==============================================================
# The runtime container will be used to build the final image that runs Nginx
FROM nginx AS runtime
WORKDIR /data/www
# copy the compiled flutter web static files from the container at the build stage
COPY --from=build /flutter/build/web ./
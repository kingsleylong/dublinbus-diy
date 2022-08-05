# syntax=docker/dockerfile:1

# This Dockerfile uses multi-stage builds technic to both read credential files
# and save space in the final image.
# Refer to https://docs.docker.com/develop/develop-images/multistage-build/

# 1.=======================IMPORT CONFIG===========================
#
# !!!Important!!! The dublinbus-config image needs to be built locally before
# building this Dockerfile. The repo url for dublinbus-config projeact is
# https://csgitlab.ucd.ie/kingsley/dublinbus-config.
#
# This FROM instruction imports the dublinbus-config image which contains the
# compiling time configuration (web/index.html). It should be placed in the
# first stage - configs - so that the later stages can use it.
FROM dublinbus-config AS configs
RUN ls -lR /config

# 2.=======================BUILD FLUTTER===========================
# The build container will only compile flutter project
FROM fischerscode/flutter:stable AS build
# use ARG to pass the build time arguments
# https://docs.docker.com/engine/reference/builder/#arg
ARG API_HOST
WORKDIR /flutter
# switch to root user to change folder permission for flutter CLI
# https://stackoverflow.com/a/24555761/12328041
USER root
RUN chown -R flutter:flutter .
USER flutter
# permission needs to be in accordance with flutter user
COPY --chown=flutter:flutter . .
# remove ios and android folder before building to avoid errors - we only build web platform on the server
RUN rm -r android ios
# copy the compiling time config file to the source folder.
# the web/index.html can't be loaded using volume as it is required by flutter
# when building the release packages.
COPY --from=configs /config/web/index.html /flutter/web/index.html
# get the dependancies before building
RUN flutter pub get
# build flutter web platform with verbos information
# https://docs.docker.com/engine/reference/builder/#using-arg-variables
RUN flutter build -v web

# 3.=======================BUILD NGINX==========================
# The runtime container will be used to build the final image that runs Nginx
FROM nginx AS runtime
WORKDIR /data/www
# copy the compiled flutter web static files from the container at the build stage
COPY --from=build /flutter/build/web ./
# the volume for storing static files such as installation files, QR codes, etc
VOLUME ["/data/files"]
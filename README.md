# MEX

Media EXchange for Plex

## Purpose

After attempting to configure Lidarr, Radarr, Sonarr, Ombi, and Transmission 
multiple times for my Plex server, I realized that all of these are hanging
together by a thread. Rather than performing further debugging I did what any
good Software Engineer does: Start a new Open Source project to solve my
problem and share it with the world!

Since the ultimate goal is to deploy MEX in a single Docker container, I have
followed the pattern created by the Gitea self-hosted Git service. The entire
project compiles to a single Docker container that runs both the UI and all
of the backend services. There is nothing else to deploy. There is no
configuration of external services necessary. Everything "just works".

## Building

To build the executable, which is statically linked and contains all of the
static assets (HTML files, CSS files, images, etc.) simply type:

    make

Yes, that really is all that's necessary

## Docker

After building the executable, the Docker container can be built with:

    docker build

Again, that's all that's necessary.

## Running

The executable does need some information from environment variables in order
to correctly access some of the web services used to search for media. This
is necessary because most of the services, such as TVDB and TMDB require per-
user API keys. I could embed a set of keys into the application like Radarr
and Sonarr, but if for any reason the provider decided to revoke my keys for
generating too many API calls, all users would be affected. That is why each
of you is responsible for creating an account with the services and getting
an API key.

Running locally:

    TVDB_API_KEY=abc TMDB_API_KEY=xyz ./mex

Running with Docker:

    docker run -d \
        -e "TVDB_API_KEY=abc" \
        -e "TMDB_API_KEY=xyz" \
        -p 9000:9000 \
        --name mex \
        --restart unless-stopped \
        mex:latest

Alternatively, `docker-compose` may be used with the included compose file:

    TVDB_API_KEY=abc TMDB_API_KEY=xyz docker compose up

## Contributing

1.  Fork it
2.  Create a feature branch (`git checkout -b new-feature`)
3.  Commit changes (`git commit -am "Added new feature xyz"`)
4.  Push the branch (`git push origin new-feature`)
5.  Create a new pull request.

## Maintainers

* [Media Exchange](http://github.com/MediaExchange/)

## License

Copyright 2019 MediaExchange.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


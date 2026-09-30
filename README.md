# IlPost Podcast Feed

Turns the podcasts of [Il Post](https://www.ilpost.it/) into standard RSS feeds,
so you can listen to them with any podcast app (Apple Podcasts, Overcast, ...)
instead of the official one.

**An Il Post subscription is still required.** The service does not give free
access: the feed asks for your ilpost.it credentials (HTTP Basic auth) and uses
them to fetch the episodes on your behalf.

Hosted instance: https://ilpost-feed.totaro.net

If you are not comfortable typing your password on a third-party server, you can host it
yourself.

## Usage

* `/` home page with the list of podcasts
* `/podcast-list` podcast list as JSON
* `/feed?podcast-name=<slug>` RSS feed of a podcast (Basic auth with your ilpost.it username and password)

## Build and run

Requires Go 1.24 or later.

    go build
    ./ilpost-podcast-feed

It listens on port 8080 and serves `./static`, so run it from the repository root.

## systemd

A sample unit is in `packaging/systemd/ilpost.service`:

    sudo cp packaging/systemd/ilpost.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable --now ilpost

# TODO
* add link for itunes on index.html
* add my user-agent when getting ilpost api so they recognize us
* refactor the /test endpoint by putting the testdata in go:embed bordone.json (and move it to testdata folder)
* review tester reports
    * review https://www.castfeedvalidator.com/validate.php?url=https://ilpost-feed.totaro.net/test
    * review https://validator.w3.org/feed/check.cgi?url=https://ilpost-feed.totaro.net/test
    * review https://podba.se/validate/?url=https://ilpost-feed.totaro.net/test 
* unit systemd: doesn't allow comments and may not find binary
* index cookie cache also by hashed pwd
* protect CookieCache with mutex
* update GO version
* deploy to docker
# Steam History / Data collector [![Build Status](https://travis-ci.org/steamhistory/collector.svg?branch=master)](https://travis-ci.org/steamhistory/collector) [![GoDoc](https://godoc.org/github.com/steamhistory/collector?status.png)](https://godoc.org/github.com/steamhistory/collector)

**[Collector](https://github.com/steamhistory/collector)** | [Storage](https://github.com/steamhistory/storage) | [Processor](https://github.com/steamhistory/processor) | [Reporter](https://github.com/steamhistory/reporter)

Steam distributes thousands of games and used by millions of PC gamers every day to play games and interact with their friends and the rest of the community. Steam provides [some stats](http://store.steampowered.com/stats) about usage, but there is not much info available. Steam History project tries to solve part of this problem by recording usage history for *all* apps distributed on Steam.

### Data sources

There are two Steam Web API interfaces that are used to collect usage info:

1. [ISteamApps/GetAppList](https://api.steampowered.com/ISteamApps/GetAppList/v2/) - list of all apps available in Steam Store;
2. [ISteamUserStats/GetNumberOfCurrentPlayers](https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=0) - current number of users for a specified app.

First we need to get list of apps to know their IDs and names. This list is saved in a local database and serves as a reference during usage history recording process. It is worth mentioning that GetAppList method returns not only applications, but also every single DLC, video, soundtrack and other not\_so\_useful stuff. These "apps" always have 0 active users so there is no need to save their usage history.

After creating list of apps, history recording process becomes pretty straightforward: get number of users for each usable app periodically and save it in a database. List of all apps needs to be updated too. We don’t want to miss any release.

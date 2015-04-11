# Steam History [![Build Status](https://travis-ci.org/gentlecat/steamhistory.svg?branch=master)](https://travis-ci.org/gentlecat/steamhistory) [![GoDoc](https://godoc.org/github.com/gentlecat/steamhistory?status.png)](https://godoc.org/github.com/gentlecat/steamhistory)

Steam distributes thousands of games and used by millions of PC gamers every day to play games and interact with their friends and the rest of the community. Steam provides [some stats](http://store.steampowered.com/stats) about usage, but there is not much info available. Steam History project tries to solve part of this problem by allowing you to record usage history for *all* apps distributed on Steam.

### Data sources

There are two Steam Web API interfaces that are used to collect usage info:

1. [ISteamApps/GetAppList](https://api.steampowered.com/ISteamApps/GetAppList/v2/) - list of all apps available in Steam Store;
2. [ISteamUserStats/GetNumberOfCurrentPlayers](https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=0) - current number of users for a specified app.

First we need to get list of apps to know their IDs and names. This list is saved in a local database and serves as a reference during usage history recording process. It is worth mentioning that GetAppList method returns not only applications, but also every single DLC, video, soundtrack and other not\_so\_useful stuff. These "apps" always have 0 active users so there is no need to save their usage history.

After creating list of apps, history recording process becomes pretty straightforward: get number of users for each usable app periodically and save it in a database. List of all apps needs to be updated too. We don’t want to miss any release.

### Storage

For simplicity, SQLite database is created for each application. Each database contains only one table with two columns: timestamp (PK) and number of users. Most queries work with a single app. For example, adding a new usage history record or retrieving all of them.

List of all apps and information about them is saved in another database. This information includes: app ID, name, and boolean value that indicates if app is usable.

### Detecting unusable apps

Some apps returned by Steam Web API cannot be run in Steam Client. That means their user count is always 0, which makes them easily detectable. Detection of these apps is done periodically by marking apps as unusable if their average user count is less then 1.

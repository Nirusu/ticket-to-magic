# ticket-to-magic
A very basic Go / OAuth implementation that queries the API used by Disneyland's website to determine the availability of tickets for Disneyland Park & California Adventure Park in Anaheim, California.

## Story
I noticed that tickets can appear very spontaneously before a given date even if they are sold out for a long time already. This is likely due to people chancing the reservation, but might help you if you *really* want to go on a certain date as I did and want to be able to grab some tickets which are only available for a couple of minutes.

This is a very basic tool that just prints any results to your console in intervals of 10 seconds. It does not do any filtering beyond checking for available tickets before the target date nor send any notifications. 

Feel free to take this implementation as a reference for your own application, or just let it running on a second monitor if you are lazy. In case you are modifying it, keep in mind that Disney actually has some User-Agent filtering in place so you cannot just query the endpoints using curl without spoofing the agent. In case this tool prints a 403, it could be the case that Disney denylisted the user-agent, which you can easily replace with a common browser user-agent in the code.

In case of othre errors, the API might have changed..

This tool was originally developed in March 2022 and last tested in September 2022.

## Build

You need to have [Go](https://go.dev/) installed. This was tested using Go 1.18 and Go 1.19. 

Clone the repository, enter the directory and build the tool (`.exe` only for Windows):
```bash
go build -o ticket-to-magic[.exe] main.go
```

That's it.


## Usage
Linux / macOS / UNIX:
```bash
./ticket-to-magic yyyy-mm-dd
```

Windows:
```bash
.\ticket-to-magic.exe yyyy-mm-dd
```

Where **yyyy-mm-dd** is the end date for which the tool should search tickets for (from now on).

Filtering for a specific day is left to be implemented by the user, can be done by piping the stdout output of the tool into `grep` (Linux / macOS / Windows using WSL):
```bash
# If you want to search specifically for tickets on November 30, 2022
./ticket-to-magic 2022-12-01 | grep 2022-11-30

## Specifically for Disneyland Park on November 30, 2022
./ticket-to-magic 2022-12-01 | grep "2022-11-30: Disneyland Park is available"

## Specifically for California Adventure Park on November 30, 2022
./ticket-to-magic 2022-12-01 | grep "2022-11-30: California Adventure Park is available"
````

The Windows equivalent using PowerShell:
```powershell
# If you want to search specifically for tickets on November 30, 2022
.\ticket-to-magic.exe 2022-12-01 | Select-String 2022-11-30

## Specifically for Disneyland Park on November 30, 2022
.\ticket-to-magic.exe 2022-12-01 | Select-String "2022-11-30: Disneyland Park is available"

## Specifically for California Adventure Park on November 30, 2022
.\ticket-to-magic.exe 2022-12-01 | Select-String "2022-11-30: California Adventure Park is available"
````

Alternatively, just adjust the filtering in the code. The time filtering uses Go's `time` package, so tweaking it to your desire should be relatively easy and a good coding practice in case you are unfamiliar with Go :)
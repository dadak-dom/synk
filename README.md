# README

## About

This is the official Wails Vue-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

TODO:

FRONTEND:

NOW:

- Finish the PeerList and SynkButton parts of the Mainview
- Display the user's local IP
- In the main view, create a subcomponent that can list the files that are currently in the shared folder (\*) - ideally in like a drop-down menu

LATER:

- Make the navbar pretty.
- Overall, think of a cohesive aesthetic to create for the app, and then make the assets for it
- Expanding on (\*): show a comparison of the files on the local vs. on the remote machine before doing the synk (as a sort of confirmation).
  Would need to create a graphical rendering of the filesystem (nested folders, etc.)

BACKEND:
NOW:

- Generate a random string of words or numbers on startup. Use the FRONTEND to confirm that you want to synk with the computer with that ID
- only launch the API when actually running a synk. (Less drastic measure: toggle the API to off when not synking.)

LATER (LONG TERM):

- Find a way to allow users to automatically synk so long as both computers are on
- Find a way to "remember" other devices and automatically synk with them

- Folder Selector:
  - make it prettier
- Multicast
  - Make sure that the library works when used on the same LAN network
- File sharing
  - code up the file sending / receiving logic

FIXME:

- Make it so that the shared folder is the default when opening the folder selector
- For the time being, folders don't work. Will need to implement some sort of check that creates a subdirectory if it doesn't already exist
- There is a bug (on windows) where the selected local IP for the http server is on a different subnet. Will need to make it so that the API is only launched once a connection has been made via multicast
- Will need to make it so that the program is always listening for new peers in the background. Will have to rejigger the logic a little bit

NOTES:

- For some reason, the multicasting only works between windows/linux when I change the multicast address (e.g. from 224.0.0.0 to 224.0.0.1)

NOTE!:

- On Linux builds, DO NOT use CSS rules like this:
  .outer-rule {
  .inner-rule
  }

This breaks the inner rule
